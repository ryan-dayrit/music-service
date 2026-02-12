package models

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestAlbum_StructFields(t *testing.T) {
	price := decimal.NewFromFloat(9.99)
	album := Album{
		Id:     1,
		Title:  "Blue Train",
		Artist: "John Coltrane",
		Price:  price,
	}

	if album.Id != 1 {
		t.Errorf("Expected Id 1, got %d", album.Id)
	}
	if album.Title != "Blue Train" {
		t.Errorf("Expected Title 'Blue Train', got '%s'", album.Title)
	}
	if album.Artist != "John Coltrane" {
		t.Errorf("Expected Artist 'John Coltrane', got '%s'", album.Artist)
	}
	if !album.Price.Equal(price) {
		t.Errorf("Expected Price %v, got %v", price, album.Price)
	}
}

func TestAlbum_ZeroValues(t *testing.T) {
	var album Album

	if album.Id != 0 {
		t.Errorf("Expected zero Id, got %d", album.Id)
	}
	if album.Title != "" {
		t.Errorf("Expected empty Title, got '%s'", album.Title)
	}
	if album.Artist != "" {
		t.Errorf("Expected empty Artist, got '%s'", album.Artist)
	}
	if !album.Price.IsZero() {
		t.Errorf("Expected zero Price, got %v", album.Price)
	}
}

func TestAlbum_DecimalPrice(t *testing.T) {
	tests := []struct {
		name     string
		price    decimal.Decimal
		expected string
	}{
		{
			name:     "Typical price",
			price:    decimal.NewFromFloat(19.99),
			expected: "19.99",
		},
		{
			name:     "Zero price",
			price:    decimal.NewFromFloat(0),
			expected: "0",
		},
		{
			name:     "High precision price",
			price:    decimal.NewFromFloat(12.345),
			expected: "12.345",
		},
		{
			name:     "Large price",
			price:    decimal.NewFromFloat(999.99),
			expected: "999.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			album := Album{
				Id:     1,
				Title:  "Test Album",
				Artist: "Test Artist",
				Price:  tt.price,
			}

			if album.Price.String() != tt.expected {
				t.Errorf("Expected price string %s, got %s", tt.expected, album.Price.String())
			}
		})
	}
}

func TestAlbum_MultipleInstances(t *testing.T) {
	album1 := Album{
		Id:     1,
		Title:  "Album 1",
		Artist: "Artist 1",
		Price:  decimal.NewFromFloat(10.00),
	}

	album2 := Album{
		Id:     2,
		Title:  "Album 2",
		Artist: "Artist 2",
		Price:  decimal.NewFromFloat(20.00),
	}

	if album1.Id == album2.Id {
		t.Error("Album instances should have different IDs")
	}
	if album1.Title == album2.Title {
		t.Error("Album instances should have different titles")
	}
	if album1.Price.Equal(album2.Price) {
		t.Error("Album instances should have different prices")
	}
}
