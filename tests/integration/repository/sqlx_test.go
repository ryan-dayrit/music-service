//go:build integration

package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	sqlxrepo "music-service/internal/repository/postgres/sqlx"
	"music-service/tests/integration/testutil"
)

func setupPostgresSqlx(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()
	ctx := context.Background()

	tmpFile, err := os.CreateTemp("", "integration_init_*.sql")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmpFile.WriteString(testutil.InitSQL); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to write init script: %v", err)
	}
	initScriptPath := tmpFile.Name()
	tmpFile.Close()
	t.Cleanup(func() { os.Remove(initScriptPath) })

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(initScriptPath),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping: %v", err)
	}

	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func TestSqlxRepository_Read_Empty(t *testing.T) {
	db, cleanup := setupPostgresSqlx(t)
	defer cleanup()

	repository := sqlxrepo.NewRepository(db)

	albums, err := repository.Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if len(albums) != 0 {
		t.Errorf("expected 0 albums, got %d", len(albums))
	}
}

func TestSqlxRepository_Read_WithData(t *testing.T) {
	db, cleanup := setupPostgresSqlx(t)
	defer cleanup()

	_, err := db.Exec("INSERT INTO music.albums (title, artist, price) VALUES ($1, $2, $3)",
		"Blue Train", "John Coltrane", 56.99)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}
	_, err = db.Exec("INSERT INTO music.albums (title, artist, price) VALUES ($1, $2, $3)",
		"Giant Steps", "John Coltrane", 63.99)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	repository := sqlxrepo.NewRepository(db)

	albums, err := repository.Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}
	if len(albums) != 2 {
		t.Fatalf("expected 2 albums, got %d", len(albums))
	}
	if albums[0].Title != "Blue Train" || albums[0].Artist != "John Coltrane" {
		t.Errorf("expected first album Blue Train by John Coltrane, got %q by %q",
			albums[0].Title, albums[0].Artist)
	}
	if albums[1].Title != "Giant Steps" || albums[1].Artist != "John Coltrane" {
		t.Errorf("expected second album Giant Steps by John Coltrane, got %q by %q",
			albums[1].Title, albums[1].Artist)
	}
	if !albums[0].Price.Equals(decimal.NewFromFloat(56.99)) {
		t.Errorf("expected first album price 56.99, got %s", albums[0].Price.String())
	}
	if !albums[1].Price.Equals(decimal.NewFromFloat(63.99)) {
		t.Errorf("expected second album price 63.99, got %s", albums[1].Price.String())
	}
}
