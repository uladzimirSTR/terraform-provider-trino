package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func buildHTTPClient(cfg Config) (*http.Client, string, error) {
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
