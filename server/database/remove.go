package database

// RemovePoint deletes and returns an entry from the database.
func (db *DB) RemovePoint(in InflectionPoint) (result InflectionPoint, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	// Get the current value of the entry in the DB before removing
	result, err = db.getPoint(in.AbstractSpec, in.GitCommit, in.Concretizer)
	if err != nil {
		return result, err
	}

	// Remove the entry from the DB
	err = db.removePoint(in)
	return result, err
}

// removePoint deletes an entry from the DB.
func (db *DB) removePoint(in InflectionPoint) (err error) {
	stmt, err := db.conn.Prepare(
		"DELETE FROM points WHERE abstractSpec = ? AND gitCommit = ? AND concretizer = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		in.AbstractSpec,
		in.GitCommit,
		in.Concretizer,
	)
	return err
}

// RemoveArtifact deletes and returns an artifactType from the database.
func (db *DB) RemoveArtifactType(uuid string) (payloadType string, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return payloadType, err
	}

	// Get the current value of the entry in the DB before removing
	payloadType, err = db.getArtifactType(uuid)
	if err != nil {
		return payloadType, err
	}

	// Remove the entry from the DB
	err = db.removeArtifactType(uuid)
	return payloadType, err
}

// removeArtifact deletes an artifact from the DB.
func (db *DB) removeArtifactType(uuid string) (err error) {
	stmt, err := db.conn.Prepare(
		"DELETE FROM artifacts WHERE uuid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		uuid,
	)
	return err
}
