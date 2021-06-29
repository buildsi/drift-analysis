from helicase import Helicase
from datetime import *
import json
import subprocess, sys

# Define SUCCESS for comparing command return codes.
SUCCESS = 0

def spack(command):
    result = subprocess.run([sys.argv[1] + "/bin/spack"] + command.split(), 
        capture_output=True, text=True)
    return result

class Result:
    commit:str
    tags:list
    spec:str

    def __init__(self):
        self.commit = ""
        self.tags = []
        self.spec = ""


class DriftAnalysis(Helicase):
    specs = {}
    last = {}
    def __init__(self, specs):
        for spec in specs:
            self.specs[spec] = []
            self.last[spec] = ""

    def analyze(self, commit):
        for spec in self.specs:
            out = spack("spec --json " + spec)
            result = Result()
            result.commit = commit.hash
            if out.returncode == SUCCESS:
                concrete_spec = json.loads(out.stdout)["spec"]
                # In the case that the concretized spec hash changes
                # we want to record the commit where this happens.
                if self.last[spec] != "" and self.last[spec][0][spec]["full_hash"] != concrete_spec[0][spec]["full_hash"]:
                    # Check if the package version is different.
                    if self.last[spec][0][spec]["version"] != concrete_spec[0][spec]["version"]:
                        result.tags += ["pkg-updated"]

                    # +---------------------+
                    # |     Dependencies    |
                    # +---------------------+
                    # Build dependency cache map for the last spec.
                    cache_map = {}
                    for i in range(len(self.last[spec])):
                        for name, dep in self.last[spec][i].items():
                            cache_map[name] = dep

                    # Iterate through all of the dependencies in the new spec.
                    for i in range(len(concrete_spec)):
                        for name, dep in concrete_spec[i].items():
                            # Check if a dependency version is different.
                            if dep["full_hash"] != cache_map[name]["full_hash"] \
                                and name != spec and "dep-updated" not in result.tags:
                                result.tags += ["dep-updated"]
                            # Check if a dependency has been added.
                            if name not in cache_map and "dep-added" not in result.tags:
                                result.tags += ["dep-added"]
                            # Remove dependency from known dependencies to track if we don't hit any
                            # dependencies in the new version that use to be in the old spec.
                            cache_map.pop(name)
                    
                    # If there are dependencies that we didn't hit in the new version mark
                    # the commit has having lost dependencies.
                    if len(cache_map) > 0:
                        result.tags += ["dep-deleted"]

                    # +---------------------+
                    # |      Variants       |
                    # +---------------------+
                    if "parameters" in concrete_spec[0][spec]:
                        # Check for variants added or modified by comparing new spec to old.
                        for param, value in concrete_spec[0][spec]["parameters"].items():
                            if param not in self.last[spec][0][spec]["parameters"]:
                                if param == "patches" and "patches-added" not in result.tags:
                                    result.tags += ["patches-added"]
                                elif "variant-added" not in result.tags:
                                    result.tags += ["variant-added"]
                            if value != self.last[spec][0][spec]["parameters"][param]:
                                if param == "patches" and "patches-modified" not in result.tags:
                                    result.tags += ["patches-modified"]
                                elif "variant-modified" not in result.tags:
                                    result.tags += ["variant-modified"]

                        # Compare old to new to see if any variants were removed.
                        for param, value in self.last[spec][0][spec]["parameters"].items():
                            if param not in concrete_spec[0][spec]["parameters"]:
                                if param == "patches" and "patches-removed" not in result.tags:
                                    result.tags += ["patches-removed"]
                                elif "variant-removed" not in result.tags:
                                    result.tags += ["variant-removed"]


                    # Add commit to list of inflection point commits.
                    self.specs[spec] += [result]

                # Save concrete spec hash as last hash.
                self.last[spec] = concrete_spec
            else:
                # If the spec doesn't concretize properly we also want
                # to record the commit at which this occurred.
                self.specs[spec] += [result]

def main():
    dt = datetime(2021, 5, 1)
    now = datetime.now()

    da = DriftAnalysis([sys.argv[2]])
    da.traverse(sys.argv[1], since=dt, to=now, checkout=True, printTrial=True)
    output = {}
    for spec in da.specs:
        output[spec] = {}
        for result in da.specs[spec]:
            output[spec][result.commit] = result.tags

    print(json.dumps(output))

main()