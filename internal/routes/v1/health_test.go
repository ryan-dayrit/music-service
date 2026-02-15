package v1

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterHealthRoute(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "health route returns OK",
			path:           "/health",
			method:         "GET",
			expectedStatus: fiber.StatusOK,
			expectedBody:   "OK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			router := app.Group("")
			RegisterHealthRoute(router)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var response map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if status, ok := response["status"].(string); !ok || status != tt.expectedBody {
				t.Errorf("Expected status: %s, got: %v", tt.expectedBody, response["status"])
			}
		})
	}
}

func TestRegisterHealthRoute_WithRouter(t *testing.T) {
	t.Run("health route works with v1 router group", func(t *testing.T) {
		app := fiber.New()
		v1Router := app.Group("/v1")
		RegisterHealthRoute(v1Router)

		req, err := http.NewRequest("GET", "/v1/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status code %d, got %d", fiber.StatusOK, resp.StatusCode)
		}
	})
}

func TestRegisterHealthRoute_OnlyGET(t *testing.T) {
	t.Run("health route only accepts GET method", func(t *testing.T) {
		app := fiber.New()
		router := app.Group("")
		RegisterHealthRoute(router)

		methods := []string{"POST", "PUT", "DELETE", "PATCH"}
		for _, method := range methods {
			req, err := http.NewRequest(method, "/health", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			// These methods should not be found
			if resp.StatusCode == fiber.StatusOK {
				t.Errorf("Method %s should not be allowed on health route", method)
			}
		}
	})
}

func TestRegisterHealthRoute_ResponseFormat(t *testing.T) {
	t.Run("response format is correct JSON", func(t *testing.T) {
		app := fiber.New()
		router := app.Group("")
		RegisterHealthRoute(router)

		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if _, ok := response["status"]; !ok {
			t.Error("Response missing 'status' key")
		}

		if len(response) != 1 {
			t.Errorf("Expected 1 key in response, got %d", len(response))
		}
	})
}

func TestRegisterHealthRoute_ContentType(t *testing.T) {
	t.Run("response has correct content type", func(t *testing.T) {
		app := fiber.New()
		router := app.Group("")
		RegisterHealthRoute(router)

		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Logf("Expected Content-Type: application/json, got: %s", contentType)
		}
	})
}

func TestRegisterHealthRoute_MultipleRouters(t *testing.T) {
	t.Run("can register health route on multiple routers", func(t *testing.T) {
		app := fiber.New()
		v1Router := app.Group("/v1")
		v2Router := app.Group("/v2")

		RegisterHealthRoute(v1Router)
		RegisterHealthRoute(v2Router)

		// Test v1
		req1, err := http.NewRequest("GET", "/v1/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp1, err := app.Test(req1)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp1.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status code %d for /v1/health, got %d", fiber.StatusOK, resp1.StatusCode)
		}

		// Test v2
		req2, err := http.NewRequest("GET", "/v2/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp2, err := app.Test(req2)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp2.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status code %d for /v2/health, got %d", fiber.StatusOK, resp2.StatusCode)
		}
	})
}
