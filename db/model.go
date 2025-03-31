package db

import (
	"context"
	"database/sql"
	"github.com/duanchi/min/v2/db/xorm"
	"github.com/duanchi/min/v2/db/xorm/caches"
	"github.com/duanchi/min/v2/db/xorm/contexts"
	"github.com/duanchi/min/v2/db/xorm/core"
	"github.com/duanchi/min/v2/db/xorm/dialects"
	"github.com/duanchi/min/v2/db/xorm/log"
	"github.com/duanchi/min/v2/db/xorm/names"
	"github.com/duanchi/min/v2/db/xorm/schemas"
	_interface "github.com/duanchi/min/v2/interface"
	"io"
	"reflect"
	"strings"
	"time"
)

func Model(model interface{}) *ModelMapper {
	instance := ModelMapper{
		Mapper: model,
	}

	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()
	modelKind := modelType.Kind()

	if modelKind == reflect.Slice {
		modelType = modelType.Elem()
	}

	instance.Struct = reflect.New(modelType).Interface().(_interface.ModelInterface)
	instance.Options()

	return &instance
}

type ModelMapper struct {
	Mapper  interface{}
	Struct  _interface.ModelInterface
	options map[string]interface{}
	engine  *xorm.Engine
}

// func (this *ModelMapper) Options (options map[string]interface{}) {
func (this *ModelMapper) Options() {
	source := this.Struct.Source()
	table := this.Struct.Table()

	if source == "" {
		source = "default"
	}
	this.SetSource(source)

	if table != "" {
		this.options["table"] = table
	}
	/*this.options = options
	if v, ok := options["source"]; ok {
		this.SetSource(v.(string))
	} else {
		this.SetSource("default")
	}

	if v, ok := options["table"]; ok {
		this.engine.Table(v)
	}*/
}

func (this *ModelMapper) NewMapper() interface{} {
	modelValue := reflect.ValueOf(this.Mapper).Elem()
	modelType := modelValue.Type()
	modelKind := modelType.Kind()

	if modelKind == reflect.Slice {
		modelType = modelType.Elem()
	}

	return reflect.New(modelType).Interface()
}

func (this *ModelMapper) Init() *ModelMapper {
	this.Options()
	return this
}

func (this *ModelMapper) GetEngine() *xorm.Engine {
	return this.engine
}

func (this *ModelMapper) SetEngine(db *xorm.Engine) {
	this.engine = db
}

func (this *ModelMapper) SetSource(name string) {
	this.engine = Engine(name)
}

// EnableSessionID if enable session id
func (this *ModelMapper) EnableSessionID(enable bool) {
	this.engine.EnableSessionID(enable)
}

// SetCacher sets cacher for the table
func (this *ModelMapper) SetCacher(tableName string, cacher caches.Cacher) {
	this.engine.SetCacher(tableName, cacher)
}

// GetCacher returns the cachher of the special table
func (this *ModelMapper) GetCacher(tableName string) caches.Cacher {
	return this.engine.GetCacher(tableName)
}

// SetQuotePolicy sets the special quote policy
func (this *ModelMapper) SetQuotePolicy(quotePolicy dialects.QuotePolicy) {
	this.engine.SetQuotePolicy(quotePolicy)
}

// BufferSize sets buffer size for iterate
func (this *ModelMapper) BufferSize(size int) *xorm.Session {
	return this.engine.BufferSize(size)
}

// ShowSQL show SQL statement or not on logger if log level is great than INFO
func (this *ModelMapper) ShowSQL(show ...bool) {
	this.engine.ShowSQL(show...)
}

// Logger return the logger interface
func (this *ModelMapper) Logger() log.ContextLogger {
	return this.engine.Logger()
}

// SetLogger set the new logger
func (this *ModelMapper) SetLogger(logger interface{}) {
	this.engine.SetLogger(logger)
}

// SetLogLevel sets the logger level
func (this *ModelMapper) SetLogLevel(level log.LogLevel) {
	this.engine.SetLogLevel(level)
}

// SetDisableGlobalCache disable global cache or not
func (this *ModelMapper) SetDisableGlobalCache(disable bool) {
	this.engine.SetDisableGlobalCache(disable)
}

// DriverName return the current sql driver's name
func (this *ModelMapper) DriverName() string {
	return this.engine.DriverName()
}

