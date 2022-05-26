# drift-analysis
A System of Tools to Measure Package Drift Within the Spack Repository

## Table of Contents
- [Introduction](#introduction)
  - [What is Package Drift?](#what-is-package-drift)
  - [How Can we Measure & Observe Package Drift?](#how-can-we-measure-observe-package-drift)
- [Usage](#usage)

## Introduction
As research software becomes more complex to support novel scientific workflows, it also becomes more difficult to install and maintain. For although shared programming libraries have allowed for rapid software development and reduced code duplication, they’ve also created complex and fragile dependency relationships between applications. To keep track of these dependency relationships and automate the process of installing complex software stacks, package managers were developed. However, current package managers still rely on these inter package relationships to be defined by human developers which are often imprecise and time-consuming. Additionally, as packages change over time, they often unknowingly break the ways other software uses or builds off of them. This process of package relationships changing over time is what we call “package drift” and it can have lasting impacts on scientific reproducibility. To better understand package drift, we designed and implemented this system to track the effects of package updates across the Spack package ecosystem and record if a change to a package resulted in a broken software stack. In the future, we hope that this data will be used to create more intelligent package managers that may automatically learn from failed installations and optimize for the most likely configuration that will install successfully.

### What is Package Drift?
To install the complex software stacks talked about before, Spack — as well as other package managers — use databases to determine which dependency packages must also be installed for a given application. Within Spack this process is known as going from an abstract specification to a concrete specification ( “spec”). One consequence of this process however, is that “install pkg” means different things as updates to a package manager arrive. You might say “install pkg” today and get the following installations:

```
pkg@1.1, dep_a@1.0, dep_b@1.0, dep_c@2.0
```

Here, “pkg” is an abstract spec of what is to be installed, and we call these installations and their versions a concrete spec.  If you asked for “install pkg” tomorrow, you’d be providing the same abstract spec but you may see that a different concrete spec with newer versions is installed:

```
pkg@1.2, dep_a@1.1, dep_b@2.0, dep_c@3.0
```

This process of a concrete spec changing over time given the same abstract spec is what we call package drift.

### How Can we Measure & Observe Package Drift?
Since package drift occurs when the same abstract specification results in a different concrete specification, the package drift with a package manager can be recorded by marking the points over time when a package’s concrete specification changes. Here [Spack](https://github.com/spack/spack) presents a unique opportunity, since the Spack package database utilizes Git version control it is possible to iterate through checkpoints in the version control system to simulate the process of traveling through time. At each checkpoint an abstract spec is converted into a concrete spec and when a given concrete spec differs from the previous version that checkpoint is marked as an inflection point within the database. In addition to marking an inflection point, we use Spack to tag each point with the differences from one concrete spec to the next. One of these tags in specific includes which dependency packages have been updated so that in the future drift-aware package managers may be able to reason about the best set of dependencies to use to compile a working application. To traverse the checkpoints within a Git repository, we developed a library for this project called [Helicase](https://github.com/buildsi/helicase). Inspired by the protein that assists with DNA replication, our library helps developers unzip a Git repository’s history to explore changes over time. Written in Python, Helicase allows a user to specify start and end commits as well as a custom analysis function to be run on every commit within the repository. Within the drift analysis project, this custom function compared abstract specs and their concrete counterparts over time, recording and uploading to the database when differences occurred.


## Usage
### Exploring Existing Data
#### Drift Analysis Visualization
If you would like to explore the existing drift analysis data we've collected on the Spack repository in a graphical form, check out the visualization dashboard hosted at [buildsi.github.io/drift-analysis](https://buildsi.github.io/drift-analysis).

#### Accessing the Raw Data
If you would like to build your own visualizations or incorperate our data into your own analysis all of the data we've collected so far is hosted at [drift-server.spack.io/inflection-points](https://drift-server.spack.io/inflection-points). For more information on how to interact with the Drift Server check out the documentation [here](https://github.com/buildsi/drift-analysis/tree/main/server#readme).

### Running A Custom Analysis
If you'd like to setup your own drift-analysis study checkout the documentation [here](https://github.com/buildsi/drift-analysis/tree/main/server#readme) to setup your own drift-server and [here](https://github.com/buildsi/drift-analysis/tree/main/runner) to learn how to execute the drift-analysis runner application.
