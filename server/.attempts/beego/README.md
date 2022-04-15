# Drift Server

The drifer server is a simple application and database to store results for
a drifter analysis. This means that:

1. Run a [Drift Analysis](https://github.com/buildsi/drift-analysis) to generate inflection points, commits, and tags
2. Save the inflection points, commits, tags, and package string to the server.
3. Build the specs, and record the outcome. 

And if there is repeat of a commit / spec full hash (e.g., something is already in the
database) we would want to ignoreit.

## 0. Background

I started out testing using basic [database/sql](https://golang.org/pkg/database/sql/)
and developed most of the first draft, but it was bothersome to need to implement the
basics (of what I would expect any basic framework to provide!) I did some research
and stumbled on [beego](https://beego.me/) which was recommended for Django 
enthusiasts (of which I am one). Thus, I created a new API project:

```bash
$ go get github.com/beego/beego/v2@v2.0.0
$ go get github.com/beego/bee/v2
```

Add go bin to the path

```bash
$ export PATH=~/go/bin:$PATH
```

Generate the api

```bash
$ bee api drift-server
```

And here we are! I then customized everything with my database setup, and models.
For now we are going to use sqlite so I can easily throw it away, and eventually
I expect we will want sql or postgres.

## 2. Organization

The framework is organized as follows:

 - [models](models): includes database models, and functions to get/create/list etc. each.
 - [controllers](controllers): are functions associated with models to GET, CREATE, etc.
 - [routers](routers): set up how to direct controllers to specific web addresses


## 3. Run the server

It's easy to run the server locally in a development sense:

```bash
$ go run main.go

# or!
$ bee run
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

## TODOs

 - need to add authentication
