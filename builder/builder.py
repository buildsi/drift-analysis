import json
import requests
import subprocess
from argparse import ArgumentParser
from spack.spec import Spec
from git import Repo

# Define SUCCESS for comparing command return codes.
SUCCESS = 0

def get_inflection_points(host:str, abstract_spec:str):
    r = requests.get(f"{host}/inflection-points/{abstract_spec}")
    return r.json()

def build(abstract_spec:str):
    result = subprocess.run(["spack","build", abstract_spec], 
        capture_output=True, text=True)
    return result.returncode

def concretize_spec(abstract_spec:str):
    result = subprocess.run(["spack","spec", "--json", abstract_spec],
        capture_output=True, text=True)
    return result

def upload_build(point_id:int, spec_id:int, status:str):
    build = {"point_id":point_id, "spec_id":spec_id, "status":status}
    r = requests.post("http://localhost:8080/build/", json=build)
    return r.json()

def upload_spec(name:str, version:str, concrete_spec:str):
    spec = {"package": {"name":name, "version":version}, "data": json.dumps(concrete_spec)}
    r = requests.post("http://localhost:8080/spec/", json=spec)
    return r.json()["ID"]

def main():
    # Setup Argument Parsing
    parser = ArgumentParser()
    parser.add_argument("--host")
    parser.add_argument("--username")
    parser.add_argument("--password")
    parser.add_argument("--repo")
    parser.add_argument("--spack-config")
    parser.add_argument("--spec")
    parser.add_argument("--concretizer")

    # Parse Arguments
    global args 
    args = parser.parse_args()

    # Define git repository
    repo = Repo(args.repo)

    # Get abstract spec from inputs
    abstract_spec = args.spec

    # Create spec based on parsed info
    spec = Spec(abstract_spec)

    # Grab a list of inflection points from the drift server.
    inflection_points = get_inflection_points(args.host, abstract_spec)
    for point in inflection_points:
        if point["Concretizer"] == args.concretizer \
            and point["Concretized"] and point["Spec"] == "":

            print(f"Building {abstract_spec} at {point["GitCommit"]})
        # # Checkout Commit in Git
        # # repo.git.checkout(commit["Commit"]["Digest"])
        # # Concretize the Spec
        # out = concretize_spec(abstract_spec)
        # if out.returncode == SUCCESS:
        #     concrete_spec = json.loads(out.stdout)["spec"]
        #     spec_id = upload_spec(name, version, concrete_spec)
        #     # Attempt Package Build
        #     built = build(abstract_spec)
        #     if built == SUCCESS:
        #         upload_build(commit["ID"], spec_id, "success")
        #     else:
        #         upload_build(commit["ID"], spec_id, "failed")
        # else:
        #     upload_build(commit["ID"], -1, "failed")

if __name__ == "__main__":
    main()