// DataSourceName return the current connection string
func (this *ModelMapper) DataSourceName() string {
	return this.engine.DataSourceName()
}

// SetMapper set the name mapping rules
func (this *ModelMapper) SetMapper(mapper names.Mapper) {
	this.engine.SetMapper(mapper)
}

// SetTableMapper set the table name mapping rule
func (this *ModelMapper) SetTableMapper(mapper names.Mapper) {
	this.engine.SetTableMapper(mapper)
}

// SetColumnMapper set the column name mapping rule
func (this *ModelMapper) SetColumnMapper(mapper names.Mapper) {
	this.engine.SetColumnMapper(mapper)
}

// Quote Use QuoteStr quote the string sql
func (this *ModelMapper) Quote(value string) string {
	return this.engine.Quote(value)
}

// QuoteTo quotes string and writes into the buffer
func (this *ModelMapper) QuoteTo(buf *strings.Builder, value string) {
	this.engine.QuoteTo(buf, value)
}

// SQLType A simple wrapper to dialect's core.SqlType method
func (this *ModelMapper) SQLType(c *schemas.Column) string {
	return this.engine.SQLType(c)
}

// AutoIncrStr Database's autoincrement statement
func (this *ModelMapper) AutoIncrStr() string {
	return this.engine.AutoIncrStr()
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (this *ModelMapper) SetConnMaxLifetime(d time.Duration) {
	this.engine.SetConnMaxLifetime(d)
}

// SetMaxOpenConns is only available for go 1.2+
func (this *ModelMapper) SetMaxOpenConns(conns int) {
	this.engine.SetMaxOpenConns(conns)
}

// SetMaxIdleConns set the max idle connections on pool, default is 2
func (this *ModelMapper) SetMaxIdleConns(conns int) {
	this.engine.SetMaxIdleConns(conns)
}

// SetDefaultCacher set the default cacher. Xorm's default not enable cacher.
func (this *ModelMapper) SetDefaultCacher(cacher caches.Cacher) {
	this.engine.SetDefaultCacher(cacher)
}

// GetDefaultCacher returns the default cacher
func (this *ModelMapper) GetDefaultCacher() caches.Cacher {
	return this.engine.GetDefaultCacher()
}

// NoCache If you has set default cacher, and you want temporilly stop use cache,
// you can use NoCache()
func (this *ModelMapper) NoCache() *xorm.Session {
	return this.engine.NoCache()
}

// NoCascade If you do not want to auto cascade load object
func (this *ModelMapper) NoCascade() *xorm.Session {
	return this.engine.NoCascade()
}

// MapCacher Set a table use a special cacher
func (this *ModelMapper) MapCacher(bean interface{}, cacher caches.Cacher) error {
	return this.engine.MapCacher(bean, cacher)
}

// NewDB provides an interface to operate database directly
func (this *ModelMapper) NewDB() (*core.DB, error) {
	return this.engine.NewDB()
}

// DB return the wrapper of sql.DB
func (this *ModelMapper) DB() *core.DB {
	return this.engine.DB()
}

// Dialect return database dialect
func (this *ModelMapper) Dialect() dialects.Dialect {
	return this.engine.Dialect()
}

// NewSession New a session
func (this *ModelMapper) NewSession() *xorm.Session {
	if table, ok := this.options["table"]; ok {
		return this.engine.NewSession().Table(table)
	} else {
		return this.engine.NewSession()
	}
}

// Close the engine
func (this *ModelMapper) Close() error {
	return this.NewSession().Close()
}

// Ping tests if database is alive
func (this *ModelMapper) Ping() error {
	return this.NewSession().Ping()
}

// SQL method let's you manually write raw SQL and operate
// For example:
//
//	this.engine.SQL("select * from user").Find(&users)
//
// This    code will execute "select * from user" and set the records to users
func (this *ModelMapper) SQL(query interface{}, args ...interface{}) *xorm.Session {
	return this.NewSession().SQL(query, args...)
}

// NoAutoTime Default if your struct has "created" or "updated" filed tag, the fields
// will automatically be filled with current time when Insert or Update
// invoked. Call NoAutoTime if you dont' want to fill automatically.
func (this *ModelMapper) NoAutoTime() *xorm.Session {
	return this.engine.NoAutoTime()
}

// NoAutoCondition disable auto generate Where condition from bean or not
func (this *ModelMapper) NoAutoCondition(no ...bool) *xorm.Session {
	return this.engine.NoAutoCondition(no...)
}

// DBMetas Retrieve all tables, columns, indexes' informations from database.
func (this *ModelMapper) DBMetas() ([]*schemas.Table, error) {
	return this.engine.DBMetas()
}

// DumpAllToFile dump database all table structs and data to a file
func (this *ModelMapper) DumpAllToFile(fp string, tp ...schemas.DBType) error {
	return this.engine.DumpAllToFile(fp, tp...)
}

// DumpAll dump database all table structs and data to w
func (this *ModelMapper) DumpAll(w io.Writer, tp ...schemas.DBType) error {
	return this.engine.DumpAll(w, tp...)
}

// DumpTablesToFile dump specified tables to SQL file.
func (this *ModelMapper) DumpTablesToFile(tables []*schemas.Table, fp string, tp ...schemas.DBType) error {
	return this.engine.DumpTablesToFile(tables, fp, tp...)
}

// DumpTables dump specify tables to io.Writer
func (this *ModelMapper) DumpTables(tables []*schemas.Table, w io.Writer, tp ...schemas.DBType) error {
	return this.engine.DumpTables(tables, w, tp...)
}

// Cascade use cascade or not
func (this *ModelMapper) Cascade(trueOrFalse ...bool) *xorm.Session {
	return this.NewSession().Cascade(trueOrFalse...)
}

// Where method provide a condition query
func (this *ModelMapper) Where(query interface{}, args ...interface{}) *xorm.Session {
	return this.NewSession().Where(query, args...)
}

// ID method provoide a condition as (id) = ?
func (this *ModelMapper) ID(id interface{}) *xorm.Session {
	return this.NewSession().ID(id)
}

func (this *ModelMapper) Id(id interface{}) *xorm.Session {
	return this.NewSession().ID(id)
}

// Before apply before Processor, affected bean is passed to closure arg
func (this *ModelMapper) Before(closures func(interface{})) *xorm.Session {
	return this.NewSession().Before(closures)
}

// After apply after insert Processor, affected bean is passed to closure arg
func (this *ModelMapper) After(closures func(interface{})) *xorm.Session {
	return this.NewSession().After(closures)
}

// Charset set charset when create table, only support mysql now
func (this *ModelMapper) Charset(charset string) *xorm.Session {
	return this.NewSession().Charset(charset)
}

// StoreEngine set store engine when create table, only support mysql now
func (this *ModelMapper) StoreEngine(storeEngine string) *xorm.Session {
	return this.NewSession().StoreEngine(storeEngine)
}

// Distinct use for distinct columns. Caution: when you are using cache,
// distinct will not be cached because cache system need id,
// but distinct will not provide id
func (this *ModelMapper) Distinct(columns ...string) *xorm.Session {
	return this.NewSession().Distinct(columns...)
}

// Select customerize your select columns or contents
func (this *ModelMapper) Select(str string) *xorm.Session {
	return this.NewSession().Select(str)
}

// Cols only use the parameters as select or update columns
func (this *ModelMapper) Cols(columns ...string) *xorm.Session {
	return this.NewSession().Cols(columns...)
}

// AllCols indicates that all columns should be use
func (this *ModelMapper) AllCols() *xorm.Session {
	return this.NewSession().AllCols()
}

// MustCols specify some columns must use even if they are empty
func (this *ModelMapper) MustCols(columns ...string) *xorm.Session {
	return this.NewSession().MustCols(columns...)
}

// UseBool xorm automatically retrieve condition according struct, but
// if struct has bool field, it will ignore them. So use UseBool
// to tell system to do not ignore them.
// If no parameters, it will use all the bool field of struct, or
// it will use parameters's columns
func (this *ModelMapper) UseBool(columns ...string) *xorm.Session {
	return this.NewSession().UseBool(columns...)
}

// Omit only not use the parameters as select or update columns
func (this *ModelMapper) Omit(columns ...string) *xorm.Session {
	return this.NewSession().Omit(columns...)
}

// Nullable set null when column is zero-value and nullable for update
func (this *ModelMapper) Nullable(columns ...string) *xorm.Session {
	return this.NewSession().Nullable(columns...)
}

// In will generate "column IN (?, ?)"
func (this *ModelMapper) In(column string, args ...interface{}) *xorm.Session {
	return this.NewSession().In(column, args...)
}

// NotIn will generate "column NOT IN (?, ?)"
func (this *ModelMapper) NotIn(column string, args ...interface{}) *xorm.Session {
	return this.NewSession().NotIn(column, args...)
}

// Incr provides a update string like "column = column + ?"
func (this *ModelMapper) Incr(column string, args ...interface{}) *xorm.Session {
	return this.NewSession().Incr(column, args...)
}

// Decr provides a update string like "column = column - ?"
func (this *ModelMapper) Decr(column string, args ...interface{}) *xorm.Session {
	return this.NewSession().Decr(column, args...)
}

// SetExpr provides a update string like "column = {expression}"
func (this *ModelMapper) SetExpr(column string, expression interface{}) *xorm.Session {
	return this.NewSession().SetExpr(column, expression)
}

// Table temporarily change the Get, Find, Update's table
func (this *ModelMapper) Table(tableNameOrBean interface{}) *xorm.Session {
	return this.engine.Table(tableNameOrBean)
}

// Alias set the table alias
func (this *ModelMapper) Alias(alias string) *xorm.Session {
	return this.NewSession().Alias(alias)
}

// Limit will generate "LIMIT start, limit"
func (this *ModelMapper) Limit(limit int, start ...int) *xorm.Session {
	return this.NewSession().Limit(limit, start...)
}

// Desc will generate "ORDER BY column1 DESC, column2 DESC"
func (this *ModelMapper) Desc(colNames ...string) *xorm.Session {
	return this.NewSession().Desc(colNames...)
}

// Asc will generate "ORDER BY column1,column2 Asc"
// This method can chainable use.
//
//	this.engine.Desc("name").Asc("age").Find(&users)
//	// SELECT * FROM user ORDER BY name DESC, age ASC
func (this *ModelMapper) Asc(colNames ...string) *xorm.Session {
	return this.NewSession().Asc(colNames...)
}

// OrderBy will generate "ORDER BY order"
func (this *ModelMapper) OrderBy(order string) *xorm.Session {
	return this.NewSession().OrderBy(order)
}

// Prepare enables prepare statement
func (this *ModelMapper) Prepare() *xorm.Session {
	return this.NewSession().Prepare()
}

// Join the join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (this *ModelMapper) Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *xorm.Session {
	return this.NewSession().Join(joinOperator, tablename, condition, args...)
}

