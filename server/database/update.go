package database

// Update attempts to modify an existing entry in the database.
func (db *DB) Update(entry InflectionPoint) (old InflectionPoint, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return old, err
	}

	// Check for entry in DB
	old, err = db.getPoint(entry.AbstractSpec, entry.GitCommit, entry.Concretizer)
	if err != nil {
		return old, err
	}

	// Update the entry if it exists.
	err = db.updatePoint(entry)
	return old, err
}

// updatePoint changes a file's status in the database.
func (db *DB) updatePoint(entry InflectionPoint) (err error) {
	stmt, err := db.conn.Prepare(
		`UPDATE points SET
			gitAuthorDate = ?,
			gitCommitDate = ?,
			specUUID = ?,
			built = ?,
			concretized = ?,
			isprimary = ?,
			buildLogUUID = ?,
			concretizationErrUUID = ?
			WHERE abstractSpec = ? AND gitCommit = ? AND concretizer = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.GitAuthorDate,
		entry.GitCommitDate,
		entry.SpecUUID,
		entry.Built,
		entry.Concretized,
		entry.Primary,
		entry.BuildLogUUID,
		entry.ConcretizationErrUUID,
		entry.AbstractSpec,
		entry.GitCommit,
		entry.Concretizer,
	)
	return err
}
