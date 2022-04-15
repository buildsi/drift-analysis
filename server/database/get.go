package database

import (
	"database/sql"
	"encoding/json"
)

// GetPoint searches for and returns a the corresponding entry from the
// database if the entry exists.
func (db *DB) GetPoint(in InflectionPoint) (result InflectionPoint, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	// Get and return entry from DB if it exists
	return db.getPoint(in.AbstractSpec, in.GitCommit, in.Concretizer)
}

// get returns the matching entry from the db if it exists.
func (db *DB) getPoint(abstractSpec, gitCommit, concretizer string) (result InflectionPoint, err error) {
	tags := ""
	files := ""
	row, err := db.conn.Query(
		`SELECT rowid, * FROM points WHERE abstractSpec = ?
		 AND gitCommit = ?
		 AND concretizer = ?`,
		abstractSpec,
		gitCommit,
		concretizer,
	)
	if err != nil {
		return result, err
	}
	defer row.Close()
	if !row.Next() {
		return result, sql.ErrNoRows
	}
	err = row.Scan(
		&result.ID,
		&result.AbstractSpec,
		&result.GitCommit,
		&result.GitAuthorDate,
		&result.GitCommitDate,
		&result.Concretizer,
		&files,
		&tags,
		&result.SpecUUID,
		&result.Built,
		&result.Concretized,
		&result.Primary,
		&result.BuildLogUUID,
		&result.ConcretizationErrUUID,
	)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(files), &result.Files)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(tags), &result.Tags)
	return result, err
}

func (db *DB) GetAll() (result []InflectionPoint, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Create files slice with limit as size.
	result = []InflectionPoint{}

	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	rows, err := db.conn.Query(
		"SELECT rowid, * FROM points ORDER BY gitCommitDate;",
	)
	if err != nil {
		return result, err
	}

	// Iterate through rows found and insert them into the list.
	for rows.Next() {
		tags := ""
		files := ""
		var in InflectionPoint

		err = rows.Scan(
			&in.ID,
			&in.AbstractSpec,
			&in.GitCommit,
			&in.GitAuthorDate,
			&in.GitCommitDate,
			&in.Concretizer,
			&files,
			&tags,
			&in.SpecUUID,
			&in.Built,
			&in.Concretized,
			&in.Primary,
			&in.BuildLogUUID,
			&in.ConcretizationErrUUID,
		)
		if err != nil {
			rows.Close()
			return nil, err
		}

		err = json.Unmarshal([]byte(tags), &in.Tags)
		if err != nil {
			rows.Close()
			return nil, err
		}

		err = json.Unmarshal([]byte(files), &in.Files)
		if err != nil {
			rows.Close()
			return nil, err
		}

		result = append(result, in)
	}

	// Check for errors and return
	err = rows.Close()
	if err != nil {
		return result, err
	}

	if len(result) <= 0 {
		return result, sql.ErrNoRows
	}

	return result, err
}

func (db *DB) GetAllWithSpec(abstractSpec string) (result []InflectionPoint, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Create files slice with limit as size.
	result = []InflectionPoint{}

	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	rows, err := db.conn.Query(
		"SELECT rowid, * FROM points WHERE abstractSpec = ? ORDER BY gitCommitDate;",
		abstractSpec,
	)
	if err != nil {
		return result, err
	}

	// Iterate through rows found and insert them into the list.
	for rows.Next() {
		tags := ""
		files := ""
		var in InflectionPoint

		err = rows.Scan(
			&in.ID,
			&in.AbstractSpec,
			&in.GitCommit,
			&in.GitAuthorDate,
			&in.GitCommitDate,
			&in.Concretizer,
			&files,
			&tags,
			&in.SpecUUID,
			&in.Built,
			&in.Concretized,
			&in.Primary,
			&in.BuildLogUUID,
			&in.ConcretizationErrUUID,
		)
		if err != nil {
			rows.Close()
			return nil, err
		}

		err = json.Unmarshal([]byte(files), &in.Files)
		if err != nil {
			rows.Close()
			return nil, err
		}

		err = json.Unmarshal([]byte(tags), &in.Tags)
		if err != nil {
			rows.Close()
			return nil, err
		}

		result = append(result, in)
	}

	// Check for errors and return
	err = rows.Close()
	if err != nil {
		return result, err
	}

	if len(result) <= 0 {
		return result, sql.ErrNoRows
	}

	return result, err
}

// GetArtifactType searches for and returns the corresponding artifact type
// from the database if the artifact exists.
func (db *DB) GetArtifactType(uuid string) (payloadType string, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return payloadType, err
	}

	// Get and return entry from DB if it exists
	return db.getArtifactType(uuid)
}

// get returns the matching entry from the db if it exists.
func (db *DB) getArtifactType(uuid string) (payloadType string, err error) {
	row, err := db.conn.Query(
		`SELECT payloadType FROM artifacts WHERE uuid = ?`,
		uuid,
	)
	if err != nil {
		return payloadType, err
	}
	defer row.Close()
	if !row.Next() {
		return payloadType, sql.ErrNoRows
	}
	err = row.Scan(
		&payloadType,
	)
	return payloadType, err
}
