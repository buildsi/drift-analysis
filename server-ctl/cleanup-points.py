from argparse import ArgumentParser

import dateutil.parser
import requests


def main():
    parser = ArgumentParser(description='Drift Server Controller')

    parser.add_argument("--username", type=str)
    parser.add_argument("--password", type=str)
    parser.add_argument('host', type=str)

    parser.add_argument('date', type=str,
                        help='Remove everything before this date')

    args = parser.parse_args()

    # date = dateutil.parser.parse(args.date)

    r = requests.get(f"{args.host}/inflection-points")
    for point in r.json():
        # point_date = dateutil.parser.parse(point["GitCommitDate"])
        if point["Tags"] == []:
            output = {}
            output["commit"] = {"digest": point["GitCommit"]}
            output["concretizer"] = point["Concretizer"]
            output["abstract_spec"] = point["AbstractSpec"]

            pr = requests.post(f"{args.host}/delete/inflection-point",
                               json=output,
                               auth=requests.auth.HTTPBasicAuth(
                                   args.username,
                                   args.password))

            print(f"Deleting: {point['ID']}: {pr.status_code}")


if __name__ == "__main__":
    main()
