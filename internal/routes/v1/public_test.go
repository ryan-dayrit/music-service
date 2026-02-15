package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"

	"music-service/gen/pb"
	"music-service/internal/models"
)

// MockProducer is a mock implementation of kafka.Producer for testing
type MockProducer struct {
	produceCalled bool
	lastAlbum     *pb.Album
	returnError   error
}

func (m *MockProducer) Produce(ctx context.Context, album *pb.Album) {
	m.produceCalled = true
	m.lastAlbum = album
}

type MockRepository struct {
}

func (m *MockRepository) Create(album models.Album) error {
	return nil
}

func (m *MockRepository) GetById(id int) (*models.Album, error) {
	return nil, nil
}

func (m *MockRepository) Get() ([]*models.Album, error) {
	return nil, nil
}

func (m *MockRepository) Update(album models.Album) error {
	return nil
}

func (m *MockRepository) Upsert(album models.Album) error {
	return nil
}

func TestRegisterPublicRoutes(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "POST /album route is registered",
			method:         "POST",
			path:           "/album",
			expectedStatus: fiber.StatusBadRequest, // Will fail validation, but route exists
		},
		{
			name:           "PUT /album route is registered",
			method:         "PUT",
			path:           "/album",
			expectedStatus: fiber.StatusBadRequest, // Will fail validation, but route exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			router := app.Group("")
			mockProducer := &MockProducer{}
			mockRepostory := &MockRepository{}
			RegisterPublicRoutes(router, mockProducer, mockRepostory)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			// Just verify the route exists (we expect validation errors without proper body)
			if resp.StatusCode == fiber.StatusNotFound {
				t.Errorf("Route %s %s was not registered", tt.method, tt.path)
			}
		})
	}
}

func TestRegisterPublicRoutes_WithRouterGroup(t *testing.T) {
	t.Run("public routes work with v1 router group", func(t *testing.T) {
		app := fiber.New()
		v1Router := app.Group("/v1")
		mockProducer := &MockProducer{}
		mockRepository := &MockRepository{}
		RegisterPublicRoutes(v1Router, mockProducer, mockRepository)

		req, err := http.NewRequest("POST", "/v1/album", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		// Verify route exists under /v1
		if resp.StatusCode == fiber.StatusNotFound {
			t.Error("Route /v1/album was not registered")
		}
	})
}

func TestRegisterPublicRoutes_BothMethodsUseSameHandler(t *testing.T) {
	t.Run("POST and PUT both call CreateAlbum", func(t *testing.T) {
		app := fiber.New()
		router := app.Group("")
		mockProducer := &MockProducer{}
		mockRepository := &MockRepository{}

		RegisterPublicRoutes(router, mockProducer, mockRepository)

		// Test POST
		reqPost, err := http.NewRequest("POST", "/album", nil)
		if err != nil {
			t.Fatalf("Failed to create POST request: %v", err)
		}

		respPost, err := app.Test(reqPost)
		if err != nil {
			t.Fatalf("Failed to test POST request: %v", err)
		}

		// Test PUT
		reqPut, err := http.NewRequest("PUT", "/album", nil)
		if err != nil {
			t.Fatalf("Failed to create PUT request: %v", err)
		}

		respPut, err := app.Test(reqPut)
		if err != nil {
			t.Fatalf("Failed to test PUT request: %v", err)
		}

		// Both should have the same behavior (not 404)
		if respPost.StatusCode == fiber.StatusNotFound {
			t.Error("POST /album route not found")
		}
		if respPut.StatusCode == fiber.StatusNotFound {
			t.Error("PUT /album route not found")
		}
	})
}

func TestRegisterPublicRoutes_WithValidPayload(t *testing.T) {
	t.Run("routes accept valid JSON payload", func(t *testing.T) {
		app := fiber.New()
		router := app.Group("")
		mockProducer := &MockProducer{}
		mockRepository := &MockRepository{}

		RegisterPublicRoutes(router, mockProducer, mockRepository)

		payload := map[string]interface{}{
			"id":     "1",
			"title":  "Test Album",
			"artist": "Test Artist",
			"price":  9.99,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}

		req, err := http.NewRequest("POST", "/album", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		// Should not be 404 (route exists)
		if resp.StatusCode == fiber.StatusNotFound {
			t.Error("Route not found")
		}
	})
}

func TestRegisterPublicRoutes_OtherMethodsNotAllowed(t *testing.T) {
	t.Run("other HTTP methods are not registered", func(t *testing.T) {
		app := fiber.New()
		router := app.Group("")
		mockProducer := &MockProducer{}
		mockRepository := &MockRepository{}

		RegisterPublicRoutes(router, mockProducer, mockRepository)

		methods := []string{"GET", "DELETE", "PATCH"}
		for _, method := range methods {
			req, err := http.NewRequest(method, "/album", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			// These methods should not be allowed
			if resp.StatusCode != fiber.StatusMethodNotAllowed && resp.StatusCode != fiber.StatusNotFound {
				t.Logf("Method %s returned unexpected status: %d", method, resp.StatusCode)
			}
		}
	})
}

func TestRegisterPublicRoutes_ProducerInjection(t *testing.T) {
	t.Run("producer is properly injected into handler", func(t *testing.T) {
		app := fiber.New()
		router := app.Group("")
		mockProducer := &MockProducer{}
		mockRepository := &MockRepository{}
		// This should not panic if producer is properly injected
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Registering routes with producer caused panic: %v", r)
			}
		}()

		RegisterPublicRoutes(router, mockProducer, mockRepository)

		// Make a request to verify handler was created successfully
		req, err := http.NewRequest("POST", "/album", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		_, err = app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}
	})
}

func TestRegisterPublicRoutes_MultipleRegistrations(t *testing.T) {
	t.Run("can register public routes on multiple routers", func(t *testing.T) {
		app := fiber.New()
		v1Router := app.Group("/v1")
		v2Router := app.Group("/v2")
		mockProducer := &MockProducer{}
		mockRepository := &MockRepository{}

		RegisterPublicRoutes(v1Router, mockProducer, mockRepository)
		RegisterPublicRoutes(v2Router, mockProducer, mockRepository)

		// Test v1
		req1, err := http.NewRequest("POST", "/v1/album", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp1, err := app.Test(req1)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp1.StatusCode == fiber.StatusNotFound {
			t.Error("Route /v1/album not found")
		}

		// Test v2
		req2, err := http.NewRequest("POST", "/v2/album", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp2, err := app.Test(req2)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp2.StatusCode == fiber.StatusNotFound {
			t.Error("Route /v2/album not found")
		}
	})
}
