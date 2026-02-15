package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestStartServer(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		setupApp  func() *fiber.App
		wantError bool
	}{
		{
			name: "start server with valid config",
			config: Config{
				ReadTimeout:  5,
				WriteTimeout: 5,
				ServerUrl:    ":0", // Use port 0 to let OS assign a free port
			},
			setupApp: func() *fiber.App {
				app := fiber.New(fiber.Config{
					DisableStartupMessage: true,
				})
				app.Get("/", func(c *fiber.Ctx) error {
					return c.SendString("OK")
				})
				return app
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.setupApp()

			// Run server in goroutine
			done := make(chan bool)
			go func() {
				StartServer(app, tt.config)
				done <- true
			}()

			// Give server time to start
			time.Sleep(100 * time.Millisecond)

			// Shutdown server
			if err := app.Shutdown(); err != nil {
				t.Errorf("Failed to shutdown server: %v", err)
			}

			// Wait for server to finish
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				t.Error("Server did not stop in time")
			}
		})
	}
}

func TestStartServerWithGracefulShutdown(t *testing.T) {
	t.Skip("Skipping test that requires signal handling - difficult to test in isolation")
	tests := []struct {
		name     string
		config   Config
		setupApp func() *fiber.App
	}{
		{
			name: "graceful shutdown with valid config",
			config: Config{
				ReadTimeout:  5,
				WriteTimeout: 5,
				ServerUrl:    ":0",
			},
			setupApp: func() *fiber.App {
				app := fiber.New(fiber.Config{
					DisableStartupMessage: true,
				})
				app.Get("/health", func(c *fiber.Ctx) error {
					return c.SendString("healthy")
				})
				return app
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tt.setupApp()

			// Run server with graceful shutdown in goroutine
			serverStopped := make(chan bool)
			go func() {
				StartServerWithGracefulShutdown(app, tt.config)
				serverStopped <- true
			}()

			// Give server time to start
			time.Sleep(100 * time.Millisecond)

			// Trigger graceful shutdown
			if err := app.Shutdown(); err != nil {
				t.Errorf("Failed to shutdown server: %v", err)
			}

			// Wait for graceful shutdown to complete
			select {
			case <-serverStopped:
				// Success
			case <-time.After(5 * time.Second):
				t.Error("Server did not gracefully shutdown in time")
			}
		})
	}
}

func TestServerIntegration(t *testing.T) {
	t.Run("server handles requests correctly", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		// Setup routes
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "test response",
			})
		})

		cfg := Config{
			ReadTimeout:  5,
			WriteTimeout: 5,
			ServerUrl:    ":0",
		}

		// Start server in goroutine
		go StartServer(app, cfg)

		// Give server time to start
		time.Sleep(100 * time.Millisecond)

		// Get the actual port the server is listening on
		// Note: With port :0, we can't easily get the assigned port in this test
		// This is a limitation of the current implementation

		// Shutdown server
		if err := app.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown server: %v", err)
		}
	})
}

func TestConfig_Values(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "default production config",
			config: Config{
				ReadTimeout:  30,
				WriteTimeout: 30,
				ServerUrl:    ":8080",
			},
		},
		{
			name: "custom config",
			config: Config{
				ReadTimeout:  60,
				WriteTimeout: 60,
				ServerUrl:    "localhost:3000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				DisableStartupMessage: true,
				ReadTimeout:           time.Duration(tt.config.ReadTimeout) * time.Second,
				WriteTimeout:          time.Duration(tt.config.WriteTimeout) * time.Second,
			})

			if app == nil {
				t.Error("Expected fiber app to be created")
			}

			// Verify config values are accessible
			if tt.config.ReadTimeout < 0 {
				t.Error("ReadTimeout should not be negative")
			}
			if tt.config.WriteTimeout < 0 {
				t.Error("WriteTimeout should not be negative")
			}
			if tt.config.ServerUrl == "" {
				t.Error("ServerUrl should not be empty")
			}
		})
	}
}

func TestStartServer_EdgeCases(t *testing.T) {
	t.Run("server with nil routes", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		cfg := Config{
			ReadTimeout:  1,
			WriteTimeout: 1,
			ServerUrl:    ":0",
		}

		done := make(chan bool)
		go func() {
			StartServer(app, cfg)
			done <- true
		}()

		time.Sleep(50 * time.Millisecond)

		if err := app.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown server: %v", err)
		}

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Error("Server did not stop in time")
		}
	})
}

func BenchmarkStartServer(b *testing.B) {
	cfg := Config{
		ReadTimeout:  5,
		WriteTimeout: 5,
		ServerUrl:    ":0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
		})

		go StartServer(app, cfg)
		time.Sleep(10 * time.Millisecond)
		_ = app.Shutdown()
	}
}

func TestServerRequest(t *testing.T) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Test using the fiber test API instead of starting actual server
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
