// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dialects

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/duanchi/min/v2/db/xorm/core"
	"github.com/duanchi/min/v2/db/xorm/schemas"
)

var (
	tdengineReservedWords = map[string]bool{
		"ADD":               true,
		"ALL":               true,
		"ALTER":             true,
		"ANALYZE":           true,
		"AND":               true,
		"AS":                true,
		"ASC":               true,
		"ASENSITIVE":        true,
		"BEFORE":            true,
		"BETWEEN":           true,
		"BIGINT":            true,
		"BINARY":            true,
		"BLOB":              true,
		"BOTH":              true,
		"BY":                true,
		"CALL":              true,
		"CASCADE":           true,
		"CASE":              true,
		"CHANGE":            true,
		"CHAR":              true,
		"CHARACTER":         true,
		"CHECK":             true,
		"COLLATE":           true,
		"COLUMN":            true,
		"CONDITION":         true,
		"CONNECTION":        true,
		"CONSTRAINT":        true,
		"CONTINUE":          true,
		"CONVERT":           true,
		"CREATE":            true,
		"CROSS":             true,
		"CURRENT_DATE":      true,
		"CURRENT_TIME":      true,
		"CURRENT_TIMESTAMP": true,
		"CURRENT_USER":      true,
		"CURSOR":            true,
		"DATABASE":          true,
		"DATABASES":         true,
		"DAY_HOUR":          true,
		"DAY_MICROSECOND":   true,
		"DAY_MINUTE":        true,
		"DAY_SECOND":        true,
		"DEC":               true,
		"DECIMAL":           true,
		"DECLARE":           true,
		"DEFAULT":           true,
		"DELAYED":           true,
		"DELETE":            true,
		"DESC":              true,
		"DESCRIBE":          true,
		"DETERMINISTIC":     true,
		"DISTINCT":          true,
		"DISTINCTROW":       true,
		"DIV":               true,
		"DOUBLE":            true,
		"DROP":              true,
		"DUAL":              true,
		"EACH":              true,
		"ELSE":              true,
		"ELSEIF":            true,
		"ENCLOSED":          true,
		"ESCAPED":           true,
		"EXISTS":            true,
		"EXIT":              true,
		"EXPLAIN":           true,
		"FALSE":             true,
		"FETCH":             true,
		"FLOAT":             true,
		"FLOAT4":            true,
		"FLOAT8":            true,
		"FOR":               true,
		"FORCE":             true,
		"FOREIGN":           true,
		"FROM":              true,
		"FULLTEXT":          true,
		"GOTO":              true,
		"GRANT":             true,
		"GROUP":             true,
		"HAVING":            true,
		"HIGH_PRIORITY":     true,
		"HOUR_MICROSECOND":  true,
		"HOUR_MINUTE":       true,
		"HOUR_SECOND":       true,
		"IF":                true,
		"IGNORE":            true,
		"IN":                true, "INDEX": true,
		"INFILE": true, "INNER": true, "INOUT": true,
		"INSENSITIVE": true, "INSERT": true, "INT": true,
		"INT1": true, "INT2": true, "INT3": true,
		"INT4": true, "INT8": true, "INTEGER": true,
		"INTERVAL": true, "INTO": true, "IS": true,
		"ITERATE": true, "JOIN": true, "KEY": true,
		"KEYS": true, "KILL": true, "LABEL": true,
		"LEADING": true, "LEAVE": true, "LEFT": true,
		"LIKE": true, "LIMIT": true, "LINEAR": true,
		"LINES": true, "LOAD": true, "LOCALTIME": true,
		"LOCALTIMESTAMP": true, "LOCK": true, "LONG": true,
		"LONGBLOB": true, "LONGTEXT": true, "LOOP": true,
		"LOW_PRIORITY": true, "MATCH": true, "MEDIUMBLOB": true,
		"MEDIUMINT": true, "MEDIUMTEXT": true, "MIDDLEINT": true,
		"MINUTE_MICROSECOND": true, "MINUTE_SECOND": true, "MOD": true,
		"MODIFIES": true, "NATURAL": true, "NOT": true,
		"NO_WRITE_TO_BINLOG": true, "NULL": true, "NUMERIC": true,
		"ON	OPTIMIZE": true, "OPTION": true,
		"OPTIONALLY": true, "OR": true, "ORDER": true,
		"OUT": true, "OUTER": true, "OUTFILE": true,
		"PRECISION": true, "PRIMARY": true, "PROCEDURE": true,
		"PURGE": true, "RAID0": true, "RANGE": true,
		"READ": true, "READS": true, "REAL": true,
		"REFERENCES": true, "REGEXP": true, "RELEASE": true,
		"RENAME": true, "REPEAT": true, "REPLACE": true,
		"REQUIRE": true, "RESTRICT": true, "RETURN": true,
		"REVOKE": true, "RIGHT": true, "RLIKE": true,
		"SCHEMA": true, "SCHEMAS": true, "SECOND_MICROSECOND": true,
		"SELECT": true, "SENSITIVE": true, "SEPARATOR": true,
		"SET": true, "SHOW": true, "SMALLINT": true,
		"SPATIAL": true, "SPECIFIC": true, "SQL": true,
		"SQLEXCEPTION": true, "SQLSTATE": true, "SQLWARNING": true,
		"SQL_BIG_RESULT": true, "SQL_CALC_FOUND_ROWS": true, "SQL_SMALL_RESULT": true,
		"SSL": true, "STARTING": true, "STRAIGHT_JOIN": true,
		"TABLE": true, "TERMINATED": true, "THEN": true,
		"TINYBLOB": true, "TINYINT": true, "TINYTEXT": true,
		"TO": true, "TRAILING": true, "TRIGGER": true,
		"TRUE": true, "UNDO": true, "UNION": true,
		"UNIQUE": true, "UNLOCK": true, "UNSIGNED": true,
		"UPDATE": true, "USAGE": true, "USE": true,
		"USING": true, "UTC_DATE": true, "UTC_TIME": true,
		"UTC_TIMESTAMP": true, "VALUES": true, "VARBINARY": true,
		"VARCHAR":      true,
		"VARCHARACTER": true,
		"VARYING":      true,
		"WHEN":         true,
		"WHERE":        true,
		"WHILE":        true,
		"WITH":         true,
		"WRITE":        true,
		"X509":         true,
		"XOR":          true,
		"YEAR_MONTH":   true,
		"ZEROFILL":     true,
	}

	tdengineQuoter = schemas.Quoter{
		Prefix:     '\'',
		Suffix:     '\'',
		IsReserved: schemas.AlwaysReserve,
	}
)

