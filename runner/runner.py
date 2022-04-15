"""Drift Analysis HPC Runner."""

import json
from operator import contains
import re
from argparse import ArgumentParser
from dataclasses import dataclass
from datetime import datetime, timezone
import subprocess

import dateutil.parser
import requests
from helicase import Helicase
from spack.cmd.diff import compare_specs
from spack.error import SpecError
from spack.spec import Spec
from spack.version import ver

# Define SUCCESS for comparing command return codes.
SUCCESS = 0


@dataclass
class Result:
    """Result is an internal represenation of an Inflection Point."""

    abstract_spec: str
    commit: str
    tags: list
    files: list
    primary: bool
    author_date: str
    commit_date: str
    concretized: bool
    concretizer: str
    concretizer_out: str
    concrete_spec: str


class DriftAnalysis(Helicase):
    """
    DriftAnalysis is the main class of the runner.

    Building off of Helicase this class traverses Spack's git history to look
    for inflection points.
    """

    def __init__(self, specs=None):
        """Set up the last cache and specs used for the analysis."""
        self.abstract_specs = specs or []
        self.last_spec = {}
        self.last_success = {
            abstract_spec: True for abstract_spec in self.abstract_specs
        }
        self.previous_points = {
            abstract_spec: None for abstract_spec in self.abstract_specs
        }

        if args.skip_known_inflection_points == "true":
            for abstract_spec in self.abstract_specs:
                r = requests.get(f"{args.host}/inflection-points/{abstract_spec}")
                for point in r.json():
                    point_date = dateutil.parser.parse(point["GitCommitDate"])
                    if point["Concretizer"] == args.concretizer and \
                       (self.previous_points[abstract_spec] is None or
                       point_date > self.previous_points[abstract_spec]):
                        self.previous_points[abstract_spec] = point_date

    def analyze(self, commit):
        """
        Analyze defines the main history traversal function.

        Analyze defines the process that should be run on every commit
        and for every abstract spec. This is used to determine if the
        commit is an "inflection point". Inflection points occur when
        the concretization for an abstract spec changes. This most
        often happens when a package updates or one of its dependencies
        update leading to a different installation for the same input
        abstract spec.
        """
        for abstract_spec in self.abstract_specs:
            # Create spec for parsing info
            spec = Spec(abstract_spec)

            if args.skip_known_inflection_points == "true" and \
               commit.committer_date < self.previous_points[abstract_spec]:
                continue

            # Check that version exists in package.py file.
            try:
                out = run(f"spack -C {args.spack_config} versions --safe-only {abstract_spec}")
                if out.returncode != SUCCESS:
                    continue
                known_versions = set([ver(v) for v in re.split(r"\s+", out.stdout.strip())])

            except ValueError:
                # If the package doesn't exist continue on to the next commit.
                # A value error occurs because the stdout will won't
                # contain any versions.
                continue

            # Don't attempt to concretize if the version doesn't yet exist.
            is_valid = False
            spec_version = ""
            try:
                is_valid = spec.version in known_versions
                spec_version = spec.version

            # In the case that our abstract spec doesn't have a
            # version spack-python will return an error. The latest version of
            # a spec (spec with no version) will always exist and thus should
            # be considered valid.
            except SpecError:
                is_valid = True

            # If the abstract spec and version both exist at the commit we
            # should go ahead and try to concretize them.
            if is_valid:
                stdout = ""
                out = run(f"spack -C {args.spack_config} spec --yaml {abstract_spec}")

                # If concretization successful check the resulting
                # concrete specs.
                if out.returncode == SUCCESS:
                    concrete_spec = Spec.from_yaml(out.stdout)

                    # If the spec is the same move on.
                    if self.last_spec.get(abstract_spec) == concrete_spec:
                        continue
                    elif self.last_spec.get(abstract_spec) is None:
                        self.last_spec[abstract_spec] = concrete_spec
                        if args.record_first_inflection_point.lower() != "true":
                            continue

                    # If the spec is different check what is different and
                    # save as tags.
                    diff = compare_specs(
                        self.last_spec[abstract_spec],
                        concrete_spec,
                        color=False)

                    # Flatten list to tags with context.
                    tags = ["added:%s" % x for x in diff['b_not_a']]
                    tags += ["removed:%s" % x for x in diff['a_not_b']]

                    # Save concrete spec as last spec.
                    self.last_spec[abstract_spec] = concrete_spec
                    self.last_success[abstract_spec] = True

                elif self.last_success[abstract_spec] is True:
                    # Mark failing concretization points.
                    tags = ["concretization-failed"]
                    stdout = out.stderr
                    self.last_success[abstract_spec] = False

                # If spec continues to fail to concretize move on without
                # creating another inflection-point.
                else:
                    continue

                # Build list of modified files
                files = []
                for file in commit.modified_files():
                    files.append(file.change_type + ": " + file.path)

                is_primary = False
                for tag in tags:
                    if spec_version == "" \
                      and tag.startswith(f'added:version("{spec.name}"'):
                        is_primary = True
                        break
                    else:
                        if tag.startswith(f'added:version("{spec.name}","{spec_version}"'):
                            is_primary = True
                            break

                # Construct Result
                result = Result(
                    abstract_spec=abstract_spec,
                    commit=commit.hash,
                    tags=tags,
                    files=files,
                    primary=is_primary,
                    author_date=f'{commit.author_date.astimezone(tz=timezone.utc):%Y-%m-%dT%H:%M:%S+00:00}',
                    commit_date=f'{commit.committer_date.astimezone(tz=timezone.utc):%Y-%m-%dT%H:%M:%S+00:00}',
                    concretized=not contains(tags, "concretization-failed"),
                    concretizer=args.concretizer,
                    concretizer_out=stdout,
                    concrete_spec=concrete_spec)

                # Send result to drift-server
                send(result)


