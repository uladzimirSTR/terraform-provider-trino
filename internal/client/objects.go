package client

import "time"

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
