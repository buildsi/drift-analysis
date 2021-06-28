from helicase.helicase import Helicase
from datetime import *
import json
import subprocess, sys

# Define SUCCESS for comparing command return codes.
SUCCESS = 0

def spack(command):
    result = subprocess.run([sys.argv[1] + "/bin/spack"] + command.split(), 
        capture_output=True, text=True)
    return result

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
            if out.returncode == SUCCESS:
                concrete_spec = json.loads(out.stdout)["spec"][0][spec]
                if self.last[spec] != "" and self.last[spec] != concrete_spec["full_hash"]:
                    # In the case that the concretized spec hash changes
                    # we want to record the commit where this happens.
                    self.specs[spec] += [commit.hash]
                self.last[spec] = concrete_spec["full_hash"]
            else:
                # If the spec doesn't concretize properly we also want
                # to record the commit at which this occurred.
                self.specs[spec] += [commit.hash]

def main():
    dt = datetime(2021, 1, 1)
    now = datetime.now()

    da = DriftAnalysis(["abyss"])
    da.traverse(sys.argv[1], since=dt, to=now, checkout=True)
    print(json.dumps(da.specs))

main()