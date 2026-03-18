package library

import (
	"context"
	"database/sql/driver"
)

type connector struct {
	cfg               *Config // immutable private copy.
	encodedAttributes string  // Encoded connection attributes.
}

func (c connector) Connect(ctx context.Context) (driver.Conn, error) {
	return &InfluxConn{
		ctx: ctx,
		cfg: c.cfg,
	}, nil
}

func (c connector) Driver() driver.Driver {
	return &InfluxDBDriver{}
}