// GroupBy generate group by statement
func (this *ModelMapper) GroupBy(keys string) *xorm.Session {
	return this.NewSession().GroupBy(keys)
}

// Having generate having statement
func (this *ModelMapper) Having(conditions string) *xorm.Session {
	return this.NewSession().Having(conditions)
}

// TableInfo get table info according to bean's content
func (this *ModelMapper) TableInfo(bean interface{}) (*schemas.Table, error) {
	return this.engine.TableInfo(bean)
}

// IsTableEmpty if a table has any reocrd
func (this *ModelMapper) IsTableEmpty(bean interface{}) (bool, error) {
	return this.NewSession().IsTableEmpty(bean)
}

// IsTableExist if a table is exist
func (this *ModelMapper) IsTableExist(beanOrTableName interface{}) (bool, error) {
	return this.NewSession().IsTableExist(beanOrTableName)
}

// IDOf get id from one struct
/*func (this *ModelMapper) IDOf(bean interface{}) (schemas.PK, error) {
	return this.engine.IDOf(bean)
}*/

// TableName returns table name with schema prefix if has
func (this *ModelMapper) TableName(bean interface{}, includeSchema ...bool) string {
	return this.engine.TableName(bean, includeSchema...)
}

// IDOfV get id from one value of struct
/*func (this *ModelMapper) IDOfV(rv reflect.Value) (schemas.PK, error) {
	return this.engine.NewSession().
}*/

