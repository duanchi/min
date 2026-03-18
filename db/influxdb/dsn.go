package library

import (
	"net/url"
	"strings"
)

type Config struct {
	Host     string
	Ssl      bool
	Database string
	Token    string
	Version  string
}

func ParseDSN(dsn string) (*Config, error) {
	dsnUrl, _ := url.Parse(dsn)
	ssl := false
	if dsnUrl.Query().Get("ssl") == "true" {
		ssl = true
	}

	return &Config{
		Host:     dsnUrl.Host,
		Ssl:      ssl,
		Database: strings.Trim(dsnUrl.Path, "/"),
		Token:    dsnUrl.Query().Get("token"),
		Version:  "v3",
	}, nil
}
