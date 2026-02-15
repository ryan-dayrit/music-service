package rest

import (
	"testing"
)

func TestConfig_Struct(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		wantReadTO    int
		wantWriteTO   int
		wantServerURL string
	}{
		{
			name: "valid config with all fields",
			config: Config{
				ReadTimeout:  30,
				WriteTimeout: 30,
				ServerUrl:    ":8080",
			},
			wantReadTO:    30,
			wantWriteTO:   30,
			wantServerURL: ":8080",
		},
		{
			name: "config with zero values",
			config: Config{
				ReadTimeout:  0,
				WriteTimeout: 0,
				ServerUrl:    "",
			},
			wantReadTO:    0,
			wantWriteTO:   0,
			wantServerURL: "",
		},
		{
			name: "config with large timeout values",
			config: Config{
				ReadTimeout:  3600,
				WriteTimeout: 3600,
				ServerUrl:    "localhost:9090",
			},
			wantReadTO:    3600,
			wantWriteTO:   3600,
			wantServerURL: "localhost:9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.ReadTimeout != tt.wantReadTO {
				t.Errorf("ReadTimeout = %v, want %v", tt.config.ReadTimeout, tt.wantReadTO)
			}
			if tt.config.WriteTimeout != tt.wantWriteTO {
				t.Errorf("WriteTimeout = %v, want %v", tt.config.WriteTimeout, tt.wantWriteTO)
			}
			if tt.config.ServerUrl != tt.wantServerURL {
				t.Errorf("ServerUrl = %v, want %v", tt.config.ServerUrl, tt.wantServerURL)
			}
		})
	}
}
