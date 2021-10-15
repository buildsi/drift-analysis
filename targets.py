import json
import requests

def main():
    r = requests.get("https://autamus.io/registry/library.json")

    print("["+", ".join([x["name"] for x in r.json()])+"]")


if __name__ == "__main__":
    main()