type tdengine struct {
	Base
	net               string
	addr              string
	params            map[string]string
	loc               *time.Location
	timeout           time.Duration
	tls               *tls.Config
	allowAllFiles     bool
	allowOldPasswords bool
	clientFoundRows   bool
	rowFormat         string
}

func (db *tdengine) Init(uri *URI) error {
	db.quoter = influxQuoter
	return db.Base.Init(db, uri)
}

func (db *tdengine) SetParams(params map[string]string) {
	rowFormat, ok := params["rowFormat"]
	if ok {
		var t = strings.ToUpper(rowFormat)
		switch t {
		case "COMPACT":
			fallthrough
		case "REDUNDANT":
			fallthrough
		case "DYNAMIC":
			fallthrough
		case "COMPRESSED":
			db.rowFormat = t
			break
		default:
			break
		}
	}
}

func (db *tdengine) SQLType(c *schemas.Column) string {
	var res string
	switch t := c.SQLType.Name; t {
	case schemas.Bool:
		res = schemas.TinyInt
		c.Length = 1
	case schemas.Serial:
		c.IsAutoIncrement = true
		c.IsPrimaryKey = true
		c.Nullable = false
		res = schemas.Int
	case schemas.BigSerial:
		c.IsAutoIncrement = true
		c.IsPrimaryKey = true
		c.Nullable = false
		res = schemas.BigInt
	case schemas.Bytea:
		res = schemas.Blob
	case schemas.TimeStampz:
		res = schemas.Char
		c.Length = 64
	case schemas.Enum: // mysql enum
		res = schemas.Enum
		res += "("
		opts := ""
		for v := range c.EnumOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		res += strings.TrimLeft(opts, ",")
		res += ")"
	case schemas.Set: // mysql set
		res = schemas.Set
		res += "("
		opts := ""
		for v := range c.SetOptions {
			opts += fmt.Sprintf(",'%v'", v)
		}
		res += strings.TrimLeft(opts, ",")
		res += ")"
	case schemas.NVarchar:
		res = schemas.Varchar
	case schemas.Uuid:
		res = schemas.Varchar
		c.Length = 40
	case schemas.Json:
		res = schemas.Text
	default:
		res = t
	}

	hasLen1 := (c.Length > 0)
	hasLen2 := (c.Length2 > 0)

	if res == schemas.BigInt && !hasLen1 && !hasLen2 {
		c.Length = 20
		hasLen1 = true
	}

	if hasLen2 {
		res += "(" + strconv.Itoa(c.Length) + "," + strconv.Itoa(c.Length2) + ")"
	} else if hasLen1 {
		res += "(" + strconv.Itoa(c.Length) + ")"
	}
	return res
}

