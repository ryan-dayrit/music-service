package grpc

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfig_YAMLUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		want     Config
		wantErr  bool
	}{
		{
			name: "valid config with all fields",
			yamlData: `
network: tcp
host: localhost
port: "50051"
`,
			want: Config{
				Network: "tcp",
				Host:    "localhost",
				Port:    "50051",
			},
			wantErr: false,
		},
		{
			name: "valid config with empty fields",
			yamlData: `
network: ""
host: ""
port: ""
`,
			want: Config{
				Network: "",
				Host:    "",
				Port:    "",
			},
			wantErr: false,
		},
		{
			name: "valid config with only network",
			yamlData: `
network: tcp
`,
			want: Config{
				Network: "tcp",
				Host:    "",
				Port:    "",
			},
			wantErr: false,
		},
		{
			name: "valid config with numeric port",
			yamlData: `
network: tcp
host: 0.0.0.0
port: "9090"
`,
			want: Config{
				Network: "tcp",
				Host:    "0.0.0.0",
				Port:    "9090",
			},
			wantErr: false,
		},
		{
			name:     "empty yaml",
			yamlData: ``,
			want: Config{
				Network: "",
				Host:    "",
				Port:    "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Config
			err := yaml.Unmarshal([]byte(tt.yamlData), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("yaml.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("yaml.Unmarshal() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestConfig_YAMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "marshal complete config",
			config: Config{
				Network: "tcp",
				Host:    "localhost",
				Port:    "50051",
			},
			wantErr: false,
		},
		{
			name: "marshal empty config",
			config: Config{
				Network: "",
				Host:    "",
				Port:    "",
			},
			wantErr: false,
		},
		{
			name: "marshal partial config",
			config: Config{
				Network: "tcp",
				Host:    "",
				Port:    "8080",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("yaml.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Unmarshal back to verify
			var got Config
			err = yaml.Unmarshal(data, &got)
			if err != nil {
				t.Errorf("yaml.Unmarshal() after Marshal error = %v", err)
				return
			}

			if got != tt.config {
				t.Errorf("Marshal/Unmarshal roundtrip failed: got = %+v, want %+v", got, tt.config)
			}
		})
	}
}

func TestConfig_StructFields(t *testing.T) {
	config := Config{
		Network: "tcp",
		Host:    "0.0.0.0",
		Port:    "9000",
	}

	if config.Network != "tcp" {
		t.Errorf("Config.Network = %q, want %q", config.Network, "tcp")
	}
	if config.Host != "0.0.0.0" {
		t.Errorf("Config.Host = %q, want %q", config.Host, "0.0.0.0")
	}
	if config.Port != "9000" {
		t.Errorf("Config.Port = %q, want %q", config.Port, "9000")
	}
}

func TestConfig_ZeroValue(t *testing.T) {
	var config Config

	if config.Network != "" {
		t.Errorf("Zero value Config.Network = %q, want empty string", config.Network)
	}
	if config.Host != "" {
		t.Errorf("Zero value Config.Host = %q, want empty string", config.Host)
	}
	if config.Port != "" {
		t.Errorf("Zero value Config.Port = %q, want empty string", config.Port)
	}
}
