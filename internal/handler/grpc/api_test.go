package grpc

import (
	"context"
	"errors"
	"testing"

	"music-service/gen/pb"
	"music-service/internal/models"

	"github.com/shopspring/decimal"
)

type MockRepository struct {
	ReadFunc func() ([]models.Album, error)
}

func (m *MockRepository) Read() ([]models.Album, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc()
	}
	return []models.Album{}, nil
}

func TestNewHandler(t *testing.T) {
	mockRepo := &MockRepository{}
	srv := NewHandler(mockRepo)

	if srv == nil {
		t.Fatal("Expected non-nil server, got nil")
	}

	var _ pb.MusicServiceServer = srv
}

func TestHandler_GetAlbumList_Success(t *testing.T) {
	mockRepo := &MockRepository{
		ReadFunc: func() ([]models.Album, error) {
			return []models.Album{
				{
					Id:     1,
					Title:  "Blue Train",
					Artist: "John Coltrane",
					Price:  decimal.NewFromFloat(56.99),
				},
				{
					Id:     2,
					Title:  "Jeru",
					Artist: "Gerry Mulligan",
					Price:  decimal.NewFromFloat(17.99),
				},
			}, nil
		},
	}

	srv := NewHandler(mockRepo)
	req := &pb.GetAlbumsRequest{}
	resp, err := srv.GetAlbumList(context.Background(), req)

	if err != nil {
		t.Fatalf("GetAlbumList() failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if len(resp.Albums) != 2 {
		t.Errorf("Expected 2 albums, got %d", len(resp.Albums))
	}

	if resp.Albums[0].Id != 1 {
		t.Errorf("Expected album[0] Id 1, got %d", resp.Albums[0].Id)
	}
	if resp.Albums[0].Title != "Blue Train" {
		t.Errorf("Expected album[0] Title 'Blue Train', got '%s'", resp.Albums[0].Title)
	}
	if resp.Albums[0].Artist != "John Coltrane" {
		t.Errorf("Expected album[0] Artist 'John Coltrane', got '%s'", resp.Albums[0].Artist)
	}
	if resp.Albums[0].Price != float32(56.99) {
		t.Errorf("Expected album[0] Price 56.99, got %f", resp.Albums[0].Price)
	}

	if resp.Albums[1].Id != 2 {
		t.Errorf("Expected album[1] Id 2, got %d", resp.Albums[1].Id)
	}
	if resp.Albums[1].Title != "Jeru" {
		t.Errorf("Expected album[1] Title 'Jeru', got '%s'", resp.Albums[1].Title)
	}
}

func TestHandler_GetAlbumList_EmptyResult(t *testing.T) {
	mockRepo := &MockRepository{
		ReadFunc: func() ([]models.Album, error) {
			return []models.Album{}, nil
		},
	}

	srv := NewHandler(mockRepo)
	req := &pb.GetAlbumsRequest{}
	resp, err := srv.GetAlbumList(context.Background(), req)

	if err != nil {
		t.Fatalf("GetAlbumList() failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if len(resp.Albums) != 0 {
		t.Errorf("Expected 0 albums, got %d", len(resp.Albums))
	}
}

func TestGetAlbumList_Function(t *testing.T) {
	tests := []struct {
		name          string
		mockAlbums    []models.Album
		mockError     error
		expectedCount int
		expectError   bool
	}{
		{
			name: "Single album",
			mockAlbums: []models.Album{
				{
					Id:     1,
					Title:  "Test Album",
					Artist: "Test Artist",
					Price:  decimal.NewFromFloat(9.99),
				},
			},
			mockError:     nil,
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Empty album list",
			mockAlbums:    []models.Album{},
			mockError:     nil,
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "Repository error",
			mockAlbums:    nil,
			mockError:     errors.New("database connection failed"),
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				ReadFunc: func() ([]models.Album, error) {
					return tt.mockAlbums, tt.mockError
				},
			}

			result, err := getAlbumList(mockRepo)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d albums, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestGetAlbumList_PriceConversion(t *testing.T) {
	mockRepo := &MockRepository{
		ReadFunc: func() ([]models.Album, error) {
			return []models.Album{
				{
					Id:     1,
					Title:  "Test",
					Artist: "Artist",
					Price:  decimal.NewFromFloat(12.345),
				},
			}, nil
		},
	}

	result, err := getAlbumList(mockRepo)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 album, got %d", len(result))
	}

	expectedPrice := float32(12.345)
	if result[0].Price != expectedPrice {
		t.Errorf("Expected price %f, got %f", expectedPrice, result[0].Price)
	}
}

func TestServer_GetAlbumList_ContextHandling(t *testing.T) {
	mockRepo := &MockRepository{
		ReadFunc: func() ([]models.Album, error) {
			return []models.Album{}, nil
		},
	}

	srv := NewHandler(mockRepo)

	req := &pb.GetAlbumsRequest{}

	// Test with valid contexts - should succeed
	_, err := srv.GetAlbumList(context.Background(), req)
	if err != nil {
		t.Errorf("GetAlbumList() with Background context failed: %v", err)
	}

	_, err = srv.GetAlbumList(context.TODO(), req)
	if err != nil {
		t.Errorf("GetAlbumList() with TODO context failed: %v", err)
	}

	// Test with canceled context - should return error
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = srv.GetAlbumList(ctx, req)
	if err == nil {
		t.Error("GetAlbumList() with canceled context should return error, got nil")
	}
	if err != context.Canceled {
		t.Errorf("GetAlbumList() with canceled context should return context.Canceled, got: %v", err)
	}
}

func BenchmarkServer_GetAlbumList(b *testing.B) {
	mockRepo := &MockRepository{
		ReadFunc: func() ([]models.Album, error) {
			return []models.Album{
				{Id: 1, Title: "Album 1", Artist: "Artist 1", Price: decimal.NewFromFloat(9.99)},
				{Id: 2, Title: "Album 2", Artist: "Artist 2", Price: decimal.NewFromFloat(19.99)},
				{Id: 3, Title: "Album 3", Artist: "Artist 3", Price: decimal.NewFromFloat(29.99)},
			}, nil
		},
	}

	srv := NewHandler(mockRepo)
	req := &pb.GetAlbumsRequest{}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = srv.GetAlbumList(ctx, req)
	}
}
