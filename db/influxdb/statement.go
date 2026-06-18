package library

import (
	"database/sql/driver"
)

type influxStatement struct {
	connector  *InfluxConn
	paramCount int
}

func (stmt *influxStatement) Close() error {
	return nil
}

func (stmt *influxStatement) NumInput() int {
	return 0
}

func (stmt *influxStatement) Exec(args []driver.Value) (driver.Result, error) {
	return nil, nil
}

func (stmt *influxStatement) Query(args []driver.Value) (driver.Rows, error) {
	rows := new(influxRows)
	rows.mc = stmt.connector
	return rows, nil
}
