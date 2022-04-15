#!/bin/python3

"""Calculate which runners timed out and need to be re-run."""

import sys

pkgs = []

for line in sys.stdin:
    line = line.split(",")[0].split(".")[1]
    pkgs.append(line)

sys.stdout.write("["+", ".join(pkgs)+"]\n")
