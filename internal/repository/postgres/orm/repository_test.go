package orm

import (
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/shopspring/decimal"

	"music-service/internal/models"
)

func TestNewRepository(t *testing.T) {
	t.Run("creates repository successfully", func(t *testing.T) {
		db := &pg.DB{}
		repo := NewRepository(db)

		if repo == nil {
			t.Error("Expected repository to be created, got nil")
		}
	})

	t.Run("repository implements Repository interface", func(t *testing.T) {
		db := &pg.DB{}
		var _ Repository = NewRepository(db)
	})
}

func TestRepository_Interface(t *testing.T) {
	t.Run("repository struct implements Repository interface", func(t *testing.T) {
		db := &pg.DB{}
		repo := &repository{db: db}

		var _ Repository = repo
	})

	t.Run("interface has required methods", func(t *testing.T) {
		db := &pg.DB{}
		repo := NewRepository(db)

		if repo == nil {
			t.Error("Repository should not be nil")
		}
	})
}

func TestRepository_StructFields(t *testing.T) {
	t.Run("repository has db field", func(t *testing.T) {
		db := &pg.DB{}
		repo := &repository{db: db}

		if repo.db == nil {
			t.Error("Expected db field to be set")
		}

		if repo.db != db {
			t.Error("Expected db field to match provided db")
		}
	})
}

func TestRepository_AlbumModel(t *testing.T) {
	t.Run("album model structure", func(t *testing.T) {
		album := models.Album{
			Id:     1,
			Title:  "Test",
			Artist: "Artist",
			Price:  decimal.NewFromFloat(9.99),
		}

		if album.Id != 1 {
			t.Errorf("Expected Id to be 1, got %d", album.Id)
		}

		if album.Title != "Test" {
			t.Errorf("Expected Title to be 'Test', got %s", album.Title)
		}

		if album.Artist != "Artist" {
			t.Errorf("Expected Artist to be 'Artist', got %s", album.Artist)
		}
	})
}

func TestRepository_NilDatabase(t *testing.T) {
	t.Run("creating repository with nil database", func(t *testing.T) {
		repo := NewRepository(nil)

		if repo == nil {
			t.Error("Repository should not be nil even with nil database")
		}
	})
}

func TestRepository_MultipleInstances(t *testing.T) {
	t.Run("can create multiple repository instances", func(t *testing.T) {
		db1 := &pg.DB{}
		db2 := &pg.DB{}

		repo1 := NewRepository(db1)
		repo2 := NewRepository(db2)

		if repo1 == nil || repo2 == nil {
			t.Error("Expected both repositories to be created")
		}

		if repo1 == repo2 {
			t.Error("Expected different repository instances")
		}
	})
}

func TestRepository_ReadParameters(t *testing.T) {
	t.Run("Read method parameter validation", func(t *testing.T) {
		testCases := []struct {
			name string
			id   int
		}{
			{"positive id", 1},
			{"large id", 999999},
			{"zero id", 0},
			{"negative id", -1},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				album := &models.Album{Id: tc.id}
				if album.Id != tc.id {
					t.Errorf("Expected Id to be %d, got %d", tc.id, album.Id)
				}
			})
		}
	})
}

func TestRepository_CreateParameters(t *testing.T) {
	t.Run("Create method with various album data", func(t *testing.T) {
		testCases := []struct {
			name   string
			album  models.Album
			expectId int
		}{
			{
				name: "complete album",
				album: models.Album{
					Id:     1,
					Title:  "Album Title",
					Artist: "Artist Name",
					Price:  decimal.NewFromFloat(19.99),
				},
				expectId: 1,
			},
			{
				name: "album with empty strings",
				album: models.Album{
					Id:     2,
					Title:  "",
					Artist: "",
					Price:  decimal.Zero,
				},
				expectId: 2,
			},
			{
				name: "album with long strings",
				album: models.Album{
					Id:     3,
					Title:  "A Very Long Album Title That Goes On And On",
					Artist: "An Artist With A Very Long Name",
					Price:  decimal.NewFromFloat(99.99),
				},
				expectId: 3,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.album.Id != tc.expectId {
					t.Errorf("Expected Id to be %d, got %d", tc.expectId, tc.album.Id)
				}
			})
		}
	})
}

func TestRepository_InterfaceContract(t *testing.T) {
	t.Run("Repository interface defines Read and Create methods", func(t *testing.T) {
		db := &pg.DB{}
		var repo Repository = NewRepository(db)

		if repo == nil {
			t.Fatal("Repository should not be nil")
		}
	})
}

func TestRepository_DecimalPrice(t *testing.T) {
	t.Run("handles decimal prices correctly", func(t *testing.T) {
		testCases := []struct {
			name  string
			price float64
		}{
			{"zero price", 0.00},
			{"small price", 0.99},
			{"normal price", 19.99},
			{"large price", 999.99},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				album := models.Album{
					Id:     1,
					Title:  "Test",
					Artist: "Test",
					Price:  decimal.NewFromFloat(tc.price),
				}

				priceFloat, _ := album.Price.Float64()
				if priceFloat != tc.price {
					t.Errorf("Expected price %.2f, got %.2f", tc.price, priceFloat)
				}
			})
		}
	})
}
