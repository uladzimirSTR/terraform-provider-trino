package client

import "fmt"

func validateConfig(cfg Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("host is required")
	}

	if cfg.Port == 0 {
		return fmt.Errorf("port is required")
	}

	if cfg.User == "" {
		return fmt.Errorf("user is required")
	}

	if cfg.HTTPScheme == "" {
		return fmt.Errorf("http scheme is required")
	}

	if cfg.HTTPScheme != "http" && cfg.HTTPScheme != "https" {
		return fmt.Errorf("http scheme must be http or https")
	}

	if cfg.Password != "" && cfg.HTTPScheme != "https" {
		return fmt.Errorf("password authentication requires https")
	}

	hasPath := cfg.PathToPEM != ""
	hasFile := cfg.FileNamePEM != ""

	if hasPath != hasFile {
		return fmt.Errorf("path_to_pem and file_name_pem must be provided together or both omitted")
	}

	return nil
}
