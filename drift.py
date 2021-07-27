from helicase import Helicase
from datetime import *
import argparse
import json, requests
import subprocess, sys

# Define SUCCESS for comparing command return codes.
SUCCESS = 0
args = None

class Result:
    abstract:str
    commit:str
    tags:list
    timestamp:str
    spec:str
    version:str

    def __init__(self):
        self.abstract = ""
        self.commit = ""
        self.tags = []
        self.timestamp = ""
        self.spec = ""
        self.version = ""


class DriftAnalysis(Helicase):
    specs = []
    last = {}
    def __init__(self, specs):
        self.specs = specs
        for spec in specs:
            self.last[spec] = ""

    def analyze(self, commit):
        for spec in self.specs:
            out = spack("spec --json " + spec)
            verToken = spec.find("@")
            if verToken < 0: verToken = len(spec)
            name = spec[:verToken]
            result = Result()
            result.commit = commit.hash
            result.abstract = name
            result.timestamp = str(commit.author_date)
            if len(spec[verToken:]) > 0:
                result.version = spec[verToken+1:]
            if out.returncode == SUCCESS:
                concrete_spec = json.loads(out.stdout)["spec"]
                # In the case that the concretized spec hash changes
                # we want to record the commit where this happens.
                if self.last[spec] != "" and self.last[spec][0][name]["full_hash"] != concrete_spec[0][name]["full_hash"]:
                    # Check if the package version is different.
                    if self.last[spec][0][name]["version"] != concrete_spec[0][name]["version"]:
                        result.tags += ["pkg-updated"]
                    # +---------------------+
                    # |     Dependencies    |
                    # +---------------------+
                    result.tags += set(self.tag_deps(concrete_spec, spec, name))
                    # +---------------------+
                    # |      Variants       |
                    # +---------------------+
                    result.tags += set(self.tag_variants(concrete_spec, spec, name))
                    # Send Inflection Point
                    send(result)
                # Save concrete spec hash as last hash.
                self.last[spec] = concrete_spec
            else:
                # If the spec doesn't concretize properly we also want
                # to record the commit at which this occurred.
                result.tags += ["concretizaton-failed"]
                send(result)

    def tag_deps(self, concrete_spec, spec, name):
        tags = []
        # Build dependency cache map for the last spec.
        cache_map = {}
        for i in range(len(self.last[spec])):
            for name, dep in self.last[spec][i].items():
                cache_map[name] = dep
        # Iterate through all of the dependencies in the new spec.
        for i in range(len(concrete_spec)):
            for name, dep in concrete_spec[i].items():
                # Check if a dependency has been added.
                if name not in cache_map:
                    tags += ["dep-added"]
                    continue
                # Check if a dependency version is different.
                elif dep["full_hash"] != cache_map[name]["full_hash"] and name != spec:
                    tags += ["dep-updated"]
                # Remove dependency from known dependencies to track if we don't hit any
                # dependencies in the new version that use to be in the old spec.
                cache_map.pop(name)
        # If there are dependencies that we didn't hit in the new version mark
        # the commit has having lost dependencies.
        if len(cache_map) > 0:
            tags += ["dep-deleted"]
        return tags

    def tag_variants(self, concrete_spec, spec, name):
        tags = []
        if "parameters" in concrete_spec[0][name]:
            # Check for variants added or modified by comparing new spec to old.
            for param, value in concrete_spec[0][name]["parameters"].items():
                if param not in self.last[spec][0][name]["parameters"]:
                    if param == "patches":
                        tags += ["patches-added"]
                    else:
                        tags += ["variant-added"]
                elif value != self.last[spec][0][name]["parameters"][param]:
                    if param == "patches":
                        tags += ["patches-modified"]
                    else:
                        tags += ["variant-modified"]
            # Compare old to new to see if any variants were removed.
            for param, value in self.last[spec][0][name]["parameters"].items():
                if param not in concrete_spec[0][name]["parameters"]:
                    if param == "patches":
                        tags += ["patches-removed"]
                    else:
                        tags += ["variant-removed"]
            return tags
    
def send(result:Result):
    output = {}
    output["commit"] = {"digest":result.commit, "timestamp":result.timestamp}
    output["tags"] = []
    for tag in result.tags:
        output["tags"] += [{"name": tag}]
    output["package"] = {"name": result.abstract, "version": result.version}

    r = requests.post(f"{args.host}/inflection-point/", json=output, auth=requests.auth.HTTPBasicAuth(args.username, args.password)) 
    print(json.dumps(output), flush=True)
    print(r.status_code)

def spack(command):
    result = subprocess.run([args.repo + "/bin/spack"] + command.split(), 
        capture_output=True, text=True)
    return result

def main():
    # Setup Argument Parsing
    parser = argparse.ArgumentParser()
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
    # Define program parameters
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

    da = DriftAnalysis([args.specs])
    da.traverse(args.repo, since=since, to=to, from_commit=since_commit, to_commit=to_commit, checkout=True, printTrial=True)

if __name__ == "__main__":
    main()