package db

import (
	"testing"
)

func TestNewPostgresDB_Integration(t *testing.T) {
	cfg := Config{DriverName: "postgres", User: "ryandayrit", DBName: "practice", SSLMode: "disable", Password: "", Host: "localhost"}
	db, err := NewPostgresDB(cfg)
	if err != nil {
		t.Skipf("Skipping: %v", err)
		return
	}
	defer db.Close()
	if db == nil {
		t.Fatal("Expected non-nil database")
	}
}

func TestNewPostgresDB_EmptyConfig(t *testing.T) {
	_, err := NewPostgresDB(Config{})
	if err == nil {
		t.Error("Expected error with empty config")
	}
}
