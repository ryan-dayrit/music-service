package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	configContent := []byte(`grpc:
  network: tcp
  host: localhost
  port: "8080"

postgres:
  driver_name: postgres
  user: testuser
  db_name: testdb
  ssl_mode: disable
  password: testpass
  host: localhost
`)

	err := os.WriteFile(configPath, configContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory to tempDir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to restore original working directory: %v", err)
		}
	})
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Grpc.Network != "tcp" {
		t.Errorf("Expected network 'tcp', got '%s'", cfg.Grpc.Network)
	}
	if cfg.Grpc.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", cfg.Grpc.Host)
	}
	if cfg.Grpc.Port != "8080" {
		t.Errorf("Expected port '8080', got '%s'", cfg.Grpc.Port)
	}

	if cfg.Postgres.DriverName != "postgres" {
		t.Errorf("Expected driver_name 'postgres', got '%s'", cfg.Postgres.DriverName)
	}
	if cfg.Postgres.User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", cfg.Postgres.User)
	}
	if cfg.Postgres.DBName != "testdb" {
		t.Errorf("Expected db_name 'testdb', got '%s'", cfg.Postgres.DBName)
	}
	if cfg.Postgres.SSLMode != "disable" {
		t.Errorf("Expected ssl_mode 'disable', got '%s'", cfg.Postgres.SSLMode)
	}
	if cfg.Postgres.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got '%s'", cfg.Postgres.Password)
	}
	if cfg.Postgres.Host != "localhost" {
		t.Errorf("Expected database host 'localhost', got '%s'", cfg.Postgres.Host)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	tempDir := t.TempDir()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory to tempDir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to restore original working directory: %v", err)
		}
	})

	_, err = Load()
	if err == nil {
		t.Error("Expected error when config file doesn't exist, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	invalidYAML := []byte(`grpc:
  network: tcp
  host: localhost
  port: "8080"
  invalid yaml content: [[[
`)

	err := os.WriteFile(configPath, invalidYAML, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory to tempDir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to restore original working directory: %v", err)
		}
	})

	_, err = Load()
	if err == nil {
		t.Error("Expected error when parsing invalid YAML, got nil")
	}
}

func TestLoad_EmptyConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory to tempDir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to restore original working directory: %v", err)
		}
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed with empty config: %v", err)
	}

	if cfg == nil {
		t.Error("Expected non-nil config, got nil")
	}
}
