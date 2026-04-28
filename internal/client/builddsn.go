package client

import (
	"fmt"
	"net/url"
)

// buildDSN constructs the Data Source Name (DSN) for connecting to Trino based on the provided configuration and custom client name.
// It returns the DSN string or an error if the DSN cannot be constructed.
func buildDSN(cfg Config, customClientName string) (string, error) {

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
