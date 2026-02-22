//go:build integration

package repository

import (
	"context"
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/shopspring/decimal"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"music-service/internal/repository/postgres/orm"
	"music-service/tests/integration/testutil"
)

func setupPostgres(t *testing.T) (*pg.DB, func()) {
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

	opts, err := pg.ParseURL(connStr)
	if err != nil {
		t.Fatalf("failed to parse connection string: %v", err)
	}

	db := pg.Connect(opts)
	if err := db.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func TestORMRepository_Get_Integration(t *testing.T) {
	db, cleanup := setupPostgres(t)
	defer cleanup()

	repository := orm.NewRepository(db)

	// Get empty
	albums, err := repository.Get()
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if len(albums) != 0 {
		t.Errorf("expected 0 albums, got %d", len(albums))
	}

	// Insert via raw SQL
	_, err = db.Exec("INSERT INTO music.albums (title, artist, price) VALUES (?, ?, ?)", "Test Album", "Test Artist", 9.99)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	// Get via repository
	albums, err = repository.Get()
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if len(albums) != 1 {
		t.Fatalf("expected 1 album, got %d", len(albums))
	}
	if albums[0].Title != "Test Album" {
		t.Errorf("expected title 'Test Album', got %s", albums[0].Title)
	}
	if albums[0].Artist != "Test Artist" {
		t.Errorf("expected artist 'Test Artist', got %s", albums[0].Artist)
	}
}

func TestORMRepository_GetById_Integration(t *testing.T) {
	db, cleanup := setupPostgres(t)
	defer cleanup()

	repository := orm.NewRepository(db)

	_, err := db.Exec("INSERT INTO music.albums (title, artist, price) VALUES (?, ?, ?)",
		"GetById Test", "Test Artist", 19.99)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	var result struct {
		Id int
	}
	_, err = db.QueryOne(&result, "SELECT id FROM music.albums WHERE title = ?", "GetById Test")
	if err != nil {
		t.Fatalf("failed to get id: %v", err)
	}
	id := result.Id

	album, err := repository.GetById(id)
	if err != nil {
		t.Fatalf("GetById(%d) failed: %v", id, err)
	}
	if album == nil {
		t.Fatal("expected album, got nil")
	}
	if album.Title != "GetById Test" {
		t.Errorf("expected title 'GetById Test', got %s", album.Title)
	}
	if !album.Price.Equals(decimal.NewFromFloat(19.99)) {
		t.Errorf("expected price 19.99, got %s", album.Price.String())
	}
}
