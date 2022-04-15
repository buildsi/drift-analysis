# Drift Server

The drifer server is a simple application and database to store results for
a drifter analysis. This means that:

1. Run a [Drift Analysis](https://github.com/buildsi/drift-analysis) to generate inflection points, commits, and tags
2. Save the inflection points, commits, tags, and package string to the server.
3. Build the specs, and record the outcome. 

And if there is repeat of a commit / spec full hash (e.g., something is already in the
database) we would want to ignoreit.

## 1. Run the server

It's easy to run the server locally in a development sense:

```bash
$ go run main.go
```

## 2. Upload Inflection Point

We then start with an [inflection_point.json](inflection_point.json) data file, which looks like this:

```json
{
    "commit": {"name": "12345", "timestamp": "2021-12-10"},
    "tags": ["one", "two", "three"],
    "package": {"name": "packageName", "version": "packageVersion"}
}
```

You can then upload the file as follows:

```bash
$ curl -F "upload=@inflection_point.json" http://localhost:8080/inflection-point/new/
```