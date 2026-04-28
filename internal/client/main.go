package client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	trinodriver "github.com/trinodb/trino-go-client/trino"

	_ "github.com/trinodb/trino-go-client/trino"
)

type Client struct {
	db *sql.DB
}

func NewClient(cfg Config) (*Client, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	httpClient, customClientName, err := cfg.buildHTTPClient()
	if err != nil {
		return nil, err
	}

	if httpClient != nil {
		if err := trinodriver.RegisterCustomClient(customClientName, httpClient); err != nil {
			return nil, fmt.Errorf("register trino custom http client: %w", err)
		}
	}

	dsn, err := cfg.buildDSN(customClientName)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("trino", dsn)
	if err != nil {
		return nil, fmt.Errorf("open trino connection: %w", err)
	}

	if cfg.QueryTimeout == 0 {
		cfg.QueryTimeout = 60 * time.Second
	}

	return &Client{db: db}, nil
}

func (c *Client) Exec(ctx context.Context, query string) error {
	if query == "" {
		return fmt.Errorf("query is empty")
	}

	_, err := c.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("execute trino query: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	if c.db == nil {
		return nil
	}

	return c.db.Close()
}
