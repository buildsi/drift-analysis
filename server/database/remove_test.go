package database

import (
	"database/sql"
	"testing"
)

func TestRemoveNoEntry(t *testing.T) {
	// Initialize mock db
	db, err := initMock()
	if err != nil {
		t.Fatal(err)
	}

	in := InflectionPoint{
		AbstractSpec: "abyss@1.1",
		GitCommit:    "i-am-not-a-real-digest",
		Concretizer:  "original",
	}

	// Test getting an entry from an empty DB
	_, err = db.RemovePoint(in)
	if err == nil || err != sql.ErrNoRows {
		t.Fatal("Test did not return a no rows error as expected")
	}
}

func TestRemoveEntry(t *testing.T) {
	// Initialize mock db
	db, err := initMock()
	if err != nil {
		t.Fatal(err)
	}

	in := InflectionPoint{
		AbstractSpec: "abyss@1.1",
		GitCommit:    "i-am-not-a-real-digest",
		Concretizer:  "original",
	}

	// Add data to mock db
	err = db.AddPoint(in)
	if err != nil {
		t.Fatal(err)
	}

	// Ask for data back from db
	out, err := db.RemovePoint(in)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the in and out commits match
	if out.GitCommit != in.GitCommit {
		t.Fatalf("expected inflection point with commit %s but got %s instead", in.GitCommit, out.GitCommit)
	}
	// Check in and out package names match
	if out.AbstractSpec != in.AbstractSpec {
		t.Fatalf("expected package name with %s but got %s instead", in.AbstractSpec, out.AbstractSpec)
	}

	// Test that the entry was removed as expected.
	_, err = db.GetPoint(in)
	if err == nil || err != sql.ErrNoRows {
		t.Fatal("Test did not return a no rows error as expected")
	}
}
