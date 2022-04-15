package database

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import sqlite driver for database interaction.
)

type DB struct {
	conn *sql.DB
	lock *sync.Mutex
}

// InflectionPoint is the information saved in the database
// for each inflection point found for an abstract spec within
// the spack repository.
type InflectionPoint struct {
	// ID is a autoincrement value added for possible visualizations
	ID int64

	// AbstractSpec is the string form of the input spec.
	AbstractSpec string

	// GitCommit is the commit found at which there was an inflection point.
	GitCommit     string
	GitAuthorDate time.Time
	GitCommitDate time.Time

	// Which concretizer was used for the inflection point.
	Concretizer string

	// List of files changed within an inflection point.
	Files []string

	// List of tags associated with the inflection point.
	Tags []string

	// SpecUUID is a UUID to the JSON encoded concrete spec
	SpecUUID string

	// Built and Concretized booleans
	// for quickly filtering inflection points.
	Built       bool
	Concretized bool

	// Primary is a boolean signifying if the inflection was caused directly
	// by this abstract spec
	Primary bool

	// BuildLogUUID is a UUID of the output/error of a build.
	BuildLogUUID string

	// ConcretizationErrUUID is a UUID of the output of a failed concretization.
	ConcretizationErrUUID string
}

// Init opens and connects to the database.
func Init(path string) (result *DB, err error) {
	db, err := sql.Open("sqlite3", path+"?cache=shared")
	if err != nil {
		return nil, err
	}

	// Create the nodes table if is doesn't already exist.
	// This will also create the database if it doesn't exist.
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS points(
			abstractSpec TEXT NOT NULL,
			gitCommit TEXT NOT NULL,
			gitAuthorDate DATETIME,
			gitCommitDate DATETIME,
			concretizer TEXT NOT NULL,
			files TEXT,
			tags TEXT,
			specUUID TEXT,
			built BIT,
			concretized BIT,
			isprimary BIT,
			buildLogUUID TEXT,
			concretizationErrUUID TEXT,
			PRIMARY KEY(abstractSpec, gitCommit, concretizer)
		);
		CREATE TABLE IF NOT EXISTS artifacts(
			uuid TEXT NOT NULL,
			payloadType TEXT NOT NULL,
			PRIMARY KEY(uuid)
		);
		PRAGMA busy_timeout = 5000;`,
	)
	if err != nil {
		return nil, err
	}

	result = &DB{
		conn: db,
		lock: &sync.Mutex{},
	}
	return result, nil
}