// CreateIndexes create indexes
func (this *ModelMapper) CreateIndexes(bean interface{}) error {
	return this.NewSession().CreateIndexes(bean)
}

// CreateUniques create uniques
func (this *ModelMapper) CreateUniques(bean interface{}) error {
	return this.NewSession().CreateUniques(bean)
}

// ClearCacheBean if enabled cache, clear the cache bean
func (this *ModelMapper) ClearCacheBean(bean interface{}, id string) error {
	return this.engine.ClearCacheBean(bean, id)
}

// ClearCache if enabled cache, clear some tables' cache
func (this *ModelMapper) ClearCache(beans ...interface{}) error {
	return this.engine.ClearCache(beans...)
}

// UnMapType remove table from tables cache
func (this *ModelMapper) UnMapType(t reflect.Type) {
	this.engine.UnMapType(t)
}

// Sync the new struct changes to database, this method will automatically add
// table, column, index, unique. but will not delete or change anything.
// If you change some field, you should change the database manually.
func (this *ModelMapper) Sync(beans ...interface{}) error {
	return this.engine.Sync(beans...)
}

// Sync2 synchronize structs to database tables
func (this *ModelMapper) Sync2(beans ...interface{}) error {
	return this.engine.Sync2(beans...)
}

// CreateTables create tabls according bean
func (this *ModelMapper) CreateTables(beans ...interface{}) error {
	return this.engine.CreateTables(beans...)
}