func (db *tdengine) IsReserved(name string) bool {
	_, ok := mysqlReservedWords[strings.ToUpper(name)]
	return ok
}

func (db *tdengine) AutoIncrStr() string {
	return "AUTO_INCREMENT"
}

func (db *tdengine) IndexCheckSQL(tableName, idxName string) (string, []interface{}) {
	args := []interface{}{db.uri.DBName, tableName, idxName}
	sql := "SELECT `INDEX_NAME` FROM `INFORMATION_SCHEMA`.`STATISTICS`"
	sql += " WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? AND `INDEX_NAME`=?"
	return sql, args
}

func (db *tdengine) IsTableExist(queryer core.Queryer, ctx context.Context, tableName string) (bool, error) {
	sql := "SELECT `TABLE_NAME` from `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=? and `TABLE_NAME`=?"
	return db.HasRecords(queryer, ctx, sql, db.uri.DBName, tableName)
}

func (db *tdengine) AddColumnSQL(tableName string, col *schemas.Column) string {
	quoter := db.dialect.Quoter()
	s, _ := ColumnString(db, col, true)
	sql := fmt.Sprintf("ALTER TABLE %v ADD %v", quoter.Quote(tableName), s)
	if len(col.Comment) > 0 {
		sql += " COMMENT '" + col.Comment + "'"
	}
	return sql
}

