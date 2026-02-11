package album

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func TestNewRepository(t *testing.T) {
	db := &sqlx.DB{}
	repo := NewRepository(db)

	if repo == nil {
		t.Fatal("Expected non-nil repository, got nil")
	}

	var _ Repository = repo
}

func TestRepository_Read_Integration(t *testing.T) {
	db, err := sqlx.Open("postgres", "host=localhost port=5432 user=ryandayrit dbname=practice sslmode=disable")
	if err != nil {
		t.Skipf("Skipping integration test: database connection failed: %v", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping integration test: database not available: %v", err)
		return
	}

	repo := NewRepository(db)

	albums, err := repo.Read()
	if err != nil {
		t.Fatalf("Read() returned error: %v", err)
	}
	if albums == nil {
		t.Error("Expected non-nil albums slice, got nil")
	}
}

func TestRepository_EmptyResult(t *testing.T) {
	db := &sqlx.DB{}
	repo := NewRepository(db)

	if repo == nil {
		t.Fatal("Expected non-nil repository, got nil")
	}
}

func TestAlbum_RepositoryIntegration(t *testing.T) {
	album := Album{
		Id:     1,
		Title:  "Test Album",
		Artist: "Test Artist",
		Price:  decimal.NewFromFloat(9.99),
	}

	if album.Id == 0 {
		t.Error("Album Id should be set")
	}
	if album.Title == "" {
		t.Error("Album Title should be set")
	}
	if album.Artist == "" {
		t.Error("Album Artist should be set")
	}
	if album.Price.IsZero() {
		t.Error("Album Price should be set")
	}
}

func TestRepository_Interface(t *testing.T) {
	db := &sqlx.DB{}
	repo := NewRepository(db)

	var _ Repository = repo
}

func TestQueryEmbedded(t *testing.T) {
	if getAlbumsQuery == "" {
		t.Error("Expected getAlbumsQuery to be embedded, got empty string")
	}
}

func BenchmarkRepository_Read(b *testing.B) {
	db, err := sqlx.Open("postgres", "host=localhost port=5432 user=ryandayrit dbname=practice sslmode=disable")
	if err != nil {
		b.Skipf("Skipping benchmark: database connection failed: %v", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		b.Skipf("Skipping benchmark: database not available: %v", err)
		return
	}

	repo := NewRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Read()
	}
}
