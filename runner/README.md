# runner
## Table of Contents
- [Introduction](#introduction)
- [Getting Started](#getting-started)
- [Usage](#usage)

## Introduction
The core drift analysis tool system is broken into the [server](https://github.com/buildsi/drift-analysis/tree/main/server) and this runner. Where the server is meant to be a long lived persistant application, this runner meant to be run as a massively parallel set of batch jobs on an HPC cluster. Using the included Maestro job definition each runner can be spawned off as their own SLRUM job, each tackling their own abstract specs.

The drift analysis runner is a python based application that scrapes through a specifies interval of the git history of Spack. The runner then attempts to concretize each abstract spec that it is in charge of at **every** commit and records instances where either the concretization through Spack fails or the concrete spec changes given the same abstract spec after a given commit. The runner will then attempt to upload these inflection points to the specified drift analysis server.

## Dependencies
- [Helicase](https://github.com/buildsi/helicase)
- [Maestro](https://github.com/llnl/maestrowf) (Optional)

## Gettings Started
##### Running With Maestro
1. Run `git clone git@github.com:buildsi/drift-analysis.git` to download the drift-analysis git repository.
2. Open the runner directory by running `cd drift-analysis/runner`.
3. Edit the Maestro Analysis by running `$EDITOR analysis.yml` and replace the defaults with your values.
4. Launch the Maestro study by running `maestro run analysis.yml`.

#### Running Without Maestro
1. Run `git clone git@github.com:buildsi/drift-analysis.git` to download the drift-analysis git repository.
2. Open the runner directory by running `cd drift-analysis/runner`.
3. See the following usage section on how to execute the runner.



## Usage
```bash

# To improve execution time you can save the spack installation and
# caches to a RAMDisk.
mkdir -p /dev/shm/spack-$(SPEC)/spack
mkdir -p /dev/shm/spack-$(SPEC)/cache
mkdir -p /dev/shm/spack-$(SPEC)/config

# Create a custom runner overlay Spack config to modify the cache location and
# concretizer used as needed
echo "config:" > /dev/shm/spack-$(SPEC)/config/config.yaml
echo "    misc_cache: /dev/shm/spack-$(SPEC)/cache" >> /dev/shm/spack-$(SPEC)/config/config.yaml
echo "    concretizer: $(CONCRETIZER)" >> /dev/shm/spack-$(SPEC)/config/config.yaml

# Make another clone of Spack
# Note: this installation will be the one that we'll use for the git history
#       and concretizations at each point in the git history but you'll
#       need another version of Spack to execute and run the actual runner.

git clone https://github.com/spack/spack.git /dev/shm/spack-$(SPEC)/spack

# You should update this path to be where you saved the drift-analysis
# git repository earlier.

cd ~/projects/drift-analysis/runner/

# Execute the runner using spack-python to automatically
# import Spack as a dependency of the runner.
spack-python runner.py \
        --host="" \      # <-- this should be the URL to your server instance
        --username="" \  # <-- this should be the same as your server api username
        --password="" \  # <-- this should be the same as your server api password
        --repo=/dev/shm/spack-$(SPEC)/spack \
        --spack-config=/dev/shm/spack-$(SPEC)/config \
        --from-commit="" \  # <-- the commit hash of where to start in Spack's git history
        --to-commit="" \    # <-- the commit hash of where to stop in Spack's git history
        --spec="$(SPEC)"
        --concretizer="$(CONCRETIZER)"

# Clean up the directories we've created in the RAMDisk.
rm -rf /dev/shm/spack-$(SPEC)
```
