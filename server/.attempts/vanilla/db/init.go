package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver for database interaction.
)

// Open opens a database and creates one if not found.
func Open(databaseName string) (database *sql.DB) {

	var err error
	database, err = sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatal(err)
	}

	// Create the database and tables if they don't exist.
	statements := []string{
		"PRAGMA foreign_keys = ON",
		"CREATE TABLE IF NOT EXISTS Tags (id INT NOT NULL PRIMARY KEY, name TEXT NOT NULL UNIQUE)",
		"CREATE TABLE IF NOT EXISTS Packages (name TEXT NOT NULL PRIMARY KEY, version TEXT NOT NULL PRIMARY KEY, CONSTRAINT C_PackageUnique UNIQUE(name, version))",
		"CREATE TABLE IF NOT EXISTS InflectionTags (tag_id INT NOT NULL, FOREIGN KEY (tag_id) REFERENCES Tags(id), inflection_id INT NOT NULL, FOREIGN KEY (inflection_id) REFERENCES InflectionPoints(id), CONSTRAINT C_InflectionTagsUnique UNIQUE(tag_id, inflection_id))",
		"CREATE TABLE IF NOT EXISTS InflectionPoints(id INT NOT NULL PRIMARY KEY, commit TEXT NOT NULL, package_id INT NOT NULL, FOREIGN KEY (package_id) REFERENCES Packages(name), CONSTRAINT C_InflectionPointsUnique UNIQUE(commit, package_id))"}

	for _, statement := range statements {
		_, err = database.Exec(statement)
		if err != nil {
			log.Fatal(err)
		}
	}
	return database
}
