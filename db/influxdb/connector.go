package library

import "database/sql/driver"

type connector struct {
	cfg               *Config // immutable private copy.
	encodedAttributes string  // Encoded connection attributes.
}

func (c *connector) Connect() (driver.Conn, error) {
	return nil, nil
}

func (c *connector) Driver() driver.Driver {
	return &InfluxDBDriver{}
}
