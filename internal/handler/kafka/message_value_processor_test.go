package kafka

import (
	"errors"
	"testing"

	"google.golang.org/protobuf/proto"

	"music-service/gen/pb"
	"music-service/internal/models"
)

// mockRepository is a mock implementation of orm.Repository
type mockRepository struct {
	createFunc   func(album models.Album) error
	getByIdFunc  func(id int) (*models.Album, error)
	updateFunc   func(album models.Album) error
	getFunc      func() ([]*models.Album, error)
	upsertFunc   func(album models.Album) error
	createCalls  int
	getByIdCalls int
	updateCalls  int
	getCalls     int
	upsertCalls  int
}

func (m *mockRepository) Create(album models.Album) error {
	m.createCalls++
	if m.createFunc != nil {
		return m.createFunc(album)
	}
	return nil
}

func (m *mockRepository) GetById(id int) (*models.Album, error) {
	m.getByIdCalls++
	if m.getByIdFunc != nil {
		return m.getByIdFunc(id)
	}
	return nil, nil
}

func (m *mockRepository) Get() ([]*models.Album, error) {
	m.getCalls++
	if m.getFunc != nil {
		return m.getFunc()
	}
	return nil, nil
}

func (m *mockRepository) Update(album models.Album) error {
	m.updateCalls++
	if m.updateFunc != nil {
		return m.updateFunc(album)
	}
	return nil
}

func (m *mockRepository) Upsert(album models.Album) error {
	m.upsertCalls++
	if m.upsertFunc != nil {
		return m.upsertFunc(album)
	}
	return nil
}

func TestNewMessageValueProcessor(t *testing.T) {
	t.Run("creates new message value processor successfully", func(t *testing.T) {
		mockRepo := &mockRepository{}
		processor := NewMessageValueProcessor(mockRepo)

		if processor == nil {
			t.Fatal("Expected processor to be non-nil")
		}

		if processor.repository != mockRepo {
			t.Error("Expected repository to be set correctly")
		}
	})
}

func TestMessageValueProcessor_ProcessMessageValue_CreateNewAlbum(t *testing.T) {
	t.Run("creates new album when not found in database", func(t *testing.T) {
		mockRepo := &mockRepository{}
		processor := NewMessageValueProcessor(mockRepo)

		// Setup mock to return "no rows in result set" error for GetById
		mockRepo.getByIdFunc = func(id int) (*models.Album, error) {
			if id != 1 {
				t.Errorf("Expected id 1, got %d", id)
			}
			return nil, errors.New("pg: no rows in result set")
		}

		// Setup mock to verify Create is called with correct data
		mockRepo.createFunc = func(album models.Album) error {
			if album.Id != 1 {
				t.Errorf("Expected album id 1, got %d", album.Id)
			}
			if album.Title != "Blue Train" {
				t.Errorf("Expected album title 'Blue Train', got '%s'", album.Title)
			}
			if album.Artist != "John Coltrane" {
				t.Errorf("Expected artist 'John Coltrane', got '%s'", album.Artist)
			}
			return nil
		}

		// Create protobuf album and marshal it
		protoAlbum := &pb.Album{
			Id:     1,
			Title:  "Blue Train",
			Artist: "John Coltrane",
			Price:  56.99,
		}
		messageValue, err := proto.Marshal(protoAlbum)
		if err != nil {
			t.Fatalf("Failed to marshal proto album: %v", err)
		}

		// Process the message value
		processor.ProcessMessageValue(messageValue)

		// Verify GetById was called
		if mockRepo.getByIdCalls != 1 {
			t.Errorf("Expected GetById to be called once, got %d calls", mockRepo.getByIdCalls)
		}

		// Verify Create was called
		if mockRepo.createCalls != 1 {
			t.Errorf("Expected Create to be called once, got %d calls", mockRepo.createCalls)
		}

		// Verify Update was not called
		if mockRepo.updateCalls != 0 {
			t.Errorf("Expected Update not to be called, got %d calls", mockRepo.updateCalls)
		}
	})
}

func TestMessageValueProcessor_ProcessMessageValue_UpdateExistingAlbum(t *testing.T) {
	t.Run("updates existing album when found in database", func(t *testing.T) {
		mockRepo := &mockRepository{}
		processor := NewMessageValueProcessor(mockRepo)

		// Setup mock to return an existing album for GetById
		mockRepo.getByIdFunc = func(id int) (*models.Album, error) {
			if id != 2 {
				t.Errorf("Expected id 2, got %d", id)
			}
			return &models.Album{
				Id:     2,
				Title:  "Old Title",
				Artist: "Old Artist",
			}, nil
		}

		// Setup mock to verify Update is called with correct data
		mockRepo.updateFunc = func(album models.Album) error {
			if album.Id != 2 {
				t.Errorf("Expected album id 2, got %d", album.Id)
			}
			if album.Title != "Jeru" {
				t.Errorf("Expected album title 'Jeru', got '%s'", album.Title)
			}
			if album.Artist != "Gerry Mulligan" {
				t.Errorf("Expected artist 'Gerry Mulligan', got '%s'", album.Artist)
			}
			return nil
		}

		// Create protobuf album and marshal it
		protoAlbum := &pb.Album{
			Id:     2,
			Title:  "Jeru",
			Artist: "Gerry Mulligan",
			Price:  17.99,
		}
		messageValue, err := proto.Marshal(protoAlbum)
		if err != nil {
			t.Fatalf("Failed to marshal proto album: %v", err)
		}

		// Process the message value
		processor.ProcessMessageValue(messageValue)

		// Verify GetById was called
		if mockRepo.getByIdCalls != 1 {
			t.Errorf("Expected GetById to be called once, got %d calls", mockRepo.getByIdCalls)
		}

		// Verify Update was called
		if mockRepo.updateCalls != 1 {
			t.Errorf("Expected Update to be called once, got %d calls", mockRepo.updateCalls)
		}

		// Verify Create was not called
		if mockRepo.createCalls != 0 {
			t.Errorf("Expected Create not to be called, got %d calls", mockRepo.createCalls)
		}
	})
}

