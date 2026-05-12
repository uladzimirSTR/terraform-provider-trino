terraform {
  required_providers {
    trino = {
      source  = "ulstr/trino"
      version = "0.1.0"
    }
  }
}

provider "trino" {
  host        = "localhost"
  port        = 8080
  user        = "admin"
  password    = "admin"
  http_scheme = "https"

  # path_to_pem   = "/path/to/certs"
  # file_name_pem = "ca.pem"
}
