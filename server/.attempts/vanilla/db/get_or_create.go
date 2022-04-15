package db

import (
	"database/sql"
)

// Get or create a package
func GetOrCreatePackage(database *sql.DB, pkg Package) (Package, error) {

	// Make sure we still have a connection
	err := database.Ping()
	if err != nil {
		return pkg, err
	}

	// Query for Package based on name and Version
	row, err := database.Query("SELECT * FROM Packages WHERE name = ?, version = ?", pkg.Name, pkg.Version)
	if err != nil {
		return pkg, err
	}
	defer row.Close()

	// If we already have the record return it (there are only two fields)
	if row.Next() {
		return pkg, nil
	}

	// Otherwise, create it
	statement, err := database.Prepare(
		"INSERT INTO Packages(" +
			"name," +
			"version" +
			") VALUES(?,?);")

	if err != nil {
		return pkg, err
	}
	_, err = statement.Exec(
		pkg.Name,
		pkg.Version)

	return pkg, nil
}

// Get a tag or return nil
func GetTag(database *sql.DB, tag string) (string, error) {

	// Make sure we still have a connection
	err := database.Ping()
	if err != nil {
		return tag, err
	}

	// Query for Package based on name and Version
	row, err := database.Query("SELECT * FROM Tags WHERE name = ?", tag)
	if err != nil {
		return tag, err
	}
	defer row.Close()

	// If we already have the record return it
	if row.Next() {
		return tag, nil
	}

	// No error, and no tag
	return tag, nil
}

// Get or create a tag
func GetOrCreateTag(database *sql.DB, tag string) (string, error) {

	// Make sure we still have a connection
	err := database.Ping()
	if err != nil {
		return nil, err
	}

	new_tag, err := GetTag(database, tag)
	if err != nil {
		return nil, err
	}
	if new_tag != nil {
		return new_tag, nil
	}

	// Otherwise, create it
	statement, err := database.Prepare("INSERT INTO Tags(name) VALUES (?);")
	if err != nil {
		return nil, err
	}

	// Execute the statement
	tag, err = statement.Exec(tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// Add a new inflection point
func CreateInflectionPoint(database *sql.DB, point InflectionPoint) error {

	// Make sure we still have a connection
	err := database.Ping()
	if err != nil {
		return err
	}

	// Get or create the package
	pkg, err := GetOrCreatePackage(database, point.Pkg)
	if err != nil {
		return err
	}

	// Get or create the tags
	var tags_ids = []string{}
	for i, tag := range point.Tags {
		new_tag, err := GetOrCreateTag(database, tag)
		if err != nil {
			return err
		}
		tag_ids.push(new_tag.id)
	}

	// Query for InflectionPoint: unique together are commit and package
	row, err := database.Query("SELECT * FROM InflectionPoints WHERE commit = ?, package_id = ?", point.Commit, point.Pkg.Name)
	if err != nil {
		return err
	}
	defer row.Close()

	// If we already have the record, don't add again
	if row.Next() {
		return nil
	}

	// Create the new inflection point
	statement, err := db.Prepare(
		"INSERT INTO InflectionPoints(" +
			"commit," +
			"package_id" +
			") VALUES(?,?);")

	if err != nil {
		return err
	}

	_, err = statement.Exec(
		point.Commit,
		point.Pkg.Name)

	if err != nil {
		return err
	}
	return nil
}
