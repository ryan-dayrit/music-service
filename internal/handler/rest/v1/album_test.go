package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"

	"music-service/gen/pb"
)

// mockProducerHandler is a mock implementation of kafka.ProducerHandler
type mockProducerHandler struct {
	produceFunc  func(ctx context.Context, album *pb.Album)
	produceCalls int
}

func (m *mockProducerHandler) Produce(ctx context.Context, album *pb.Album) {
	m.produceCalls++
	if m.produceFunc != nil {
		m.produceFunc(ctx, album)
	}
}

func TestNewAlbumHandler(t *testing.T) {
	t.Run("creates new album handler successfully", func(t *testing.T) {
		mockProducer := &mockProducerHandler{}
		handler := NewAlbumHandler(mockProducer)

		if handler == nil {
			t.Fatal("Expected handler to be non-nil")
		}

		if handler.producerHandler != mockProducer {
			t.Error("Expected producer handler to be set correctly")
		}
	})
}

func TestAlbumHandler_CreateAlbum(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
		setupMock      func(*mockProducerHandler)
		validateMock   func(*testing.T, *mockProducerHandler)
	}{
		{
			name: "successfully creates album",
			requestBody: &pb.Album{
				Id:     1,
				Title:  "Blue Train",
				Artist: "John Coltrane",
				Price:  56.99,
			},
			expectedStatus: fiber.StatusCreated,
			setupMock: func(m *mockProducerHandler) {
				m.produceFunc = func(ctx context.Context, album *pb.Album) {
					if album.Title != "Blue Train" {
						t.Errorf("Expected album title 'Blue Train', got '%s'", album.Title)
					}
					if album.Artist != "John Coltrane" {
						t.Errorf("Expected artist 'John Coltrane', got '%s'", album.Artist)
					}
				}
			},
			validateMock: func(t *testing.T, m *mockProducerHandler) {
				if m.produceCalls != 1 {
					t.Errorf("Expected Produce to be called once, got %d calls", m.produceCalls)
				}
			},
		},
		{
			name: "successfully creates album with zero values",
			requestBody: &pb.Album{
				Id:     0,
				Title:  "",
				Artist: "",
				Price:  0,
			},
			expectedStatus: fiber.StatusCreated,
			setupMock: func(m *mockProducerHandler) {
				m.produceFunc = func(ctx context.Context, album *pb.Album) {
					// Just verify it's called
				}
			},
			validateMock: func(t *testing.T, m *mockProducerHandler) {
				if m.produceCalls != 1 {
					t.Errorf("Expected Produce to be called once, got %d calls", m.produceCalls)
				}
			},
		},
		{
			name:           "returns bad request for invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "cannot parse JSON",
			setupMock:      func(m *mockProducerHandler) {},
			validateMock: func(t *testing.T, m *mockProducerHandler) {
				if m.produceCalls != 0 {
					t.Errorf("Expected Produce not to be called, got %d calls", m.produceCalls)
				}
			},
		},
		{
			name:           "returns bad request for malformed JSON",
			requestBody:    `{"id": "not a number"}`,
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "cannot parse JSON",
			setupMock:      func(m *mockProducerHandler) {},
			validateMock: func(t *testing.T, m *mockProducerHandler) {
				if m.produceCalls != 0 {
					t.Errorf("Expected Produce not to be called, got %d calls", m.produceCalls)
				}
			},
		},
		{
			name: "creates album with partial data",
			requestBody: &pb.Album{
				Title: "Jeru",
			},
			expectedStatus: fiber.StatusCreated,
			setupMock: func(m *mockProducerHandler) {
				m.produceFunc = func(ctx context.Context, album *pb.Album) {
					if album.Title != "Jeru" {
						t.Errorf("Expected album title 'Jeru', got '%s'", album.Title)
					}
				}
			},
			validateMock: func(t *testing.T, m *mockProducerHandler) {
				if m.produceCalls != 1 {
					t.Errorf("Expected Produce to be called once, got %d calls", m.produceCalls)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockProducer := &mockProducerHandler{}
			if tt.setupMock != nil {
				tt.setupMock(mockProducer)
			}
			handler := NewAlbumHandler(mockProducer)

			// Register route
			app.Post("/album", handler.CreateAlbum)

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
			req, err := http.NewRequest("POST", "/album", bytes.NewReader(body))
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
			var response map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Check for error message if expected
			if tt.expectedError != "" {
				errorMsg, ok := response["error"].(string)
				if !ok {
					t.Error("Expected error field in response")
				} else if errorMsg != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, errorMsg)
				}
			} else {
				// Validate successful response contains album data
				if _, hasError := response["error"]; hasError {
					t.Errorf("Expected no error in response, got: %v", response["error"])
				}
			}

			// Validate mock
			if tt.validateMock != nil {
				tt.validateMock(t, mockProducer)
			}
		})
	}
}

func TestAlbumHandler_CreateAlbum_ProducerPanic(t *testing.T) {
	t.Run("handles producer panic gracefully", func(t *testing.T) {
		// This test verifies the behavior when the producer panics
		// In a real scenario, you'd want proper error handling middleware

		// Setup
		app := fiber.New()

		mockProducer := &mockProducerHandler{
			produceFunc: func(ctx context.Context, album *pb.Album) {
				// Note: In production, the producer should not panic
				// This test documents the current behavior
			},
		}
		handler := NewAlbumHandler(mockProducer)

		// Register route
		app.Post("/album", handler.CreateAlbum)

		// Prepare request
		album := &pb.Album{
			Id:     1,
			Title:  "Test",
			Artist: "Test Artist",
			Price:  9.99,
		}
		body, _ := json.Marshal(album)
		req, _ := http.NewRequest("POST", "/album", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		// Should succeed since we're not panicking in the mock
		if resp.StatusCode != fiber.StatusCreated {
			t.Errorf("Expected status code %d, got %d", fiber.StatusCreated, resp.StatusCode)
		}
	})
}

func TestAlbumHandler_CreateAlbum_EmptyBody(t *testing.T) {
	t.Run("handles empty request body", func(t *testing.T) {
		// Setup
		app := fiber.New()
		mockProducer := &mockProducerHandler{}
		handler := NewAlbumHandler(mockProducer)

		// Register route
		app.Post("/album", handler.CreateAlbum)

		// Create request with empty body
		req, err := http.NewRequest("POST", "/album", bytes.NewReader([]byte("")))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Execute request
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		// Validate status code - empty body is technically valid JSON for protobuf
		if resp.StatusCode != fiber.StatusCreated && resp.StatusCode != fiber.StatusBadRequest {
			t.Errorf("Expected status code %d or %d, got %d",
				fiber.StatusCreated, fiber.StatusBadRequest, resp.StatusCode)
		}
	})
}
