package postgres

import "testing"

func TestConfig_StructFields(t *testing.T) {
	cfg := Config{
		DriverName: "postgres",
		User:       "testuser",
		DBName:     "testdb",
		SSLMode:    "disable",
		Password:   "testpass",
		Host:       "localhost",
	}

	if cfg.DriverName != "postgres" {
		t.Errorf("Expected DriverName 'postgres', got '%s'", cfg.DriverName)
	}
	if cfg.User != "testuser" {
		t.Errorf("Expected User 'testuser', got '%s'", cfg.User)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("Expected DBName 'testdb', got '%s'", cfg.DBName)
	}
	if cfg.SSLMode != "disable" {
		t.Errorf("Expected SSLMode 'disable', got '%s'", cfg.SSLMode)
	}
	if cfg.Password != "testpass" {
		t.Errorf("Expected Password 'testpass', got '%s'", cfg.Password)
	}
	if cfg.Host != "localhost" {
		t.Errorf("Expected Host 'localhost', got '%s'", cfg.Host)
	}
}

func TestConfig_ZeroValues(t *testing.T) {
	var cfg Config

	if cfg.DriverName != "" {
		t.Errorf("Expected empty DriverName, got '%s'", cfg.DriverName)
	}
	if cfg.User != "" {
		t.Errorf("Expected empty User, got '%s'", cfg.User)
	}
	if cfg.DBName != "" {
		t.Errorf("Expected empty DBName, got '%s'", cfg.DBName)
	}
	if cfg.SSLMode != "" {
		t.Errorf("Expected empty SSLMode, got '%s'", cfg.SSLMode)
	}
	if cfg.Password != "" {
		t.Errorf("Expected empty Password, got '%s'", cfg.Password)
	}
	if cfg.Host != "" {
		t.Errorf("Expected empty Host, got '%s'", cfg.Host)
	}
}

func TestConfig_DifferentSSLModes(t *testing.T) {
	tests := []struct {
		name    string
		sslMode string
	}{
		{"Disable SSL", "disable"},
		{"Require SSL", "require"},
		{"Verify CA", "verify-ca"},
		{"Verify Full", "verify-full"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				DriverName: "postgres",
				User:       "user",
				DBName:     "db",
				SSLMode:    tt.sslMode,
				Password:   "pass",
				Host:       "host",
			}

			if cfg.SSLMode != tt.sslMode {
				t.Errorf("Expected SSLMode '%s', got '%s'", tt.sslMode, cfg.SSLMode)
			}
		})
	}
}

func TestConfig_ValidPostgresConfig(t *testing.T) {
	cfg := Config{
		DriverName: "postgres",
		User:       "admin",
		DBName:     "production",
		SSLMode:    "require",
		Password:   "securepass123",
		Host:       "db.example.com",
	}

	if cfg.DriverName != "postgres" {
		t.Error("Expected postgres driver")
	}

	if cfg.User == "" || cfg.DBName == "" || cfg.Host == "" {
		t.Error("Essential config fields should not be empty")
	}
}
