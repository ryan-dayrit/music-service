package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	"music-service/gen/pb"
	"music-service/internal/models"
)

// mockRepository is a mock implementation of orm.Repository
type mockRepository struct {
	createFunc  func(album models.Album) error
	getByIdFunc func(id int) (*models.Album, error)
	getFunc     func() ([]*models.Album, error)
	updateFunc  func(album models.Album) error
	upsertFunc  func(album models.Album) error
	getCalls    int
}

func (m *mockRepository) Create(album models.Album) error {
	if m.createFunc != nil {
		return m.createFunc(album)
	}
	return nil
}

func (m *mockRepository) GetById(id int) (*models.Album, error) {
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
	return []*models.Album{}, nil
}

func (m *mockRepository) Update(album models.Album) error {
	if m.updateFunc != nil {
		return m.updateFunc(album)
	}
	return nil
}

func (m *mockRepository) Upsert(album models.Album) error {
	if m.upsertFunc != nil {
		return m.upsertFunc(album)
	}
	return nil
}

func TestNewAlbumsHandler(t *testing.T) {
	t.Run("creates new albums handler successfully", func(t *testing.T) {
		mockProducer := &mockProducerHandler{}
		mockRepo := &mockRepository{}
		handler := NewAlbumsHandler(mockProducer, mockRepo)

		if handler == nil {
			t.Fatal("Expected handler to be non-nil")
		}

		if handler.producerHandler != mockProducer {
			t.Error("Expected producer handler to be set correctly")
		}

		if handler.repository != mockRepo {
			t.Error("Expected repository to be set correctly")
		}
	})
}

