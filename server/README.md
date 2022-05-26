# server

The drift server is a simple api and database to store results from a
drift analysis study.

## Table of Contents
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [Environment](#environment)
- [Endpoints](#endpoints)

## Getting Started
To make deployment as easy as possible the drift server is automatically built and distributed as a OCI container from GHCR. To run a local instance you can run the following (replacing your environment variables for the ones listed here),
```
podman run --name drift-server \
  -e DRIFTSERVER_API_USERNAME=myusername \
  -e DRIFTSERVER_API_PASSWORD=mypassword \
  -e DRIFTSERVER_API_PORT=8080 \
  -e DRIFTSERVER_DB_PATH=/data/drift-server.db
  -e DRIFTSERVER_S3_BUCKET=my-awesome-s3-bucket \
  -e DRIFTSERVER_S3_ACCESSKEY=reDAcTed \
  -e DRIFTSERVER_S3_SECRETKEY=reDAcTed \
  -v /path/to/drift-server.db:/data/drift-server.db \
  -p 8080:8080 \
  ghcr.io/buildsi/drift-server \
```

## Authentication

By default, authentication is disabled, meaning that the username and password
are empty strings. To change this, export the following variables when you
run the server:

```bash
export DRIFTSERVER_API_USERNAME=myusername
export DRIFTSERVER_API_PASSWORD=mypassword
```

For any endpoint described below, you can add basic authentication to curl as follows:

```bash
$ curl -u $DRIFTSERVER_API_USERNAME:$DRIFTSERVER_API_PASSWORD ...
```

## Environment

By default, the server will run on port 8080 and use an sqlite database. However, you
can set the following environment variables to change that. Notice that the username
and password shown above are part of the environment variable set. This is because
they are part of the same server config.

| variable | description | default |
|-----------|------------|----------|
| DRIFTSERVER_API_USERNAME | The username to authenticate with basic http auth | "" |
| DRIFTSERVER_API_PASSWORD | The password to authenticate with basic http auth | "" |
| DRIFTSERVER_API_PORT | The port to run the server on | "8080" |
| DRIFTSERVER_DB_PATH | The path to the drift server sqlite database.| "points.db" |
| DRIFTSERVER_S3_BUCKET | The name of the s3 bucket to use to store artifacts. | "" |
| DRIFTSERVER_S3_ACCESSKEY | The s3 access key for the configured bucket. | "" |
| DRIFTSERVER_S3_SECRETKEY | The s3 secret key for the configured bucket. | "" |
| DRIFTSERVER_S3_ENDPOINT | The endpoint location for the s3 server. | "" |
| DRIFTSERVER_S3_REGION | The region for the s3 server. | "us-east-1" |


## Inflection Point
Below is an example of an inflection point that can be uploaded or downloaded from the drift server.
```json
{
  "ID":200,
  "AbstractSpec":"openmpi",
  "GitCommit":"f28ca41d02ce9af1b83c16b1d9a0d00ab4a2ad12",
  "GitAuthorDate":"2021-02-26T09:19:32Z",
  "GitCommitDate":"2021-02-26T09:19:32Z",
  "Concretizer":"original",
  "Files": [
    "M: var/spack/repos/builtin/packages/hwloc/package.py"
  ],
  "Tags":[
    "added:hash(\"hwloc\", \"h4yjvwl3gawsi35nwaf3dnfqdjyx34vf\")",
    "added:hash(\"openmpi\", \"sxm226crsf5t5p65mnc3bgw5oqs7tfru\")",
    "added:version(\"hwloc\", \"2.4.1\")","removed:hash(\"hwloc\", \"glezzhsmy6vfdyf67ak7iu3yyuzrifrw\")",
    "removed:hash(\"openmpi\", \"pnqtwi5mja4nn3juxtcdu7gwkg7zegw2\")",
    "removed:version(\"hwloc\", \"2.4.0\")"
  ],
  "SpecUUID":"83a67604-a37b-4ae8-a1ca-b38d67d42db1",
  "Built":true,
  "Concretized":true,
  "Primary":false,
  "BuildLogUUID":"ffba49d3-d4d4-48e5-bf95-c61b378747d4",
  "ConcretizationErrUUID":""
}
```

## Endpoints
### `GET /artifact/{UUID}`
Retrieve an artifact (a JSON encoded concrete spec, concretization failure log, or build log) from the backend.

Example,
```bash
$ curl -X GET https://drift-server.spack.io/artifact/6c01dfe6-c6d5-4858-982b-e1f14ac7bc54
```

### `GET /concretizer-diff`
Retrieve a list of all of the inflection points found in one concretizer or another but not in both.

Example,
```bash
$ curl -X GET https://drift-server.spack.io/concretizer-diff
```

### `GET /concretizer-diff/{Abstract-Spec}`
Retrieve a list of all of the inflection points found in one concretizer or another but not in both **limited to a single abstract spec**.

Example,
```bash
$ curl -X GET https://drift-server.spack.io/concretizer-diff/ginkgo
```

### `GET /inflection-point/{ID}`
Retrieve a specific inflection point based on its' row ID in the database. Note: these ID's are **NOT** stable and may change depending on `PUT` and `GET` opperations.

Example,
```bash
$ curl -X GET https://drift-server.spack.io/inflection-point/205
```

### `GET /inflection-points`
Retrieve a list of all inflection points recorded in the database.

Example,
```
$ curl -X GET https://drift-server.spack.io/inflection-points
```

### `GET /inflection-points/{Abstract-Spec}`
Retrieve a list of all inflection points recorded in the database **for a specific abstract spec**.

Example,
```bash
$ curl -X GET https://drift-server.spack.io/inflection-points/ginkgo
```

### `GET /specs`
Retrieve a set of metadata about an abstract spec including the number of inflection points found and the number of known dependencies.

Example,
```bash
$ curl -X GET https://drift-server.spack.io/specs
```

### `GET /specs/{Abstract-Spec}`
Retrieve a set of metadata about an abstract spec including the number of inflection points found and the number of known dependencies. **Limit to a specific abstract spec.**

Example,
```bash
$ curl -X GET https://drift-server.spack.io/specs/ginkgo
```

### `PUT /inflection-point`
Upload a new inflection point or re-upload a modified inflection point to the database.

Example,
```bash
$ curl -X PUT -H "Content-Type: application/json" -d @point.json https://drift-server.spack.io/inflection-point
```

### `PUT /artifact`
Upload an artifact to the backend datastore (s3) and retrieve back a UUID to refrence the artifact by.

Example,
```bash
$ curl -X PUT -H "Content-Type: application/json" -d @artifact.json https://drift-server.spack.io/artifact
```

Response,
```json
{"uuid": 6c01dfe6-c6d5-4858-982b-e1f14ac7bc54}
```


### `DELETE /inflection-point`
Delete an inflection point from the database. Note: Make sure to delete any associated artifacts before you delete a inflection point.

Example,
```bash
$ curl -X DELETE -H "Content-Type: application/json" -d @point.json https://drift-server.spack.io/inflection-point
```

### `DELETE /artifact/UUID`
Delete an artifact from the backend.

Example,
```bash
$ curl -X DELETE https://drift-server.spack.io/artifact/6c01dfe6-c6d5-4858-982b-e1f14ac7bc54
```
