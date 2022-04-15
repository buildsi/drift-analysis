package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3" // Import sqlite driver for database interaction.
)

// initMock opens up a mock DB for unit testing purposes.
func initMock() (result *DB, err error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
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
		PRAGMA busy_timeout = 5000;`,
	)
	if err != nil {
		return nil, err
	}

	result = &DB{
		conn: db,
		lock: &sync.Mutex{},
	}

	return result, err
}
