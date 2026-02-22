//go:build integration

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"music-service/gen/pb"
	"music-service/internal/repository/postgres/orm"
	"music-service/internal/routes"
	v1routes "music-service/internal/routes/v1"
	"music-service/tests/integration/testutil"
)

// noopProducer is a no-op Kafka producer for integration tests (avoids Kafka dependency)
type noopProducer struct{}

func (n *noopProducer) Produce(ctx context.Context, album *pb.Album) {}

func (n *noopProducer) Close() {}

func setupPostgres(t *testing.T) (*pg.DB, func()) {
	t.Helper()
	ctx := context.Background()

	// Write init script to temp file (testcontainers needs a file path)
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

func setupApp(db *pg.DB) *fiber.App {
	repository := orm.NewRepository(db)
	producer := &noopProducer{}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(cors.New())
	routes.RegisterSwaggerRoute(app)

	v1Router := app.Group("/api/v1")
	v1routes.RegisterHealthRoute(v1Router)
	v1routes.RegisterPublicRoutes(v1Router, producer, repository)

	return app
}

func TestRESTAlbums_GetAlbums_Empty(t *testing.T) {
	db, cleanup := setupPostgres(t)
	defer cleanup()

	app := setupApp(db)

	req, err := http.NewRequest("GET", "/api/v1/albums", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var albums []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&albums); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(albums) != 0 {
		t.Errorf("expected 0 albums, got %d", len(albums))
	}
}

func TestRESTAlbums_GetAlbums_WithData(t *testing.T) {
	db, cleanup := setupPostgres(t)
	defer cleanup()

	_, err := db.Exec("INSERT INTO music.albums (title, artist, price) VALUES (?, ?, ?)",
		"Blue Train", "John Coltrane", 56.99)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	_, err = db.Exec("INSERT INTO music.albums (title, artist, price) VALUES (?, ?, ?)",
		"Giant Steps", "John Coltrane", 63.99)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	app := setupApp(db)

	req, err := http.NewRequest("GET", "/api/v1/albums", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var albums []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&albums); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(albums) != 2 {
		t.Fatalf("expected 2 albums, got %d", len(albums))
	}

	// Verify first album (models.Album serializes with capitalized keys: Title, Artist)
	if title, ok := albums[0]["Title"].(string); !ok || title != "Blue Train" {
		t.Errorf("expected first album title 'Blue Train', got %v", albums[0]["Title"])
	}
	if artist, ok := albums[0]["Artist"].(string); !ok || artist != "John Coltrane" {
		t.Errorf("expected first album artist 'John Coltrane', got %v", albums[0]["Artist"])
	}

	// Verify second album
	if title, ok := albums[1]["Title"].(string); !ok || title != "Giant Steps" {
		t.Errorf("expected second album title 'Giant Steps', got %v", albums[1]["Title"])
	}
}

func TestRESTAlbums_PostAlbums(t *testing.T) {
	db, cleanup := setupPostgres(t)
	defer cleanup()

	app := setupApp(db)

	newAlbums := []*pb.Album{
		{Id: 1, Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{Id: 2, Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	}
	body, err := json.Marshal(newAlbums)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/albums", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestRESTAlbums_PostAlbums_InvalidJSON(t *testing.T) {
	db, cleanup := setupPostgres(t)
	defer cleanup()

	app := setupApp(db)

	req, err := http.NewRequest("POST", "/api/v1/albums", bytes.NewReader([]byte("invalid json")))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if errMsg, ok := result["error"].(string); !ok || errMsg != "cannot parse JSON" {
		t.Errorf("expected error 'cannot parse JSON', got %v", result["error"])
	}
}
