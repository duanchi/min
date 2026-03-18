package library

import (
	"bufio"
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"maps"
	"strings"

	"github.com/duanchi/min/v2/requests/http"
)

type session struct {
	request   *http.Request
	method    string
	database  string
	sql       string
	params    []interface{}
	table     string
	tags      map[string]string
	timestamp int64
}

type InfluxConn struct {
	ctx         context.Context
	cfg         *Config
	result      [][]byte
	resultLen   int64
	resultIndex int64
	setResult   bool
}

func (mc *InfluxConn) Prepare(query string) (driver.Stmt, error) {
	return &influxStatement{
		connector: mc,
	}, nil
}

func (mc *InfluxConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func (mc *InfluxConn) Close() (err error) {
	mc.clear()
	return
}

func (mc *InfluxConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	sess := request(mc)
	databaseName := mc.cfg.Database
	// 拆分dbname
	if strings.HasPrefix(query, "USE ") {
		queryStack := strings.SplitN(query, ";", 2)
		if len(queryStack) == 2 {
			query = queryStack[1]
			databaseName = strings.Trim(queryStack[0][3:], " '")
		}
	}

	response, err := sess.request.Url("/query_sql").JSON(QuerySQLRequestV3{
		Q:      query,
		Db:     databaseName,
		Format: FORMAT_JSONL,
	}).Response()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(response.Payload))
	for scanner.Scan() {
		mc.result = append(mc.result, scanner.Bytes())
	}
	mc.resultLen = int64(len(mc.result))
	rows := new(influxRows)
	rows.mc = mc
	rows.rs.columns, err = mc.readColumns()
	rows.mc.setResult = true
	return rows, nil
}

func (mc *InfluxConn) readColumns() ([]influxField, error) {
	columns := []influxField{}
	if mc.resultLen > 0 {
		line := map[string]any{}
		err := json.Unmarshal(mc.result[0], &line)
		if err != nil {
			return nil, err
		}
		for name := range maps.Keys(line) {
			columns = append(columns, influxField{
				name: name,
			})
		}
	}
	return columns, nil
}

func (mc *InfluxConn) clear() {
	mc.result = nil
	mc.setResult = false
	mc.resultIndex = 0
	mc.resultLen = 0
}

func request(c *InfluxConn) *session {
	schema := "http"
	if c.cfg.Ssl {
		schema = "https"
	}

	baseUrl := schema + "://" + c.cfg.Host
	req := http.New()
	req.Method(http.METHOD_POST).BearerToken(c.cfg.Token)
	switch c.cfg.Version {
	case "v3":
		req.BaseUrl(baseUrl + "/api/v3/")
	}
	return &session{
		request: &req,
		tags:    map[string]string{},
	}
}