func TestMessageValueProcessor_ProcessMessageValue_WithZeroValues(t *testing.T) {
	t.Run("handles album with zero values", func(t *testing.T) {
		mockRepo := &mockRepository{}
		processor := NewMessageValueProcessor(mockRepo)

		// Setup mock to return "no rows in result set" error for GetById
		mockRepo.getByIdFunc = func(id int) (*models.Album, error) {
			return nil, errors.New("pg: no rows in result set")
		}

		// Setup mock to verify Create is called
		mockRepo.createFunc = func(album models.Album) error {
			if album.Id != 0 {
				t.Errorf("Expected album id 0, got %d", album.Id)
			}
			if album.Title != "" {
				t.Errorf("Expected empty title, got '%s'", album.Title)
			}
			if album.Artist != "" {
				t.Errorf("Expected empty artist, got '%s'", album.Artist)
			}
			return nil
		}

		// Create protobuf album with zero values and marshal it
		protoAlbum := &pb.Album{
			Id:     0,
			Title:  "",
			Artist: "",
			Price:  0,
		}
		messageValue, err := proto.Marshal(protoAlbum)
		if err != nil {
			t.Fatalf("Failed to marshal proto album: %v", err)
		}

		// Process the message value
		processor.ProcessMessageValue(messageValue)

		// Verify GetById was called
		if mockRepo.getByIdCalls != 1 {
			t.Errorf("Expected GetById to be called once, got %d calls", mockRepo.getByIdCalls)
		}

		// Verify Create was called
		if mockRepo.createCalls != 1 {
			t.Errorf("Expected Create to be called once, got %d calls", mockRepo.createCalls)
		}
	})
}

func TestMessageValueProcessor_ProcessMessageValue_PriceGeneration(t *testing.T) {
	t.Run("generates random price for album", func(t *testing.T) {
		mockRepo := &mockRepository{}
		processor := NewMessageValueProcessor(mockRepo)

		// Setup mock to return "no rows in result set" error for GetById
		mockRepo.getByIdFunc = func(id int) (*models.Album, error) {
			return nil, errors.New("pg: no rows in result set")
		}

		// Setup mock to capture the price
		var capturedPrice float64
		mockRepo.createFunc = func(album models.Album) error {
			priceFloat, _ := album.Price.Float64()
			capturedPrice = priceFloat
			return nil
		}

		// Create protobuf album and marshal it
		protoAlbum := &pb.Album{
			Id:     1,
			Title:  "Test Album",
			Artist: "Test Artist",
			Price:  100.00, // This will be replaced with random price
		}
		messageValue, err := proto.Marshal(protoAlbum)
		if err != nil {
			t.Fatalf("Failed to marshal proto album: %v", err)
		}

		// Process the message value
		processor.ProcessMessageValue(messageValue)

		// Verify price is within range [0, 1)
		if capturedPrice < 0 || capturedPrice >= 1 {
			t.Errorf("Expected price to be in range [0, 1), got %f", capturedPrice)
		}
	})
}

func TestMessageValueProcessor_ProcessMessageValue_MultipleAlbums(t *testing.T) {
	t.Run("processes multiple albums correctly", func(t *testing.T) {
		mockRepo := &mockRepository{}
		processor := NewMessageValueProcessor(mockRepo)

		albums := []*pb.Album{
			{Id: 1, Title: "Album 1", Artist: "Artist 1", Price: 10.99},
			{Id: 2, Title: "Album 2", Artist: "Artist 2", Price: 20.99},
			{Id: 3, Title: "Album 3", Artist: "Artist 3", Price: 30.99},
		}

		// Track which albums were created
		createdAlbums := make(map[int]bool)

		// Setup mock to return "no rows in result set" for all albums
		mockRepo.getByIdFunc = func(id int) (*models.Album, error) {
			return nil, errors.New("pg: no rows in result set")
		}

		// Setup mock to track created albums
		mockRepo.createFunc = func(album models.Album) error {
			createdAlbums[album.Id] = true
			return nil
		}

		// Process each album
		for _, protoAlbum := range albums {
			messageValue, err := proto.Marshal(protoAlbum)
			if err != nil {
				t.Fatalf("Failed to marshal proto album: %v", err)
			}
			processor.ProcessMessageValue(messageValue)
		}

		// Verify all albums were created
		if len(createdAlbums) != 3 {
			t.Errorf("Expected 3 albums to be created, got %d", len(createdAlbums))
		}

		for _, album := range albums {
			if !createdAlbums[int(album.Id)] {
				t.Errorf("Expected album %d to be created", album.Id)
			}
		}

		// Verify GetById and Create were called 3 times each
		if mockRepo.getByIdCalls != 3 {
			t.Errorf("Expected GetById to be called 3 times, got %d calls", mockRepo.getByIdCalls)
		}
		if mockRepo.createCalls != 3 {
			t.Errorf("Expected Create to be called 3 times, got %d calls", mockRepo.createCalls)
		}
	})
}
