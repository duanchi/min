package library

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/duanchi/min/v2/requests/http"
	"github.com/duanchi/min/v2/util"
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
	c := newConnector(cfg)
	return c.Connect(context.Background())
}

func (d InfluxDBDriver) OpenConnector(dsn string) (driver.Connector, error) {
	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return newConnector(cfg), nil
}

type Client struct {
	token    string
	schema   string
	url      string
	format   string
	version  string
	database string
}

func New(host string, token string, database ...string) *Client {
	c := &Client{
		token:   token,
		url:     host,
		schema:  "http",
		format:  FORMAT_JSON,
		version: "v3",
	}

	if len(database) > 0 {
		c.database = database[0]
	}
	return c
}

func (c *Client) SetSchema(schema string) *Client {
	c.schema = schema
	return c
}

func (c *Client) Database(database string) *Session {
	c.database = database
	return c.NewSession().Database(database)
}

func (c *Client) Table(table string) *Session {
	return c.NewSession().Table(table)
}

func (c *Client) SetFormat(format string) *Client {
	c.format = format
	return c
}

func (c *Client) SetVersion(version string) *Client {
	c.version = version
	return c
}

func (c *Client) Query(sql string, beans any) error {
	return c.NewSession().Query(sql, beans)
}

func (c *Client) Insert(bean any, timestamp ...int64) error {
	return c.NewSession().Insert(bean, timestamp...)
}

func (c *Client) NewSession() *Session {
	baseUrl := c.schema + "://" + c.url
	request := http.New()
	request.Method(http.METHOD_POST).BearerToken(c.token).Url(c.url)
	switch c.version {
	case "v3":
		request.BaseUrl(baseUrl + "/api/v3/")
	}
	return &Session{
		client:   c,
		database: c.database,
		request:  &request,
		tags:     map[string]string{},
	}
}

type Session struct {
	client    *Client
	request   *http.Request
	method    string
	database  string
	sql       string
	params    []interface{}
	table     string
	tags      map[string]string
	timestamp int64
}

func (s *Session) Query(query string, rows any) (err error) {
	response, err := s.request.Url("/query_sql").JSON(QuerySQLRequestV3{
		Q:      query,
		Db:     s.client.database,
		Params: s.params,
		Format: s.client.format,
	}).Response()

	if err != nil {
		return err
	}

	err = response.BindJSON(rows)
	return
}

func (s *Session) Database(database string) *Session {
	s.database = database
	return s
}

func (s *Session) Table(table string) *Session {
	s.table = table
	return s
}

func (s *Session) Tag(name, value string) *Session {
	s.tags[name] = value
	return s
}

func (s *Session) Insert(bean any, timestamp ...int64) (err error) {
	ts := int64(0)
	if len(timestamp) == 0 {
		ts = time.Now().UnixMicro()
	} else {
		ts = timestamp[0]
	}
	if s.table == "" {
		k := reflect.TypeOf(bean)
		s.table = util.ToSnake(k.Name())
	}
	lineData := buildLineData(s.table, s.tags, bean, ts)
	response, err := s.request.Url("/write_lp?db="+s.database+"&no_sync=true&precision=auto").Body([]byte(lineData)).Header("Content-Type", "text/plain; charset=utf-8").Response()
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		errStruct := WriteErrorResponse{}
		e := response.BindJSON(&errStruct)
		if e == nil {
			return errors.New(errStruct.Error + ", at line " + strconv.Itoa(errStruct.Data.LineNumber) + ", [" + errStruct.Data.OriginalLine + "], " + errStruct.Data.ErrorMessage)
		}
		err = errors.New(string(response.Payload))
	}
	return err
}

func (s *Session) BatchInsert(beans any, timestamp ...int64) (err error) {

	beanValue := reflect.ValueOf(beans)

	if beanValue.Len() == 0 {
		return nil
	}

	ts := int64(0)
	lineDataList := []string{}
	if len(timestamp) == 0 {
		ts = time.Now().UnixMicro()
	} else {
		ts = timestamp[0]
	}
	if s.table == "" {
		k := beanValue.Index(0).Type()
		s.table = util.ToSnake(k.Name())
	}
	for i := range beanValue.Len() {
		lineDataList = append(lineDataList, buildLineData(s.table, s.tags, beanValue.Index(i).Interface(), ts))
	}

	response, err := s.request.Url("/write_lp?db="+s.database+"&no_sync=true&precision=auto").Body([]byte(strings.Join(lineDataList, "\n"))).Header("Content-Type", "text/plain; charset=utf-8").Response()
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		errStruct := WriteErrorResponse{}
		e := response.BindJSON(&errStruct)
		if e == nil {
			return errors.New(errStruct.Error + ", at line " + strconv.Itoa(errStruct.Data.LineNumber) + ", [" + errStruct.Data.OriginalLine + "], " + errStruct.Data.ErrorMessage)
		}
		err = e
	}
	return err
}

func buildLineData(database string, tags map[string]string, bean any, timestamp int64) string {
	line := database
	if len(tags) > 0 {
		for tagKey, tag := range tags {
			line += "," + escapeString(tagKey, []string{",", "=", " "}) + "=" + escapeString(tag, []string{",", "=", " "})
		}
	}
	line += " " + structToDataLine(bean) + " " + strconv.FormatInt(timestamp, 10)
	return line
}

func structToDataLine(obj any) (line string) {
	var slice []string
	k := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Struct:
		{
			for f := range k.Fields() {
				tag := f.Tag.Get("xorm")
				if !strings.Contains(tag, "notnull") && !v.FieldByName(f.Name).IsZero() {
					slice = append(slice, util.ToSnake(f.Name)+"="+parseValue(v.FieldByName(f.Name).Interface()))
				}
			}
		}
	case reflect.Map:
		{
			for _, key := range v.MapKeys() {
				slice = append(slice, key.String()+"="+parseValue(v.MapIndex(key).Interface()))
			}
		}
	case reflect.String:
		{
			slice = append(slice, v.String())
		}

	}
	return strings.Join(slice, ",")
}

func parseValue(value interface{}) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return "\"" + escapeString(value.(string), []string{"\\", "\""}) + "\""
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