// DropTables drop specify tables
func (this *ModelMapper) DropTables(beans ...interface{}) error {
	return this.engine.DropTables(beans...)
}

// DropIndexes drop indexes of a table
func (this *ModelMapper) DropIndexes(bean interface{}) error {
	return this.NewSession().DropIndexes(bean)
}

// Exec raw sql
func (this *ModelMapper) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	return this.NewSession().Exec(sqlOrArgs...)
}

// Query a raw sql and return records as []map[string][]byte
func (this *ModelMapper) QueryBytes(sqlOrArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	return this.NewSession().QueryBytes(sqlOrArgs...)
}

// Query a raw sql and return records as []map[string]Value
func (this *ModelMapper) QueryValue(sqlOrArgs ...interface{}) (resultsSlice []map[string]xorm.Value, err error) {
	return this.NewSession().QueryValue(sqlOrArgs...)
}

// Query a raw sql and return records as Result
func (this *ModelMapper) QueryResult(sqlOrArgs ...interface{}) (result *xorm.ResultValue) {
	return this.NewSession().QueryResult(sqlOrArgs...)
}

// QueryString runs a raw sql and return records as []map[string]string
func (this *ModelMapper) QueryString(sqlOrArgs ...interface{}) ([]map[string]string, error) {
	return this.NewSession().QueryString(sqlOrArgs...)
}

// QueryInterface runs a raw sql and return records as []map[string]interface{}
func (this *ModelMapper) QueryInterface(sqlOrArgs ...interface{}) ([]map[string]interface{}, error) {
	return this.NewSession().QueryInterface(sqlOrArgs...)
}

// Insert one or more records
func (this *ModelMapper) Insert(beans ...interface{}) (int64, error) {
	return this.NewSession().Insert(this.Mapper.([]interface{})...)
}

// InsertOne insert only one record
func (this *ModelMapper) InsertOne() (int64, error) {
	return this.NewSession().InsertOne(this.Mapper)
}

// Update records, bean's non-empty fields are updated contents,
// condiBean' non-empty filds are conditions
// CAUTION:
//
//	1.bool will defaultly be updated content nor conditions
//	 You should call UseBool if you have bool to use.
//	2.float32 & float64 may be not inexact as conditions
func (this *ModelMapper) Update(condiBeans ...interface{}) (int64, error) {
	return this.NewSession().Update(this.Mapper, condiBeans...)
}

// Delete records, bean's non-empty fields are conditions
func (this *ModelMapper) Delete() (int64, error) {
	return this.NewSession().Delete(this.Mapper)
}

// Get retrieve one record from table, bean's non-empty fields
// are conditions
func (this *ModelMapper) Get() (bool, error) {
	return this.NewSession().Get(this.Mapper)
}

// Exist returns true if the record exist otherwise return false
func (this *ModelMapper) Exist() (bool, error) {
	return this.NewSession().Exist(this.Mapper)
}

// Find retrieve records from table, condiBeans's non-empty fields
// are conditions. beans could be []Struct, []*Struct, map[int64]Struct
// map[int64]*Struct
func (this *ModelMapper) Find(condiBeans ...interface{}) error {
	return this.NewSession().Find(this.Mapper, condiBeans...)
}

