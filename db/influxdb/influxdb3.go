package library

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

const (
	FORMAT_JSON    = "json"
	FORMAT_JSONL   = "jsonl"
	FORMAT_CSV     = "csv"
	FORMAT_PERTTY  = "pretty"
	FORMAT_PARQUET = "parquet"
)

var driverName = "influxdb"

func init() {
	if driverName != "" {
		sql.Register(driverName, &InfluxDBDriver{})
	}
}

type QuerySQLRequestV3 struct {
	Q      string        `json:"q"`
	Db     string        `json:"db"`
	Params []interface{} `json:"params"`
	Format string        `json:"format"`
}

type WriteErrorResponse struct {
	Error string `json:"error"`
	Data  struct {
		OriginalLine string `json:"original_line"`
		LineNumber   int    `json:"line_number"`
		ErrorMessage string `json:"error_message"`
	} `json:"data"`
}

type InfluxDBDriver struct{}

func (d InfluxDBDriver) Open(dsn string) (driver.Conn, error) {

	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	c := connector{
		cfg: cfg,
	}
	return c.Connect(context.Background())
}

func (d InfluxDBDriver) OpenConnector(dsn string) (driver.Connector, error) {
	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return connector{
		cfg: cfg,
	}, nil
}
