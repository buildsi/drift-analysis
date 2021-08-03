from argparse import ArgumentParser
from datetime import datetime
from helicase import Helicase
from requests import post
from spack import *
import spack.cmd.diff
from subprocess import run as subprocess_run

# Define SUCCESS for comparing command return codes.
SUCCESS = 0

class Result:
    def __init__(self, name, version, commit, tags, timestamp):
        self.name = ""
        self.commit = ""
        self.tags = []
        self.timestamp = ""
        self.version = ""


class DriftAnalysis(Helicase):
    def __init__(self, specs=None):
        self.last = {}
        self.specs = specs or []

    def analyze(self, commit):
        for abstract_spec in self.specs:
            tags = []
            # Separate out spec version and name.
            name = abstract_spec
            version = None
            if "@" in abstract_spec:
                name, version = abstract_spec.split("@", 1)
            # Check that version exists in package.py file.
            out = run(f"spack versions --safe-only {abstract_spec}")
            # Don't attempt to concretize if the version doesn't yet exist.
            if abstract_spec[verToken+1:] in out.stdout:
                out = run(f"spack spec --yaml {abstract_spec}")
                # If concretization successful check the resulting concrete specs.
                if out.returncode == SUCCESS:
                    concrete_spec = spack.spec.Spec().from_yaml(out.stdout)
                    # If the spec is the same move on.
                    if abstract_spec not in self.last or self.last[abstract_spec] == concrete_spec:
                        continue
                    # If the spec is different check what is different and save as tags.
                    diff = spack.cmd.diff.compare_specs(
                        self.last[abstract_spec],
                        concrete_spec,
                        colorful=False)
                    # Flatten list to tags with context.
                    tags = ["%s:%s" %(x[0],x[1]) for x in diff['b_not_a']]
                    # Save concrete spec as last spec.
                    self.last[abstract_spec] = concrete_spec
                # Mark failing concretization points.
                else:
                    tags = [("concretization-failed","")]
                # Construct Result
                result = Result(
                    name,
                    version,
                    commit.hash,
                    tags,
                    str(commit.author_date))
                # Send result to drift-server
                send(result)
    
def send(result:Result):
    output = {}
    output["commit"] = {"digest":result.commit, "timestamp":result.timestamp}
    output["tags"] = [{"name": x} for x in tags]
    output["package"] = {"name": result.name, "version": result.version}

    r = requests.post(f"{args.host}/inflection-point/", json=output, auth=requests.auth.HTTPBasicAuth(args.username, args.password)) 
    print(json.dumps(output), flush=True)
    print(r.status_code)

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

    specs = None
    if args.specs != None:
        specs = args.specs.split()
    # Create new Drift Analysis
    da = DriftAnalysis(specs)
    da.traverse(args.repo, since=since, to=to, from_commit=since_commit, to_commit=to_commit, checkout=True, printTrial=True)

if __name__ == "__main__":
    main()