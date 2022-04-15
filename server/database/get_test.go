package database

import (
	"database/sql"
	"testing"
)

func TestGetNoEntry(t *testing.T) {
	// Initialize mock db
	db, err := initMock()
	if err != nil {
		t.Fatal(err)
	}

	// Test getting an entry from an empty DB
	_, err = db.GetPoint(InflectionPoint{})
	if err == nil || err != sql.ErrNoRows {
		t.Fatal("Test did not return a no rows error as expected")
	}
}

func TestGetEntry(t *testing.T) {
	// Initialize mock db
	db, err := initMock()
	if err != nil {
		t.Fatal(err)
	}

	in := InflectionPoint{
		AbstractSpec: "abyss@1.1",
		GitCommit:    "i-am-not-a-real-digest",
		Concretizer:  "original",
		Tags:         []string{"hello", "world"},
		Files:        []string{"i", "am", "files"},
	}

	// Add data to mock db
	err = db.AddPoint(in)
	if err != nil {
		t.Fatal(err)
	}

	// Ask for data back from db
	out, err := db.GetPoint(in)
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
}