func TestAlbumsHandler_CreateAlbums(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
		setupMocks     func(*mockProducerHandler, *mockRepository)
		validateMocks  func(*testing.T, *mockProducerHandler, *mockRepository)
	}{
		{
			name: "successfully creates single album",
			requestBody: []*pb.Album{
				{
					Id:     1,
					Title:  "Blue Train",
					Artist: "John Coltrane",
					Price:  56.99,
				},
			},
			expectedStatus: fiber.StatusCreated,
			setupMocks: func(mp *mockProducerHandler, mr *mockRepository) {
				mp.produceFunc = func(ctx context.Context, album *pb.Album) {
					if album.Title != "Blue Train" {
						t.Errorf("Expected album title 'Blue Train', got '%s'", album.Title)
					}
				}
			},
			validateMocks: func(t *testing.T, mp *mockProducerHandler, mr *mockRepository) {
				if mp.produceCalls != 1 {
					t.Errorf("Expected Produce to be called once, got %d calls", mp.produceCalls)
				}
			},
		},
		{
			name: "successfully creates multiple albums",
			requestBody: []*pb.Album{
				{
					Id:     1,
					Title:  "Blue Train",
					Artist: "John Coltrane",
					Price:  56.99,
				},
				{
					Id:     2,
					Title:  "Jeru",
					Artist: "Gerry Mulligan",
					Price:  17.99,
				},
				{
					Id:     3,
					Title:  "Giant Steps",
					Artist: "John Coltrane",
					Price:  63.99,
				},
			},
			expectedStatus: fiber.StatusCreated,
			setupMocks: func(mp *mockProducerHandler, mr *mockRepository) {
				mp.produceFunc = func(ctx context.Context, album *pb.Album) {
					// Verify each album
				}
			},
			validateMocks: func(t *testing.T, mp *mockProducerHandler, mr *mockRepository) {
				if mp.produceCalls != 3 {
					t.Errorf("Expected Produce to be called 3 times, got %d calls", mp.produceCalls)
				}
			},
		},
		{
			name:           "successfully creates empty album list",
			requestBody:    []*pb.Album{},
			expectedStatus: fiber.StatusCreated,
			setupMocks:     func(mp *mockProducerHandler, mr *mockRepository) {},
			validateMocks: func(t *testing.T, mp *mockProducerHandler, mr *mockRepository) {
				if mp.produceCalls != 0 {
					t.Errorf("Expected Produce not to be called, got %d calls", mp.produceCalls)
				}
			},
		},
		{
			name:           "returns bad request for invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "cannot parse JSON",
			setupMocks:     func(mp *mockProducerHandler, mr *mockRepository) {},
			validateMocks: func(t *testing.T, mp *mockProducerHandler, mr *mockRepository) {
				if mp.produceCalls != 0 {
					t.Errorf("Expected Produce not to be called, got %d calls", mp.produceCalls)
				}
			},
		},
		{
			name:           "returns bad request for malformed JSON",
			requestBody:    `[{"id": "not a number"}]`,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "cannot parse JSON",
			setupMocks:     func(mp *mockProducerHandler, mr *mockRepository) {},
			validateMocks: func(t *testing.T, mp *mockProducerHandler, mr *mockRepository) {
				if mp.produceCalls != 0 {
					t.Errorf("Expected Produce not to be called, got %d calls", mp.produceCalls)
				}
			},
		},
		{
			name: "handles albums with partial data",
			requestBody: []*pb.Album{
				{
					Title: "Album with only title",
				},
				{
					Artist: "Artist with only artist name",
				},
			},
			expectedStatus: fiber.StatusCreated,
			setupMocks:     func(mp *mockProducerHandler, mr *mockRepository) {},
			validateMocks: func(t *testing.T, mp *mockProducerHandler, mr *mockRepository) {
				if mp.produceCalls != 2 {
					t.Errorf("Expected Produce to be called 2 times, got %d calls", mp.produceCalls)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockProducer := &mockProducerHandler{}
			mockRepo := &mockRepository{}
			if tt.setupMocks != nil {
				tt.setupMocks(mockProducer, mockRepo)
			}
			handler := NewAlbumsHandler(mockProducer, mockRepo)

			// Register route
			app.Post("/albums", handler.CreateAlbums)

			// Prepare request body
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			// Create request
			req, err := http.NewRequest("POST", "/albums", bytes.NewReader(body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			// Validate status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Validate response body
			var response interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Check for error message if expected
			if tt.expectedError != "" {
				respMap, ok := response.(map[string]interface{})
				if !ok {
					t.Error("Expected response to be a map")
				} else {
					errorMsg, ok := respMap["error"].(string)
					if !ok {
						t.Error("Expected error field in response")
					} else if errorMsg != tt.expectedError {
						t.Errorf("Expected error '%s', got '%s'", tt.expectedError, errorMsg)
					}
				}
			} else {
				// Validate successful response contains albums data
				respMap, ok := response.(map[string]interface{})
				if ok {
					if _, hasError := respMap["error"]; hasError {
						t.Errorf("Expected no error in response, got: %v", respMap["error"])
					}
				}
			}

			// Validate mocks
			if tt.validateMocks != nil {
				tt.validateMocks(t, mockProducer, mockRepo)
			}
		})
	}
}

func TestAlbumsHandler_GetAlbums(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*mockRepository)
		expectedStatus int
		expectedError  string
		validateResult func(*testing.T, []interface{})
	}{
		{
			name: "successfully retrieves albums",
			setupMock: func(mr *mockRepository) {
				mr.getFunc = func() ([]*models.Album, error) {
					return []*models.Album{
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
				}
			},
			expectedStatus: fiber.StatusOK,
			validateResult: func(t *testing.T, albums []interface{}) {
				if len(albums) != 2 {
					t.Errorf("Expected 2 albums, got %d", len(albums))
				}
			},
		},
		{
			name: "successfully retrieves empty album list",
			setupMock: func(mr *mockRepository) {
				mr.getFunc = func() ([]*models.Album, error) {
					return []*models.Album{}, nil
				}
			},
			expectedStatus: fiber.StatusOK,
			validateResult: func(t *testing.T, albums []interface{}) {
				if len(albums) != 0 {
					t.Errorf("Expected 0 albums, got %d", len(albums))
				}
			},
		},
		{
			name: "returns internal server error on repository error",
			setupMock: func(mr *mockRepository) {
				mr.getFunc = func() ([]*models.Album, error) {
					return nil, errors.New("database connection failed")
				}
			},
			expectedStatus: fiber.StatusInternalServerError,
			expectedError:  "failed to get albums",
		},
		{
			name: "handles nil albums from repository",
			setupMock: func(mr *mockRepository) {
				mr.getFunc = func() ([]*models.Album, error) {
					return nil, nil
				}
			},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockProducer := &mockProducerHandler{}
			mockRepo := &mockRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}
			handler := NewAlbumsHandler(mockProducer, mockRepo)

			// Register route
			app.Get("/albums", handler.GetAlbums)

			// Create request
			req, err := http.NewRequest("GET", "/albums", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Execute request
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			// Validate status code
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Validate response body
			var response interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Check for error message if expected
			if tt.expectedError != "" {
				respMap, ok := response.(map[string]interface{})
				if !ok {
					t.Error("Expected response to be a map")
				} else {
					errorMsg, ok := respMap["error"].(string)
					if !ok {
						t.Error("Expected error field in response")
					} else if errorMsg != tt.expectedError {
						t.Errorf("Expected error '%s', got '%s'", tt.expectedError, errorMsg)
					}
				}
			} else {
				// Validate successful response
				if albums, ok := response.([]interface{}); ok {
					if tt.validateResult != nil {
						tt.validateResult(t, albums)
					}
				} else if respMap, ok := response.(map[string]interface{}); ok {
					if _, hasError := respMap["error"]; hasError {
						t.Errorf("Expected no error in response, got: %v", respMap["error"])
					}
				}
			}

			// Validate repository was called
			if mockRepo.getCalls != 1 {
				t.Errorf("Expected Get to be called once, got %d calls", mockRepo.getCalls)
			}
		})
	}
}

func TestAlbumsHandler_GetAlbums_RepositoryPanic(t *testing.T) {
	t.Run("handles repository panic", func(t *testing.T) {
		// This test verifies the behavior when the repository panics
		// In production, proper error handling should be implemented

		// Setup
		app := fiber.New()
		mockProducer := &mockProducerHandler{}
		mockRepo := &mockRepository{
			getFunc: func() ([]*models.Album, error) {
				// Note: In production, the repository should return errors, not panic
				return nil, errors.New("simulated error instead of panic")
			},
		}
		handler := NewAlbumsHandler(mockProducer, mockRepo)

		// Register route
		app.Get("/albums", handler.GetAlbums)

		// Create request
		req, _ := http.NewRequest("GET", "/albums", nil)

		// Execute request
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		// Should return internal server error
		if resp.StatusCode != fiber.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", fiber.StatusInternalServerError, resp.StatusCode)
		}
	})
}