def send(result: Result):
    """Send reformats the input result and uploads it to the drift server."""
    # Reformat the inputted result object into the expected json.
    output = {}
    output["AbstractSpec"] = result.abstract_spec
    output["GitCommit"] = result.commit
    output["GitAuthorDate"] = result.author_date
    output["GitCommitDate"] = result.commit_date
    output["Concretizer"] = result.concretizer
    output["Files"] = result.files
    output["Tags"] = result.tags
    output["Built"] = False
    output["Concretized"] = result.concretized
    output["Primary"] = result.primary

    if result.concretizer_out is not "":
        output["ConcretizationErrUUID"] = send_artifact(result.concretizer_out)

    output["SpecUUID"] = send_artifact(result.concrete_spec, datatype="application/json")

    # Upload json data to the drift server.
    r = requests.put(
        f"{args.host}/inflection-point",
        json=output,
        auth=requests.auth.HTTPBasicAuth(args.username, args.password),
    )
    print(json.dumps(output), flush=True)
    print(r.status_code)

def send_artifact(artifact: str, datatype="text/plain"):
    r = requests.put(
        f"{args.host}/artifact",
        headers={'Content-Type':datatype},
        body=artifact,
        auth=requests.auth.HTTPBasicAuth(args.username, args.password),
    )

    if r.status_code is not 200:
        print("Error uploading artifact")
        print(f"{args.host}/artifact: {r.status_code}")
        exit()

    resp = r.getresponse().json()
    return resp["uuid"]

def run(command):
    """Run executes a spack command using the known spack bin."""
    result = subprocess.run(
        [args.repo + "/bin/"+command.split()[0]] + command.split()[1:],
        stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True
    )
    return result


def main():
    """Set up the drift-analysis runner."""
    # Setup Argument Parsing
    parser = ArgumentParser()
    parser.add_argument("--host")
    parser.add_argument("--username")
    parser.add_argument("--password")
    parser.add_argument("--repo")
    parser.add_argument("--spack-config")
    parser.add_argument("--since-commit")
    parser.add_argument("--since")
    parser.add_argument("--to-commit")
    parser.add_argument("--to")
    parser.add_argument("--specs")
    parser.add_argument("--concretizer")
    parser.add_argument("--record-first-inflection-point")
    parser.add_argument("--skip-known-inflection-points")

    # Parse Arguments
    global args
    args = parser.parse_args()

    # Define program parameter defaults
    since = None
    to = None
    since_commit = None
    to_commit = None
    if args.since is not None:
        data = args.since.split("/")
        since = datetime(int(data[2]), int(data[0]), int(data[1]))
    if args.to is not None:
        data = args.to.split("/")
        to = datetime(int(data[2]), int(data[0]), int(data[1]))
    if args.since_commit is not None:
        since_commit = args.since_commit
    if args.to_commit is not None:
        to_commit = args.to_commit
    if args.record_first_inflection_point is None:
        args.record_first_inflection_point = "true"
    if args.skip_known_inflection_points is None:
        args.skip_known_inflection_points = "true"

    # If specs are actually defined split them up based on commas.
    specs = None
    if args.specs is not None:
        specs = args.specs.split(",")

    # Create new Drift Analysis
    da = DriftAnalysis(specs)
    da.traverse(
        args.repo,
        since=since,
        to=to,
        from_commit=since_commit,
        to_commit=to_commit,
        checkout=True,
        printTrial=True,
    )


if __name__ == "__main__":
    main()
