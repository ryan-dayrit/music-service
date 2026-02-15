package routes

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterNotFoundRoute(t *testing.T) {
	t.Run("not found for undefined route", func(t *testing.T) {
		app := fiber.New()
		RegisterNotFoundRoute(app)

		req, err := http.NewRequest("GET", "/undefined", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp.StatusCode != fiber.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", fiber.StatusNotFound, resp.StatusCode)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if errorVal, ok := response["error"].(bool); !ok || errorVal != true {
			t.Errorf("Expected error: true, got: %v", response["error"])
		}

		if msgVal, ok := response["msg"].(string); !ok || msgVal != "endpoint is not found" {
			t.Errorf("Expected msg: 'endpoint is not found', got: %s", response["msg"])
		}
	})

	t.Run("not found for POST to undefined route", func(t *testing.T) {
		app := fiber.New()
		RegisterNotFoundRoute(app)

		req, err := http.NewRequest("POST", "/api/undefined", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp.StatusCode != fiber.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", fiber.StatusNotFound, resp.StatusCode)
		}
	})
}

func TestRegisterNotFoundRoute_ResponseFormat(t *testing.T) {
	t.Run("response format is correct", func(t *testing.T) {
		app := fiber.New()
		RegisterNotFoundRoute(app)

		req, err := http.NewRequest("GET", "/test", nil)
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

		if _, ok := response["error"]; !ok {
			t.Error("Response missing 'error' key")
		}

		if _, ok := response["msg"]; !ok {
			t.Error("Response missing 'msg' key")
		}

		if len(response) != 2 {
			t.Errorf("Expected 2 keys in response, got %d", len(response))
		}
	})
}