// FindAndCount find the results and also return the counts
func (this *ModelMapper) FindAndCount(condiBean ...interface{}) (int64, error) {
	return this.NewSession().FindAndCount(this.Mapper, condiBean...)
}

// Iterate record by record handle records from table, bean's non-empty fields
// are conditions.
func (this *ModelMapper) Iterate(fun xorm.IterFunc) error {
	return this.NewSession().Iterate(this.Mapper, fun)
}

// Rows return sql.Rows compatible Rows obj, as a forward Iterator object for iterating record by record, bean's non-empty fields
// are conditions.
func (this *ModelMapper) Rows() (*xorm.Rows, error) {
	return this.NewSession().Rows(this.Mapper)
}

// Count counts the records. bean's non-empty fields are conditions.
func (this *ModelMapper) Count() (int64, error) {
	return this.NewSession().Count(this.Mapper)
}

// Sum sum the records by some column. bean's non-empty fields are conditions.
func (this *ModelMapper) Sum(colName string) (float64, error) {
	return this.NewSession().Sum(this.Mapper, colName)
}

// SumInt sum the records by some column. bean's non-empty fields are conditions.
func (this *ModelMapper) SumInt(colName string) (int64, error) {
	return this.NewSession().SumInt(this.Mapper, colName)
}

// Sums sum the records by some columns. bean's non-empty fields are conditions.
func (this *ModelMapper) Sums(colNames ...string) ([]float64, error) {
	return this.NewSession().Sums(this.Mapper, colNames...)
}

// SumsInt like Sums but return slice of int64 instead of float64.
func (this *ModelMapper) SumsInt(colNames ...string) ([]int64, error) {
	return this.NewSession().SumsInt(this.Mapper, colNames...)
}

// ImportFile SQL DDL file
func (this *ModelMapper) ImportFile(ddlPath string) ([]sql.Result, error) {
	return this.NewSession().ImportFile(ddlPath)
}

// Import SQL DDL from io.Reader
func (this *ModelMapper) Import(r io.Reader) ([]sql.Result, error) {
	return this.NewSession().Import(r)
}

// GetColumnMapper returns the column name mapper
func (this *ModelMapper) GetColumnMapper() names.Mapper {
	return this.engine.GetColumnMapper()
}

// GetTableMapper returns the table name mapper
func (this *ModelMapper) GetTableMapper() names.Mapper {
	return this.engine.GetTableMapper()
}

// GetTZLocation returns time zone of the application
func (this *ModelMapper) GetTZLocation() *time.Location {
	return this.engine.GetTZLocation()
}

// SetTZLocation sets time zone of the application
func (this *ModelMapper) SetTZLocation(tz *time.Location) {
	this.engine.SetTZLocation(tz)
}

// GetTZDatabase returns time zone of the database
func (this *ModelMapper) GetTZDatabase() *time.Location {
	return this.engine.GetTZDatabase()
}

// SetTZDatabase sets time zone of the database
func (this *ModelMapper) SetTZDatabase(tz *time.Location) {
	this.engine.SetTZDatabase(tz)
}

// SetSchema sets the schema of database
func (this *ModelMapper) SetSchema(schema string) {
	this.engine.SetSchema(schema)
}

func (this *ModelMapper) AddHook(hook contexts.Hook) {
	this.engine.AddHook(hook)
}

// Unscoped always disable struct tag "deleted"
func (this *ModelMapper) Unscoped() *xorm.Session {
	return this.engine.Unscoped()
}

// ContextHook creates a session with the context
func (this *ModelMapper) Context(ctx context.Context) *xorm.Session {
	return this.engine.Context(ctx)
}

// SetDefaultContext set the default context
func (this *ModelMapper) SetDefaultContext(ctx context.Context) {
	this.engine.SetDefaultContext(ctx)
}

// PingContext tests if database is alive
func (this *ModelMapper) PingContext(ctx context.Context) error {
	return this.engine.PingContext(ctx)
}

// Transaction Execute sql wrapped in a transaction(abbr as tx), tx will automatic commit if no errors occurred
func (this *ModelMapper) Transaction(f func(*xorm.Session) (interface{}, error)) (interface{}, error) {
	return this.engine.Transaction(f)
}
