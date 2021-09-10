import json
from argparse import ArgumentParser
from dataclasses import dataclass
from datetime import datetime
from subprocess import run as subprocess_run

import requests
import re
from helicase import Helicase
from spack.cmd.diff import compare_specs
from spack.error import SpecError
from spack.spec import Spec
from spack.version import ver

# Define SUCCESS for comparing command return codes.
SUCCESS = 0

@dataclass
class Result:
    name:str
    version:str
    commit:str
    tags:list
    timestamp:str

class DriftAnalysis(Helicase):
    def __init__(self, specs=None):
        """Setup the last cache and specs used for the analysis."""

        self.last = {}
        self.specs = specs or []

    def analyze(self, commit):
        """Analyze defines the process that should be run on every commit
           and for every abstract spec. This is used to determine if the 
           commit is an "inflection point". Inflection points occur when
           the concretization for an abstract spec changes. This most 
           often happens when a package updates or one of its dependencies
           update leading to a different installation for the same input
           abstract spec."""

        for abstract_spec in self.specs:
            # Create spec for parsing info
            spec = Spec(abstract_spec)

            # Check that version exists in package.py file.
            out = run(f"spack versions --safe-only {abstract_spec}")
            known_versions = set([ver(v) for v in re.split(r"\s+", out.stdout.strip())])

            # Don't attempt to concretize if the version doesn't yet exist.
            is_valid = False
            spec_version = ""
            try:
                is_valid = spec.version in known_versions
                spec_version = spec.version
            except SpecError:
                is_valid = True

            if is_valid:
                out = run(f"spack spec --yaml {abstract_spec}")

                # If concretization successful check the resulting concrete specs.
                if out.returncode == SUCCESS:
                    concrete_spec = Spec.from_yaml(out.stdout)

                    # If the spec is the same move on.
                    if self.last.get(abstract_spec) == concrete_spec:
                        continue
                    else:
                        self.last[abstract_spec] = concrete_spec

                    # If the spec is different check what is different and save as tags.
                    diff = compare_specs(
                        self.last[abstract_spec],
                        concrete_spec,
                        color=False)

                    # Flatten list to tags with context.
                    tags = ["added:%s:%s" %(x[0],x[1]) for x in diff['b_not_a']]
                    tags += ["removed:%s:%s" %(x[0],x[1]) for x in diff['a_not_b']]

                    # Save concrete spec as last spec.
                    self.last[abstract_spec] = concrete_spec

                else:
                    # Mark failing concretization points.
                    tags = ["concretization-failed"]

                # Construct Result
                result = Result(
                    spec.name,
                    spec_version,
                    commit.hash,
                    tags,
                    str(commit.author_date))

                # Send result to drift-server
                send(result)
    
def send(result:Result):
    """Send reformats the input result and attempts to upload the data as
       as a json message to the supplied drift server."""

    # Reformat the inputted result object into the expected json.
    output = {}
    output["commit"] = {"digest":result.commit, "timestamp":result.timestamp}
    output["tags"] = [{"name": x} for x in result.tags]
    output["package"] = {"name": result.name, "version": result.version}

    # Upload json data to the drift server.
    # r = requests.post(f"{args.host}/inflection-point/", json=output, auth=requests.auth.HTTPBasicAuth(args.username, args.password)) 
    print(json.dumps(output), flush=True)
    # print(r.status_code)

def run(command):
    result = subprocess_run([args.repo + "/bin/"+command.split()[0]] + command.split()[1:], 
        capture_output=True, text=True)
    return result

def main():
    # Setup Argument Parsing
    parser = ArgumentParser()
    parser.add_argument("--host")
    parser.add_argument("--username")
    parser.add_argument("--password")
    parser.add_argument("--repo")
    parser.add_argument("--since-commit")
    parser.add_argument("--since")
    parser.add_argument("--to-commit")
    parser.add_argument("--to")
    parser.add_argument("--specs")

    # Parse Arguments
    global args 
    args = parser.parse_args()

    # Define program parameter defaults
    since = None
    to = None
    since_commit = None
    to_commit = None
    if args.since != None:
        data = args.since.split("/")
        since = datetime(int(data[2]), int(data[0]), int(data[1]))
    if args.to != None:
        data = args.since.split("/")
        to = datetime(int(data[2]), int(data[0]), int(data[1])) 
    if args.since_commit != None:
        since_commit = args.since_commit
    if args.to_commit != None:
        to_commit = args.to_commit
    
    # If specs are actually defined split them up based on commas.
    specs = None
    if args.specs != None:
        specs = args.specs.split(",")

    # Create new Drift Analysis
    da = DriftAnalysis(specs)
    da.traverse(args.repo, 
        since=since,
        to=to,
        from_commit=since_commit,
        to_commit=to_commit,
        checkout=True,
        printTrial=True,
    )

if __name__ == "__main__":
    main()
