package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	Catalog    string
	Schema     string
	HTTPScheme string

	PathToPEM   string
	FileNamePEM string

	QueryTimeout time.Duration
}

func (cfg *Config) buildHTTPClient() (*http.Client, string, error) {
	if cfg.PathToPEM == "" && cfg.FileNamePEM == "" {
		return nil, "", nil
	}

	if cfg.PathToPEM == "" || cfg.FileNamePEM == "" {
		return nil, "", fmt.Errorf("both path_to_pem and file_name_pem must be provided together")
	}

	pemPath := filepath.Join(cfg.PathToPEM, cfg.FileNamePEM)

	certBytes, err := os.ReadFile(pemPath)
	if err != nil {
		return nil, "", fmt.Errorf("read pem file %q: %w", pemPath, err)
	}

	rootCAs := x509.NewCertPool()

	if ok := rootCAs.AppendCertsFromPEM(certBytes); !ok {
		return nil, "", fmt.Errorf("failed to append certificates from pem file %q", pemPath)
	}

	tlsConfig := &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	customClientName := "trino-provider-custom-tls"

	return httpClient, customClientName, nil
}

// buildDSN constructs the Data Source Name (DSN) for connecting to Trino based on the provided configuration and custom client name.
// It returns the DSN string or an error if the DSN cannot be constructed.
func (cfg *Config) buildDSN(customClientName string) (string, error) {

	u := &url.URL{
		Scheme: cfg.HTTPScheme,
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		User:   url.User(cfg.User),
	}

	if cfg.Password != "" {
		u.User = url.UserPassword(cfg.User, cfg.Password)
	}

	q := url.Values{}

	if cfg.Catalog != "" {
		q.Set("catalog", cfg.Catalog)
	}

	if cfg.Schema != "" {
		q.Set("schema", cfg.Schema)
	}

	if customClientName != "" {
		q.Set("custom_client", customClientName)
	}

	q.Set("source", "terraform-provider-trino")

	u.RawQuery = q.Encode()

	return u.String(), nil
}
