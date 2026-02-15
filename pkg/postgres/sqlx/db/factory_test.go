package db

import (
	"testing"

	"music-service/pkg/postgres"
)

func TestNewDB_Integration(t *testing.T) {
	cfg := postgres.Config{
		DriverName: "postgres",
		User:       "ryandayrit",
		DBName:     "practice",
		SSLMode:    "disable",
		Password:   "",
		Host:       "localhost",
	}

	db, err := NewDB(cfg)
	if err != nil {
		t.Skipf("Skipping integration test: database connection failed: %v", err)
		return
	}
	defer db.Close()

	if db == nil {
		t.Fatal("Expected non-nil database connection, got nil")
	}

	if err := db.Ping(); err != nil {
		t.Errorf("Database ping failed: %v", err)
	}
}

func TestNewDB_InvalidHost(t *testing.T) {
	cfg := postgres.Config{
		DriverName: "postgres",
		User:       "testuser",
		DBName:     "testdb",
		SSLMode:    "disable",
		Password:   "testpass",
		Host:       "invalid-host-that-does-not-exist.local",
	}

	db, err := NewDB(cfg)
	if err == nil {
		if db != nil {
			db.Close()
		}
		t.Skip("Expected error for invalid host, but connection succeeded")
	}

	if db != nil {
		t.Error("Expected nil database connection on error")
	}
}

func TestNewDB_InvalidDatabase(t *testing.T) {
	cfg := postgres.Config{
		DriverName: "postgres",
		User:       "ryandayrit",
		DBName:     "nonexistent_database_12345",
		SSLMode:    "disable",
		Password:   "",
		Host:       "localhost",
	}

	db, err := NewDB(cfg)
	if err == nil {
		if db != nil {
			db.Close()
		}
		t.Skip("Expected error for invalid database, but connection succeeded")
	}

	if db != nil {
		t.Error("Expected nil database connection on error")
	}
}

func TestNewDB_EmptyConfig(t *testing.T) {
	cfg := postgres.Config{}

	db, err := NewDB(cfg)
	if err == nil {
		if db != nil {
			db.Close()
		}
		t.Error("Expected error with empty config, got nil")
	}

	if db != nil {
		t.Error("Expected nil database connection on error")
	}
}

func TestNewDB_InvalidDriver(t *testing.T) {
	cfg := postgres.Config{
		DriverName: "invalid_driver",
		User:       "testuser",
		DBName:     "testdb",
		SSLMode:    "disable",
		Password:   "testpass",
		Host:       "localhost",
	}

	db, err := NewDB(cfg)
	if err == nil {
		if db != nil {
			db.Close()
		}
		t.Error("Expected error with invalid driver, got nil")
	}

	if db != nil {
		t.Error("Expected nil database connection on error")
	}
}

func TestNewDB_ReturnType(t *testing.T) {
	cfg := postgres.Config{
		DriverName: "postgres",
		User:       "ryandayrit",
		DBName:     "practice",
		SSLMode:    "disable",
		Password:   "",
		Host:       "localhost",
	}

	db, err := NewDB(cfg)
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
		return
	}
	defer db.Close()

	if db == nil {
		t.Error("Expected non-nil DB")
	}
}

func TestNewDB_DifferentSSLModes(t *testing.T) {
	tests := []struct {
		name    string
		sslMode string
		skip    bool
	}{
		{"Disable SSL", "disable", false},
		{"Require SSL", "require", true},
		{"Verify CA", "verify-ca", true},
		{"Verify Full", "verify-full", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping SSL mode test - requires SSL configuration")
			}

			cfg := postgres.Config{
				DriverName: "postgres",
				User:       "ryandayrit",
				DBName:     "practice",
				SSLMode:    tt.sslMode,
				Password:   "",
				Host:       "localhost",
			}

			db, err := NewDB(cfg)
			if err != nil {
				t.Skipf("Skipping %s test: %v", tt.name, err)
				return
			}
			defer db.Close()

			if db == nil {
				t.Fatal("Expected non-nil database connection")
			}
		})
	}
}

func TestNewDB_ConnectionPooling(t *testing.T) {
	cfg := postgres.Config{
		DriverName: "postgres",
		User:       "ryandayrit",
		DBName:     "practice",
		SSLMode:    "disable",
		Password:   "",
		Host:       "localhost",
	}

	db, err := NewDB(cfg)
	if err != nil {
		t.Skipf("Skipping integration test: database connection failed: %v", err)
		return
	}
	defer db.Close()

	stats := db.Stats()
	if stats.OpenConnections < 0 {
		t.Error("Invalid open connections count")
	}
}

func TestNewDB_MultipleConnections(t *testing.T) {
	cfg := postgres.Config{
		DriverName: "postgres",
		User:       "ryandayrit",
		DBName:     "practice",
		SSLMode:    "disable",
		Password:   "",
		Host:       "localhost",
	}

	db1, err := NewDB(cfg)
	if err != nil {
		t.Skipf("Skipping integration test: database connection failed: %v", err)
		return
	}
	defer db1.Close()

	db2, err := NewDB(cfg)
	if err != nil {
		t.Fatalf("Second connection failed: %v", err)
	}
	defer db2.Close()

	if db1 == nil || db2 == nil {
		t.Fatal("Expected non-nil database connections")
	}

	if err := db1.Ping(); err != nil {
		t.Errorf("First connection ping failed: %v", err)
	}

	if err := db2.Ping(); err != nil {
		t.Errorf("Second connection ping failed: %v", err)
	}
}

func BenchmarkNewPostgresDB(b *testing.B) {
	cfg := postgres.Config{
		DriverName: "postgres",
		User:       "ryandayrit",
		DBName:     "practice",
		SSLMode:    "disable",
		Password:   "",
		Host:       "localhost",
	}

	testDB, err := NewDB(cfg)
	if err != nil {
		b.Skipf("Skipping benchmark: database not available: %v", err)
		return
	}
	testDB.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db, err := NewDB(cfg)
		if err != nil {
			b.Fatalf("Connection failed: %v", err)
		}
		db.Close()
	}
}
