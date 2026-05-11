package client

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validConfig() Config {
	return Config{
		Host:       "localhost",
		Port:       8080,
		User:       "test_user",
		HTTPScheme: "http",
		Catalog:    "memory",
		Schema:     "default",
	}
}

func TestValidateConfigValid(t *testing.T) {
	cfg := validConfig()

	if err := validateConfig(cfg); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestValidateConfigRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{
			name: "missing host",
			mutate: func(cfg *Config) {
				cfg.Host = ""
			},
			wantErr: "host is required",
		},
		{
			name: "missing port",
			mutate: func(cfg *Config) {
				cfg.Port = 0
			},
			wantErr: "port is required",
		},
		{
			name: "missing user",
			mutate: func(cfg *Config) {
				cfg.User = ""
			},
			wantErr: "user is required",
		},
		{
			name: "missing scheme",
			mutate: func(cfg *Config) {
				cfg.HTTPScheme = ""
			},
			wantErr: "http scheme is required",
		},
		{
			name: "invalid scheme",
			mutate: func(cfg *Config) {
				cfg.HTTPScheme = "ftp"
			},
			wantErr: "http scheme must be http or https",
		},
		{
			name: "password with http",
			mutate: func(cfg *Config) {
				cfg.Password = "secret"
			},
			wantErr: "password authentication requires https",
		},
		{
			name: "pem path without file",
			mutate: func(cfg *Config) {
				cfg.PathToPEM = "/tmp"
			},
			wantErr: "path_to_pem and file_name_pem must be provided together or both omitted",
		},
		{
			name: "pem file without path",
			mutate: func(cfg *Config) {
				cfg.FileNamePEM = "ca.pem"
			},
			wantErr: "path_to_pem and file_name_pem must be provided together or both omitted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			tt.mutate(&cfg)

			err := validateConfig(cfg)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestBuildDSN(t *testing.T) {
	cfg := Config{
		Host:       "trino.example.com",
		Port:       8443,
		User:       "alice",
		Password:   "secret",
		HTTPScheme: "https",
		Catalog:    "iceberg",
		Schema:     "analytics",
	}

	dsn, err := cfg.buildDSN("custom-client")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("invalid dsn: %v", err)
	}

	if u.Scheme != "https" {
		t.Fatalf("expected scheme https, got %q", u.Scheme)
	}

	if u.Host != "trino.example.com:8443" {
		t.Fatalf("expected host trino.example.com:8443, got %q", u.Host)
	}

	username := u.User.Username()
	if username != "alice" {
		t.Fatalf("expected username alice, got %q", username)
	}

	password, ok := u.User.Password()
	if !ok || password != "secret" {
		t.Fatalf("expected password secret, got %q", password)
	}

	q := u.Query()

	if q.Get("catalog") != "iceberg" {
		t.Fatalf("expected catalog iceberg, got %q", q.Get("catalog"))
	}

	if q.Get("schema") != "analytics" {
		t.Fatalf("expected schema analytics, got %q", q.Get("schema"))
	}

	if q.Get("custom_client") != "custom-client" {
		t.Fatalf("expected custom_client custom-client, got %q", q.Get("custom_client"))
	}

	if q.Get("source") != "terraform-provider-trino" {
		t.Fatalf("expected source terraform-provider-trino, got %q", q.Get("source"))
	}
}

func TestBuildHTTPClientWithoutPEM(t *testing.T) {
	cfg := validConfig()

	httpClient, customClientName, err := cfg.buildHTTPClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpClient != nil {
		t.Fatalf("expected nil http client")
	}

	if customClientName != "" {
		t.Fatalf("expected empty custom client name, got %q", customClientName)
	}
}

func TestBuildHTTPClientMissingPEMPair(t *testing.T) {
	cfg := validConfig()
	cfg.PathToPEM = "/tmp"

	_, _, err := cfg.buildHTTPClient()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "both path_to_pem and file_name_pem must be provided together") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildHTTPClientInvalidPEM(t *testing.T) {
	dir := t.TempDir()
	pemPath := filepath.Join(dir, "ca.pem")

	if err := os.WriteFile(pemPath, []byte("not a pem file"), 0600); err != nil {
		t.Fatalf("write pem file: %v", err)
	}

	cfg := validConfig()
	cfg.PathToPEM = dir
	cfg.FileNamePEM = "ca.pem"

	_, _, err := cfg.buildHTTPClient()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to append certificates") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecEmptyQuery(t *testing.T) {
	c := &Client{}

	err := c.Exec(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "query is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCloseNilDB(t *testing.T) {
	c := &Client{}

	if err := c.Close(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestNewClientInvalidConfig(t *testing.T) {
	cfg := validConfig()
	cfg.Host = ""

	client, err := NewClient(cfg)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if client != nil {
		t.Fatalf("expected nil client")
	}
}
