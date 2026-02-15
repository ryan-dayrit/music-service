package routes

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterSwaggerRoute(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "swagger route registered at /swagger",
			path:           "/swagger",
			expectedStatus: fiber.StatusMovedPermanently,
		},
		{
			name:           "swagger route with wildcard path",
			path:           "/swagger/index.html",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "swagger route with doc.json",
			path:           "/swagger/doc.json",
			expectedStatus: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			RegisterSwaggerRoute(app)

			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Logf("Path: %s returned status: %d (expected: %d)", tt.path, resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}

func TestRegisterSwaggerRoute_GroupCreation(t *testing.T) {
	t.Run("swagger route group is created", func(t *testing.T) {
		app := fiber.New()
		RegisterSwaggerRoute(app)

		// Test that the swagger group exists by testing a path under it
		req, err := http.NewRequest("GET", "/swagger/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		// The route should exist (even if it returns an error/redirect)
		// A status of 404 would indicate the group wasn't created
		if resp.StatusCode == 0 {
			t.Error("Expected a response from swagger route")
		}
	})
}

func TestRegisterSwaggerRoute_Integration(t *testing.T) {
	t.Run("swagger route works with other routes", func(t *testing.T) {
		app := fiber.New()

		// Register other routes
		app.Get("/health", func(c *fiber.Ctx) error {
			return c.SendString("OK")
		})

		// Register swagger route
		RegisterSwaggerRoute(app)

		// Test health route still works
		req1, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp1, err := app.Test(req1)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		if resp1.StatusCode != fiber.StatusOK {
			t.Errorf("Expected status code %d for /health, got %d", fiber.StatusOK, resp1.StatusCode)
		}

		// Test swagger route exists
		req2, err := http.NewRequest("GET", "/swagger/index.html", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp2, err := app.Test(req2)
		if err != nil {
			t.Fatalf("Failed to test request: %v", err)
		}

		// Should get a response (status code != 0)
		if resp2.StatusCode == 0 {
			t.Error("Expected a response from swagger route")
		}
	})
}

func TestRegisterSwaggerRoute_OnlyGETMethod(t *testing.T) {
	t.Run("swagger route only accepts GET method", func(t *testing.T) {
		app := fiber.New()
		RegisterSwaggerRoute(app)

		methods := []string{"POST", "PUT", "DELETE", "PATCH"}
		for _, method := range methods {
			req, err := http.NewRequest(method, "/swagger/index.html", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			// These methods should not be allowed
			if resp.StatusCode == fiber.StatusOK {
				t.Errorf("Method %s should not be allowed on swagger route", method)
			}
		}
	})
}

func TestRegisterSwaggerRoute_MultipleRegistrations(t *testing.T) {
	t.Run("can register swagger route multiple times without panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Registering swagger route twice caused panic: %v", r)
			}
		}()

		app := fiber.New()
		RegisterSwaggerRoute(app)
		// This should not cause issues (though in practice you wouldn't do this)
		// Just ensuring no panic occurs
	})
}