func (db *tdengine) GetColumns(queryer core.Queryer, ctx context.Context, tableName string) ([]string, map[string]*schemas.Column, error) {
	args := []interface{}{db.uri.DBName, tableName}
	alreadyQuoted := "(INSTR(VERSION(), 'maria') > 0 && " +
		"(SUBSTRING_INDEX(VERSION(), '.', 1) > 10 || " +
		"(SUBSTRING_INDEX(VERSION(), '.', 1) = 10 && " +
		"(SUBSTRING_INDEX(SUBSTRING(VERSION(), 4), '.', 1) > 2 || " +
		"(SUBSTRING_INDEX(SUBSTRING(VERSION(), 4), '.', 1) = 2 && " +
		"SUBSTRING_INDEX(SUBSTRING(VERSION(), 6), '-', 1) >= 7)))))"
	s := "SELECT `COLUMN_NAME`, `IS_NULLABLE`, `COLUMN_DEFAULT`, `COLUMN_TYPE`," +
		" `COLUMN_KEY`, `EXTRA`, `COLUMN_COMMENT`, " +
		alreadyQuoted + " AS NEEDS_QUOTE " +
		"FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?" +
		" ORDER BY `COLUMNS`.ORDINAL_POSITION"

	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols := make(map[string]*schemas.Column)
	colSeq := make([]string, 0)
	for rows.Next() {
		col := new(schemas.Column)
		col.Indexes = make(map[string]int)

		var columnName, isNullable, colType, colKey, extra, comment string
		var alreadyQuoted bool
		var colDefault *string
		err = rows.Scan(&columnName, &isNullable, &colDefault, &colType, &colKey, &extra, &comment, &alreadyQuoted)
		if err != nil {
			return nil, nil, err
		}
		col.Name = strings.Trim(columnName, "` ")
		col.Comment = comment
		if "YES" == isNullable {
			col.Nullable = true
		}

		if colDefault != nil && (!alreadyQuoted || *colDefault != "NULL") {
			col.Default = *colDefault
			col.DefaultIsEmpty = false
		} else {
			col.DefaultIsEmpty = true
		}

		cts := strings.Split(colType, "(")
		colName := cts[0]
		colType = strings.ToUpper(colName)
		var len1, len2 int
		if len(cts) == 2 {
			idx := strings.Index(cts[1], ")")
			if colType == schemas.Enum && cts[1][0] == '\'' { // enum
				options := strings.Split(cts[1][0:idx], ",")
				col.EnumOptions = make(map[string]int)
				for k, v := range options {
					v = strings.TrimSpace(v)
					v = strings.Trim(v, "'")
					col.EnumOptions[v] = k
				}
			} else if colType == schemas.Set && cts[1][0] == '\'' {
				options := strings.Split(cts[1][0:idx], ",")
				col.SetOptions = make(map[string]int)
				for k, v := range options {
					v = strings.TrimSpace(v)
					v = strings.Trim(v, "'")
					col.SetOptions[v] = k
				}
			} else {
				lens := strings.Split(cts[1][0:idx], ",")
				len1, err = strconv.Atoi(strings.TrimSpace(lens[0]))
				if err != nil {
					return nil, nil, err
				}
				if len(lens) == 2 {
					len2, err = strconv.Atoi(lens[1])
					if err != nil {
						return nil, nil, err
					}
				}
			}
		}
		if colType == "FLOAT UNSIGNED" {
			colType = "FLOAT"
		}
		if colType == "DOUBLE UNSIGNED" {
			colType = "DOUBLE"
		}
		col.Length = len1
		col.Length2 = len2
		if _, ok := schemas.SqlTypes[colType]; ok {
			col.SQLType = schemas.SQLType{Name: colType, DefaultLength: len1, DefaultLength2: len2}
		} else {
			return nil, nil, fmt.Errorf("Unknown colType %v", colType)
		}

		if colKey == "PRI" {
			col.IsPrimaryKey = true
		}
		if colKey == "UNI" {
			// col.is
		}

		if extra == "auto_increment" {
			col.IsAutoIncrement = true
		}

		if !col.DefaultIsEmpty {
			if !alreadyQuoted && col.SQLType.IsText() {
				col.Default = "'" + col.Default + "'"
			} else if col.SQLType.IsTime() && !alreadyQuoted && col.Default != "CURRENT_TIMESTAMP" {
				col.Default = "'" + col.Default + "'"
			}
		}
		cols[col.Name] = col
		colSeq = append(colSeq, col.Name)
	}
	return colSeq, cols, nil
}

func (db *tdengine) GetTables(queryer core.Queryer, ctx context.Context) ([]*schemas.Table, error) {
	args := []interface{}{db.uri.DBName}
	s := "SELECT `TABLE_NAME`, `ENGINE`, `AUTO_INCREMENT`, `TABLE_COMMENT` from " +
		"`INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=? AND (`ENGINE`='MyISAM' OR `ENGINE` = 'InnoDB' OR `ENGINE` = 'TokuDB')"

	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]*schemas.Table, 0)
	for rows.Next() {
		table := schemas.NewEmptyTable()
		var name, engine string
		var autoIncr, comment *string
		err = rows.Scan(&name, &engine, &autoIncr, &comment)
		if err != nil {
			return nil, err
		}

		table.Name = name
		if comment != nil {
			table.Comment = *comment
		}
		table.StoreEngine = engine
		tables = append(tables, table)
	}
	return tables, nil
}

