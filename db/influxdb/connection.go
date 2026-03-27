package library

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func (mc *InfluxConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	// ts := time.Now().UnixNano()
	// lineDataList := []string{}

	sess := request(mc)
	databaseName := mc.cfg.Database
	tableName := ""
	columns := []string{}
	// 拆分dbname
	if strings.HasPrefix(query, "USE ") {
		queryStack := strings.SplitN(query, ";", 2)
		if len(queryStack) == 2 {
			databaseName = strings.Trim(queryStack[0][3:], " '")
			query = queryStack[1]
		}
	}
	// 拆分table

	if splitQuery := strings.SplitN(query, "VALUES", 2); strings.HasPrefix(query, "INSERT") && len(splitQuery) == 2 {
		splitTable := strings.SplitN(splitQuery[0][12:], " ", 2)

		if len(splitTable) == 2 {
			ts := time.Now().UnixNano()
			tableName = strings.Trim(splitTable[0], " ")
			columns = strings.Split(strings.Trim(splitTable[1], " ()"), ",")
			columnSize := len(columns)
			groups := len(args) / columnSize
			lineDataList := []string{}
			ids := []int64{}
			if (len(args) % columnSize) > 0 {
				groups += 1
			}
			for i := range groups {
				line, tsid := buildLineData(tableName, columns, args[columnSize*i:columnSize*(i+1)], ts)
				lineDataList = append(lineDataList, line)
				ids = append(ids, tsid)
			}
			response, err := sess.request.Url("/write_lp?db="+databaseName+"&no_sync=true&precision=auto").Body([]byte(strings.Join(lineDataList, "\n"))).Header("Content-Type", "text/plain; charset=utf-8").Response()

			if err != nil {
				return nil, err
			}

			if response.StatusCode >= 400 {
				errStruct := WriteErrorResponse{}
				e := response.BindJSON(&errStruct)
				if e == nil {
					return nil, errors.New(errStruct.Error + ", at line " + strconv.Itoa(errStruct.Data.LineNumber) + ", [" + errStruct.Data.OriginalLine + "], " + errStruct.Data.ErrorMessage)
				}
				err = e
			}

			return &influxResult{
				affectedRows: int64(groups),
				insertIds:    ids,
			}, nil
		}
	}

	return nil, nil
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

	re := regexp.MustCompile("\\?")
	index := 0
	l := len(args)
	query = re.ReplaceAllStringFunc(query, func(match string) string {
		index++
		if index > l {
			return match
		}
		return parseValue(args[index-1], "'")
	})

	response, err := sess.request.Url("/query_sql").JSON(QuerySQLRequestV3{
		Q:      query,
		Db:     databaseName,
		Format: FORMAT_JSONL,
	}).Response()
	if err != nil {
		return nil, err
	}

	mc.result = bytes.Split(response.Payload, []byte("\n"))

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

func buildLineData(table string, columns []string, beans []driver.Value, timestamp int64) (string, int64) {
	line := table
	ts := int64(0)
	tagLine := []string{}
	dataLine := []string{}
	for i := range columns {
		if strings.HasPrefix(columns[i], "TAG:") {
			tagLine = append(tagLine, ","+columns[i][4:]+"="+escapeString(parseValue(beans[i]), []string{",", "=", " "}))
		} else if columns[i] == "time" {
			ts = parseTimestamp(beans[i])
		} else {
			dataLine = append(dataLine, columns[i]+"="+parseValue(beans[i], "\""))
		}
	}
	if ts == 0 {
		ts = timestamp
	}
	line += strings.Join(tagLine, "") + " " + strings.Join(dataLine, ",") + " " + strconv.FormatInt(ts, 10)
	return line, ts
}

func parseTimestamp(timeObject any) int64 {
	switch timeObject.(type) {
	case int64:
		return timeObject.(int64)
	case string:
		timing, _ := time.Parse("2006-01-02 15:04:05.999999999", timeObject.(string))
		return timing.UnixNano()
	case time.Time:
		return timeObject.(time.Time).UnixNano()
	}
	return 0
}

func parseValue(value interface{}, quote ...string) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		if len(quote) > 0 {
			return quote[0] + escapeString(value.(string), []string{"\\", "\""}) + quote[0]
		}

		return escapeString(value.(string), []string{"\\", "\""})

	case reflect.Int:
		return strconv.Itoa(value.(int)) + "i"
	case reflect.Int64:
		return strconv.FormatInt(value.(int64), 10) + "i"
	case reflect.Uint32:
		return strconv.Itoa(int(value.(uint32))) + "u"
	case reflect.Uint64:
		return strconv.Itoa(int(value.(uint64))) + "u"
	case reflect.Bool:
		if value.(bool) {
			return "true"
		}
		return "false"
	case reflect.Float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	default:
		return fmt.Sprintf("%s", value)
	}
}

func escapeString(s string, search []string) string {
	for _, element := range search {
		s = strings.ReplaceAll(s, element, "\\"+element)
	}
	return s
}
