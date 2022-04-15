# server

The drift server is a simple application and database to store results for
a drifter analysis. This means that:

1. Run a [Drift Analysis](https://github.com/buildsi/drift-analysis) to generate inflection points, commits, and tags
2. Save the inflection points, commits, tags, and package string to the server.
3. Build the specs, and record the outcome.

And if there is repeat of a commit / spec full hash (e.g., something is already in the
database) we would want to ignore it.

## Authentication

By default, authentication is disabled, meaning that the username and password
are empty strings. To change this, export the following variables when you
run the server:

```bash
export DRIFTSERVER_USER=myusername
export DRIFTSERVER_PASS=mypassword
```

For any endpoint described below, you can add basic authentication to curl as follows:

```bash
$ curl -u $DRIFTSERVER_USER:$DRIFTSERVER_PASS ...
```

## Environment

By default, the server will run on port 8080 and use an sqlite databsae. However, you
can set the following environment variables to change that! Notice that the username
and password shown above are part of the environment variable set. This is because
they are part of the same server config.

| variable | description | default |
|-----------|------------|----------|
| DRIFTSERVER_PORT | The port to run the server on | 8080 |
| DRIFTSERVER_DEBUG| Currently not used, will be integrated when logging is better designed | false |
| DRIFTSERVER_DB_PATH| The path to the drift server sqlite database.| "points.db" |
| DRIFTSERVER_USER | The username to authenticate with basic | "" |
| DRIFTSERVER_PASS | The password to authenticate with basic | "" |


## Endpoints

### 1. Upload Inflection Point

We then start with an [inflection_point.json](data/inflection_point.json) data file, which looks like this:

```json
{
    "commit": {"digest": "12345", "timestamp": "2021-12-10T21:20:15+00:00"},
    "tags": [{"name": "one"}, {"name": "two"}, {"name": "three"}],
    "abstract_spec": "abyss@1.1",
    "concretizer": "original"
}
```

You can then upload the file as follows:

```bash
$ curl -X POST -H "Content-Type: application/json" -d @data/inflection_point.json http://localhost:8080/inflection-point/
```

The result of this operation will give you an ID for the inflection point, which you
need to save. If you ever forget it, just upload the same data structure to look it up
again. Note that we are reading the file and sending json data, and not uploading the file
as a multipart (form) upload. This is to allow sending directly from Python without
needing to save to file first.

### 2. List inflection points

You can list inflection points as follows:

```bash
$ curl -X GET -H "Content-Type: application/json" http://localhost:8080/inflection-points/
```
