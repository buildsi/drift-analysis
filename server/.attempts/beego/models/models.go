package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// TODO: can we have an abstract base model that includes times?
// TODO: we probably want to think about m2m reverse names / on delete

// Commits are unique by their digest
type Commit struct {
	Id        int    `orm:"auto"`
	Digest    string `orm:"unique"`
	Timestamp string
	Created   time.Time `orm:"auto_now_add;type(datetime)"`
	Updated   time.Time `orm:"auto_now;type(datetime)"`
}

// Model to serialize a commit request
type CommitRequest struct {
	Digest    string `json:"digest"`
	Timestamp string `json:"timestamp"`
}

// A tag is just a unique string
type Tag struct {
	Name    string    `orm:"pk;unique"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

// Model to serialize a tag
type TagRequest struct {
	Name string `json:"name"`
}

// A Spack package has a name and version
type Package struct {
	Id      int
	Name    string
	Version string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

// packages have name and version unique together
func (p *Package) TableUnique() [][]string {
	return [][]string{
		{"Name", "Version"},
	}
}

// Model to serialize a package request
type PackageRequest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// A spack spec is an instance of a package
type Spec struct {
	Package *Package  `orm:"unique,rel(fk)";json:"package"` // foreign key
	Data    string    `orm:"type(text)";json:"spec"`        // we will need to store json as text
	Created time.Time `orm:"auto_now_add;type(datetime)";json:"created"`
	Updated time.Time `orm:"auto_now;type(datetime)";json:"updated"`
}

type InflectionPoint struct {
	Id      int
	Commit  *Commit
	Package *Package  `orm:"rel(fk);"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

// Model for many to many relationship
type InflectionTag struct {
	Id              int
	InflectionPoint *InflectionPoint `orm:"rel(fk)"`
	Tag             *Tag             `orm:"rel(fk)"`
	CreatedAt       time.Time        `orm:"type(datetime)"`
	UpdatedAt       time.Time        `orm:"type(datetime)"`
}

// Model to serialize inflection point request
type InflectionPointRequest struct {
	Commit  CommitRequest   `json:"commit"`
	Tags    []*TagRequest   `json:"tags"` // m2m == many to many relationship
	Package *PackageRequest `json:"package"`
}

// A commit and package are unique together for a point
func (p *InflectionPoint) TableUnique() [][]string {
	return [][]string{
		{"Commit", "Package"},
	}
}

type Build struct {
	Spec             *Spec            `orm:"rel(fk)";json:"spec"`
	Inflection_point *InflectionPoint `orm:"rel(fk)";json:"poing"`
	Status           string           `json:"status"`
	Created          time.Time        `orm:"auto_now_add;type(datetime)";json:"created"`
	Updated          time.Time        `orm:"auto_now;type(datetime)";json:"updated"`
}

func init() {

	// Register sqlite driver
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "drift-server.sqlite")

	// Register the new models
	orm.RegisterModel(new(Commit), new(Tag), new(Package), new(InflectionTag), new(InflectionPoint))
	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		logs.Warn(err)
	}
}
