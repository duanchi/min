package library

import (
	"database/sql/driver"
	"encoding/json"
	"io"
	"time"
)

type influxField struct {
	name       string
	columnType string
}

type resultSet struct {
	columns     []influxField
	columnNames []string
	done        bool
}

type influxRows struct {
	mc     *InfluxConn
	rs     resultSet
	finish func()
}

func (rows *influxRows) Columns() []string {
	if rows.rs.columnNames != nil {
		return rows.rs.columnNames
	}

	columns := make([]string, len(rows.rs.columns))
	for i := range columns {
		columns[i] = rows.rs.columns[i].name
	}

	rows.rs.columnNames = columns
	return columns
}

func (rows *influxRows) Close() error {
	rows.mc = nil
	rows.rs.columnNames = nil
	rows.rs.columns = nil
	return nil
}

func (rows *influxRows) Next(dest []driver.Value) error {
	if rows.mc.resultLen > 0 && rows.mc.resultIndex+1 < rows.mc.resultLen {
		if len(rows.rs.columns) > 0 {
			o := map[string]driver.Value{}
			json.Unmarshal(rows.mc.result[rows.mc.resultIndex], &o)

			for i := range dest {
				// key := rows.rs.columns[i].name
				if _, has := o[rows.rs.columns[i].name]; has {
					switch rows.rs.columns[i].name {
					case "time":
						{
							t, _ := time.Parse("2006-01-02T15:04:05.000000", o[rows.rs.columns[i].name].(string))
							o[rows.rs.columns[i].name] = t.Local().UnixNano()
						}
					}

					dest[i] = o[rows.rs.columns[i].name]
				}
			}
		}

		rows.mc.resultIndex++
		return nil
	}

	// dest = append(dest, rows.rs.columns[0])
	return io.EOF
}

func (rows *influxRows) HasNextResultSet() bool {
	if rows.mc == nil || rows.mc.resultLen == 0 || (rows.mc.resultIndex+1 >= rows.mc.resultLen) {
		return false
	}
	return true
}

func (rows *influxRows) NextResultSet() (err error) {
	rows.rs.columns, err = rows.mc.readColumns()
	if err == nil {
		rows.mc.setResult = true
	}
	return nil
}