func (db *tdengine) SetQuotePolicy(quotePolicy QuotePolicy) {
	switch quotePolicy {
	case QuotePolicyNone:
		var q = influxQuoter
		q.IsReserved = schemas.AlwaysNoReserve
		db.quoter = q
	case QuotePolicyReserved:
		var q = influxQuoter
		q.IsReserved = db.IsReserved
		db.quoter = q
	case QuotePolicyAlways:
		fallthrough
	default:
		db.quoter = influxQuoter
	}
}

func (db *tdengine) GetIndexes(queryer core.Queryer, ctx context.Context, tableName string) (map[string]*schemas.Index, error) {
	args := []interface{}{db.uri.DBName, tableName}
	s := "SELECT `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"

	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexes := make(map[string]*schemas.Index, 0)
	for rows.Next() {
		var indexType int
		var indexName, colName, nonUnique string
		err = rows.Scan(&indexName, &nonUnique, &colName)
		if err != nil {
			return nil, err
		}

		if indexName == "PRIMARY" {
			continue
		}

		if "YES" == nonUnique || nonUnique == "1" {
			indexType = schemas.IndexType
		} else {
			indexType = schemas.UniqueType
		}

		colName = strings.Trim(colName, "` ")
		var isRegular bool
		if strings.HasPrefix(indexName, "IDX_"+tableName) || strings.HasPrefix(indexName, "UQE_"+tableName) {
			indexName = indexName[5+len(tableName):]
			isRegular = true
		}

		var index *schemas.Index
		var ok bool
		if index, ok = indexes[indexName]; !ok {
			index = new(schemas.Index)
			index.IsRegular = isRegular
			index.Type = indexType
			index.Name = indexName
			indexes[indexName] = index
		}
		index.AddColumn(colName)
	}
	return indexes, nil
}

func (db *tdengine) CreateTableSQL(table *schemas.Table, tableName string) ([]string, bool) {
	var sql = "CREATE TABLE IF NOT EXISTS "
	if tableName == "" {
		tableName = table.Name
	}

	quoter := db.Quoter()

	sql += quoter.Quote(tableName)
	sql += " ("

	if len(table.ColumnsSeq()) > 0 {
		pkList := table.PrimaryKeys

		for _, colName := range table.ColumnsSeq() {
			col := table.GetColumn(colName)
			s, _ := ColumnString(db, col, col.IsPrimaryKey && len(pkList) == 1)
			sql += s
			sql = strings.TrimSpace(sql)
			if len(col.Comment) > 0 {
				sql += " COMMENT '" + col.Comment + "'"
			}
			sql += ", "
		}

		if len(pkList) > 1 {
			sql += "PRIMARY KEY ( "
			sql += quoter.Join(pkList, ",")
			sql += " ), "
		}

		sql = sql[:len(sql)-2]
	}
	sql += ")"

	if table.StoreEngine != "" {
		sql += " ENGINE=" + table.StoreEngine
	}

	var charset = table.Charset
	if len(charset) == 0 {
		charset = db.URI().Charset
	}
	if len(charset) != 0 {
		sql += " DEFAULT CHARSET " + charset
	}

	if db.rowFormat != "" {
		sql += " ROW_FORMAT=" + db.rowFormat
	}
	return []string{sql}, true
}

func (db *tdengine) Filters() []Filter {
	return []Filter{}
}

type tdengineDriver struct {
}

func (p *tdengineDriver) Parse(driverName, dataSourceName string) (*URI, error) {
	dsnUrl, _ := url.Parse(dataSourceName)
	timeout := int64(0)

	if dsnUrl.Query().Has("timeout") {
		timeout, _ = strconv.ParseInt(dsnUrl.Query().Get("timeout"), 10, 64)
	}

	uri := &URI{
		DBType:  schemas.TDENGINE,
		Proto:   "",
		Host:    "",
		Port:    "",
		DBName:  strings.Trim(dsnUrl.Path, "/"),
		Timeout: time.Duration(timeout) * time.Second,
		Schema:  "",
	}

	return uri, nil
}
