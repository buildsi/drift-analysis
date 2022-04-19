"""This program implements a minimal drift-analysis builder."""

import json
import subprocess
from argparse import ArgumentParser

import requests
from helicase import Repo

# Define SUCCESS for comparing command return codes.
SUCCESS = 0


def get_inflection_points(host: str, abstract_spec: str):
    """Return a list of inflection points represented as dictionaries."""
    r = requests.get(f"{host}/inflection-points/{abstract_spec}")
    return r.json()


def run(command):
    """Run executes a spack command using the known spack bin."""
    result = subprocess.run(
        [args.repo + "/bin/"+command.split()[0]] + command.split()[1:],
        stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True
    )
    return result


def send(point: dict):
    """Upload an updated inflection point to the server."""
    # Upload json data to the drift server.
    r = requests.put(
        f"{args.host}/inflection-point",
        json=point,
        auth=requests.auth.HTTPBasicAuth(args.username, args.password),
    )
    print(json.dumps(point), flush=True)
    print(r.status_code)


def send_artifact(artifact: str, datatype="text/plain"):
    """Upload an artifact to a drift-server and return its UUID."""
    r = requests.put(
        f"{args.host}/artifact",
        headers={'Content-Type': datatype},
        body=artifact,
        auth=requests.auth.HTTPBasicAuth(args.username, args.password),
    )

    if r.status_code != 200:
        print("Error uploading artifact")
        print(f"{args.host}/artifact: {r.status_code}")
        exit()

    resp = r.getresponse().json()
    return resp["uuid"]


def main():
    """Pull inflection points and attempt to build them."""
    # Setup Argument Parsing
    parser = ArgumentParser()
    parser.add_argument("--host")
    parser.add_argument("--username")
    parser.add_argument("--password")
    parser.add_argument("--repo")
    parser.add_argument("--spack-config")
    parser.add_argument("--specs")
    parser.add_argument("--concretizer")

    # Parse Arguments
    global args
    args = parser.parse_args()

    # Define git repository
    repo = Repo(args.repo)

    # Repeat the build process for all assigned specs.
    for abstract_spec in args.specs:

        # Get a list of known inflection points from a drift server instance.
        inflection_points = get_inflection_points(args.host, abstract_spec)
        for point in inflection_points:
            if point["Concretizer"] == args.concretizer \
              and point["Concretized"] and point["BuildLogUUID"] == "":
                print(f"[BUILDING] {abstract_spec} at {point['GitCommit']}")

                # Checkout inflection point git commit in repository.
                repo.git("checkout", point['GitCommit'])

                # Attempt to build the abstract_spec at the inflection_point
                out = run(f"spack -C {args.spack_config} install --fail-fast {abstract_spec}")

                # Upload build log to a drift-server instance.
                point["BuildLogUUID"] = send_artifact(out.stdout)

                # Update "Build" status for point
                point["Built"] = out.returncode == SUCCESS

                # Upload updated point to a drift-server instance.
                send(point)


if __name__ == "__main__":
    main()
