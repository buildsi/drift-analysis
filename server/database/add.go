package database

import (
	"database/sql"
	"encoding/json"
	"errors"
)

// AddPoint inserts an InflectionPoint entry into the database if it doesn't exist already.
func (db *DB) AddPoint(input InflectionPoint) (err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return err
	}

	_, err = db.getPoint(input.AbstractSpec, input.GitCommit, input.Concretizer)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.insertPoint(input)
		} else {
			return err
		}
	} else {
		err = db.updatePoint(input)
	}
	return err
}

// insertPoint adds a point entry to the database.
func (db *DB) insertPoint(entry InflectionPoint) (err error) {
	files, err := json.Marshal(entry.Files)
	if err != nil {
		return err
	}

	tags, err := json.Marshal(entry.Tags)
	if err != nil {
		return err
	}

	stmt, err := db.conn.Prepare(
		`INSERT INTO points(
			abstractSpec,
			gitCommit,
			gitAuthorDate,
			gitCommitDate,
			concretizer,
			files,
			tags,
			specUUID,
			built,
			concretized,
			isprimary,
			buildLogUUID,
			concretizationErrUUID
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.AbstractSpec,
		entry.GitCommit,
		entry.GitAuthorDate,
		entry.GitCommitDate,
		entry.Concretizer,
		string(files),
		string(tags),
		entry.SpecUUID,
		entry.Built,
		entry.Concretized,
		entry.Primary,
		entry.BuildLogUUID,
		entry.ConcretizationErrUUID,
	)
	return err
}

func (db *DB) AddArtifactType(uuid, payloadType string) (err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return err
	}

	_, err = db.getArtifactType(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.insertArtifactType(uuid, payloadType)
		} else {
			return err
		}
	} else {
		return errors.New("artifact uuid already in db")
	}
	return err
}

// insertSpec adds a point entry to the database.
func (db *DB) insertArtifactType(uuid, payloadType string) (err error) {
	stmt, err := db.conn.Prepare(
		`INSERT INTO artifacts(
			uuid,
			payloadType
		) VALUES(?,?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		uuid,
		payloadType,
	)
	return err
}
