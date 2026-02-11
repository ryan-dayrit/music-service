package album

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
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

func TestRepository_Read(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	// Set up expected query and results
	rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
		AddRow(1, "Blue Train", "John Coltrane", decimal.NewFromFloat(56.99)).
		AddRow(2, "Giant Steps", "John Coltrane", decimal.NewFromFloat(63.99))

	mock.ExpectQuery("SELECT id, title, artist, price FROM music.albums").
		WillReturnRows(rows)

	albums, err := repo.Read()
	if err != nil {
		t.Fatalf("Read() returned unexpected error: %v", err)
	}

	if len(albums) != 2 {
		t.Errorf("Expected 2 albums, got %d", len(albums))
	}

	if albums[0].Title != "Blue Train" {
		t.Errorf("Expected first album title 'Blue Train', got '%s'", albums[0].Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRepository_Read_EmptyResult(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	// Set up expected query with no results
	rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"})

	mock.ExpectQuery("SELECT id, title, artist, price FROM music.albums").
		WillReturnRows(rows)

	albums, err := repo.Read()
	if err != nil {
		t.Fatalf("Read() returned unexpected error: %v", err)
	}
	if albums == nil {
		t.Error("Expected non-nil albums slice, got nil")
	}

	if len(albums) != 0 {
		t.Errorf("Expected 0 albums, got %d", len(albums))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRepository_Read_ScanError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	// Set up expected query with invalid data that will cause scan error
	rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
		AddRow("invalid", "Blue Train", "John Coltrane", decimal.NewFromFloat(56.99))

	mock.ExpectQuery("SELECT id, title, artist, price FROM music.albums").
		WillReturnRows(rows)

	_, err = repo.Read()
	if err == nil {
		t.Error("Expected error from invalid data, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	repo := NewRepository(db)

	// Set up expected query and results for each iteration
	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow(1, "Blue Train", "John Coltrane", decimal.NewFromFloat(56.99)).
			AddRow(2, "Giant Steps", "John Coltrane", decimal.NewFromFloat(63.99))

		mock.ExpectQuery("SELECT id, title, artist, price FROM music.albums").
			WillReturnRows(rows)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Read()
	}
}
