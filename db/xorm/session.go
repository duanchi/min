// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/duanchi/min/v2/db/xorm/contexts"
	"github.com/duanchi/min/v2/db/xorm/convert"
	"github.com/duanchi/min/v2/db/xorm/core"
	"github.com/duanchi/min/v2/db/xorm/log"
	"github.com/duanchi/min/v2/db/xorm/private/json"
	"github.com/duanchi/min/v2/db/xorm/private/statements"
	"github.com/duanchi/min/v2/db/xorm/schemas"
)

type sessionType bool

const (
	engineSession sessionType = false
	groupSession  sessionType = true
)

// Session keep a pointer to sql.DB and provides all execution of all
// kind of database operations.
type Session struct {
	engine                 *Engine
	tx                     *core.Tx
	statement              *statements.Statement
	currentTransaction     *Transaction
	isAutoCommit           bool
	isCommitedOrRollbacked bool
	isSqlFunc              bool
	isAutoClose            bool
	isClosed               bool
	prepareStmt            bool
	// Automatically reset the statement after operations that execute a SQL
	// query such as Count(), Find(), Get(), ...
	autoResetStatement bool

	// !nashtsai! storing these beans due to yet committed tx
	afterInsertBeans map[interface{}]*[]func(interface{})
	afterUpdateBeans map[interface{}]*[]func(interface{})
	afterDeleteBeans map[interface{}]*[]func(interface{})
	// --

	beforeClosures  []func(interface{})
	afterClosures   []func(interface{})
	afterProcessors []executedProcessor

	stmtCache map[uint32]*core.Stmt //key: hash.Hash32 of (queryStr, len(queryStr))

	lastSQL     string
	lastSQLArgs []interface{}
	showSQL     bool

	rollbackSavePointID string

	ctx         context.Context
	sessionType sessionType

	err error
}

func newSessionID() string {
	hash := sha256.New()
	_, err := io.CopyN(hash, rand.Reader, 50)
	if err != nil {
		return "????????????????????"
	}
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr[0:20]
}

func newSession(engine *Engine) *Session {
	var ctx context.Context
	if engine.logSessionID {
		ctx = context.WithValue(engine.defaultContext, log.SessionIDKey, newSessionID())
	} else {
		ctx = engine.defaultContext
	}

	session := &Session{
		ctx:    ctx,
		engine: engine,
		tx:     nil,
		statement: statements.NewStatement(
			engine.dialect,
			engine.tagParser,
			engine.DatabaseTZ,
		),
		isClosed:               false,
		isAutoCommit:           true,
		isCommitedOrRollbacked: false,
		isAutoClose:            false,
		isSqlFunc:              false,
		autoResetStatement:     true,
		prepareStmt:            false,

		afterInsertBeans: make(map[interface{}]*[]func(interface{}), 0),
		afterUpdateBeans: make(map[interface{}]*[]func(interface{}), 0),
		afterDeleteBeans: make(map[interface{}]*[]func(interface{}), 0),
		beforeClosures:   make([]func(interface{}), 0),
		afterClosures:    make([]func(interface{}), 0),
		afterProcessors:  make([]executedProcessor, 0),
		stmtCache:        make(map[uint32]*core.Stmt),

		lastSQL:     "",
		lastSQLArgs: make([]interface{}, 0),

		sessionType: engineSession,
	}
	if engine.logSessionID {
		session.ctx = context.WithValue(session.ctx, log.SessionKey, session)
	}
	return session
}

// Close release the connection from pool
func (session *Session) Close() error {
	for _, v := range session.stmtCache {
		if err := v.Close(); err != nil {
			return err
		}
	}

	if !session.isClosed {
		// When Close be called, if session is a transaction and do not call
		// Commit or Rollback, then call Rollback.
		if session.tx != nil && !session.isCommitedOrRollbacked {
			if err := session.Rollback(); err != nil {
				return err
			}
		}
		session.tx = nil
		session.stmtCache = nil
		session.isClosed = true
	}
	return nil
}

func (session *Session) db() *core.DB {
	return session.engine.db
}

func (session *Session) Engine() *Engine {
	return session.engine
}

func (session *Session) getQueryer() core.Queryer {
	if session.tx != nil {
		return session.tx
	}
	return session.db()
}

// ContextCache enable context cache or not
func (session *Session) ContextCache(context contexts.ContextCache) *Session {
	session.statement.SetContextCache(context)
	return session
}

// IsClosed returns if session is closed
func (session *Session) IsClosed() bool {
	return session.isClosed
}

func (session *Session) resetStatement() {
	if session.autoResetStatement {
		session.statement.Reset()
	}
	session.isSqlFunc = false
}

// Prepare set a flag to session that should be prepare statement before execute query
func (session *Session) Prepare() *Session {
	session.prepareStmt = true
	return session
}

// Before Apply before Processor, affected bean is passed to closure arg
func (session *Session) Before(closures func(interface{})) *Session {
	if closures != nil {
		session.beforeClosures = append(session.beforeClosures, closures)
	}
	return session
}

// After Apply after Processor, affected bean is passed to closure arg
func (session *Session) After(closures func(interface{})) *Session {
	if closures != nil {
		session.afterClosures = append(session.afterClosures, closures)
	}
	return session
}

// Table can input a string or pointer to struct for special a table to operate.
func (session *Session) Table(tableNameOrBean interface{}) *Session {
	if err := session.statement.SetTable(tableNameOrBean); err != nil {
		session.statement.LastError = err
	}
	return session
}

// Alias set the table alias
func (session *Session) Alias(alias string) *Session {
	session.statement.Alias(alias)
	return session
}

// NoCascade indicate that no cascade load child object
func (session *Session) NoCascade() *Session {
	session.statement.UseCascade = false
	return session
}

// ForUpdate Set Read/Write locking for UPDATE
func (session *Session) ForUpdate() *Session {
	session.statement.IsForUpdate = true
	return session
}

// NoAutoCondition disable generate SQL condition from beans
func (session *Session) NoAutoCondition(no ...bool) *Session {
	session.statement.SetNoAutoCondition(no...)
	return session
}

// Limit provide limit and offset query condition
func (session *Session) Limit(limit int, start ...int) *Session {
	session.statement.Limit(limit, start...)
	return session
}

// OrderBy provide order by query condition, the input parameter is the content
// after order by on a sql statement.
func (session *Session) OrderBy(order string) *Session {
	session.statement.OrderBy(order)
	return session
}

// Desc provide desc order by query condition, the input parameters are columns.
func (session *Session) Desc(colNames ...string) *Session {
	session.statement.Desc(colNames...)
	return session
}

// Asc provide asc order by query condition, the input parameters are columns.
func (session *Session) Asc(colNames ...string) *Session {
	session.statement.Asc(colNames...)
	return session
}

// StoreEngine is only avialble mysql dialect currently
func (session *Session) StoreEngine(storeEngine string) *Session {
	session.statement.StoreEngine = storeEngine
	return session
}

// Charset is only avialble mysql dialect currently
func (session *Session) Charset(charset string) *Session {
	session.statement.Charset = charset
	return session
}

// Cascade indicates if loading sub Struct
func (session *Session) Cascade(trueOrFalse ...bool) *Session {
	if len(trueOrFalse) >= 1 {
		session.statement.UseCascade = trueOrFalse[0]
	}
	return session
}

// MustLogSQL means record SQL or not and don't follow engine's setting
func (session *Session) MustLogSQL(logs ...bool) *Session {
	var showSQL = true
	if len(logs) > 0 {
		showSQL = logs[0]
	}
	session.ctx = context.WithValue(session.ctx, log.SessionShowSQLKey, showSQL)
	return session
}

// NoCache ask this session do not retrieve data from cache system and
// get data from database directly.
func (session *Session) NoCache() *Session {
	session.statement.UseCache = false
	return session
}

// Join join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (session *Session) Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *Session {
	session.statement.Join(joinOperator, tablename, condition, args...)
	return session
}

// GroupBy Generate Group By statement
func (session *Session) GroupBy(keys string) *Session {
	session.statement.GroupBy(keys)
	return session
}

// Having Generate Having statement
func (session *Session) Having(conditions string) *Session {
	session.statement.Having(conditions)
	return session
}

// DB db return the wrapper of sql.DB
func (session *Session) DB() *core.DB {
	return session.db()
}

func (session *Session) canCache() bool {
	if session.statement.RefTable == nil ||
		session.statement.JoinStr != "" ||
		session.statement.RawSQL != "" ||
		!session.statement.UseCache ||
		session.statement.IsForUpdate ||
		session.tx != nil ||
		len(session.statement.SelectStr) > 0 {
		return false
	}
	return true
}

func (session *Session) doPrepare(db *core.DB, sqlStr string) (stmt *core.Stmt, err error) {
	crc := crc32.ChecksumIEEE([]byte(sqlStr))
	// TODO try hash(sqlStr+len(sqlStr))
	var has bool
	stmt, has = session.stmtCache[crc]
	if !has {
		stmt, err = db.PrepareContext(session.ctx, sqlStr)
		if err != nil {
			return nil, err
		}
		session.stmtCache[crc] = stmt
	}
	return
}

func (session *Session) getField(dataStruct *reflect.Value, key string, table *schemas.Table, idx int) (*reflect.Value, error) {
	var col *schemas.Column
	if col = table.GetColumnIdx(key, idx); col == nil {
		return nil, ErrFieldIsNotExist{key, table.Name}
	}

	fieldValue, err := col.ValueOfV(dataStruct)
	if err != nil {
		return nil, err
	}

	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return nil, ErrFieldIsNotValid{key, table.Name}
	}

	return fieldValue, nil
}

// Cell cell is a result of one column field
type Cell *interface{}

func (session *Session) rows2Beans(rows *core.Rows, fields []string,
	table *schemas.Table, newElemFunc func([]string) reflect.Value,
	sliceValueSetFunc func(*reflect.Value, schemas.PK) error) error {
	for rows.Next() {
		var newValue = newElemFunc(fields)
		bean := newValue.Interface()
		dataStruct := newValue.Elem()

		// handle beforeClosures
		scanResults, err := session.row2Slice(rows, fields, bean)
		if err != nil {
			return err
		}
		pk, err := session.slice2Bean(scanResults, fields, bean, &dataStruct, table)
		if err != nil {
			return err
		}
		session.afterProcessors = append(session.afterProcessors, executedProcessor{
			fun: func(*Session, interface{}) error {
				return sliceValueSetFunc(&newValue, pk)
			},
			session: session,
			bean:    bean,
		})
	}
	return nil
}

func (session *Session) row2Slice(rows *core.Rows, fields []string, bean interface{}) ([]interface{}, error) {
	for _, closure := range session.beforeClosures {
		closure(bean)
	}

	scanResults := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var cell interface{}
		scanResults[i] = &cell
	}
	if err := rows.Scan(scanResults...); err != nil {
		return nil, err
	}

	executeBeforeSet(bean, fields, scanResults)

	return scanResults, nil
}

func (session *Session) slice2Bean(scanResults []interface{}, fields []string, bean interface{}, dataStruct *reflect.Value, table *schemas.Table) (schemas.PK, error) {
	defer func() {
		executeAfterSet(bean, fields, scanResults)
	}()

	buildAfterProcessors(session, bean)

	var tempMap = make(map[string]int)
	var pk schemas.PK
	for ii, key := range fields {
		var idx int
		var ok bool
		var lKey = strings.ToLower(key)
		if idx, ok = tempMap[lKey]; !ok {
			idx = 0
		} else {
			idx = idx + 1
		}
		tempMap[lKey] = idx

		fieldValue, err := session.getField(dataStruct, key, table, idx)
		if err != nil {
			if !strings.Contains(err.Error(), "is not valid") {
				session.engine.logger.Warnf("%v", err)
			}
			continue
		}
		if fieldValue == nil {
			continue
		}
		rawValue := reflect.Indirect(reflect.ValueOf(scanResults[ii]))

		// if row is null then ignore
		if rawValue.Interface() == nil {
			continue
		}

		if fieldValue.CanAddr() {
			if structConvert, ok := fieldValue.Addr().Interface().(convert.Conversion); ok {
				if data, err := value2Bytes(&rawValue); err == nil {
					if err := structConvert.FromDB(data); err != nil {
						return nil, err
					}
				} else {
					return nil, err
				}
				continue
			}
		}

		if _, ok := fieldValue.Interface().(convert.Conversion); ok {
			if data, err := value2Bytes(&rawValue); err == nil {
				if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue.Interface().(convert.Conversion).FromDB(data)
			} else {
				return nil, err
			}
			continue
		}

		rawValueType := reflect.TypeOf(rawValue.Interface())
		vv := reflect.ValueOf(rawValue.Interface())
		col := table.GetColumnIdx(key, idx)
		if col.IsPrimaryKey {
			pk = append(pk, rawValue.Interface())
		}
		fieldType := fieldValue.Type()
		hasAssigned := false

		if col.SQLType.IsJson() {
			var bs []byte
			if rawValueType.Kind() == reflect.String {
				bs = []byte(vv.String())
			} else if rawValueType.ConvertibleTo(schemas.BytesType) {
				bs = vv.Bytes()
			} else {
				return nil, fmt.Errorf("unsupported database data type: %s %v", key, rawValueType.Kind())
			}

			hasAssigned = true

			if len(bs) > 0 {
				if fieldType.Kind() == reflect.String {
					fieldValue.SetString(string(bs))
					continue
				}
				if fieldValue.CanAddr() {
					err := json.DefaultJSONHandler.Unmarshal(bs, fieldValue.Addr().Interface())
					if err != nil {
						return nil, err
					}
				} else {
					x := reflect.New(fieldType)
					err := json.DefaultJSONHandler.Unmarshal(bs, x.Interface())
					if err != nil {
						return nil, err
					}
					fieldValue.Set(x.Elem())
				}
			}

			continue
		}

		switch fieldType.Kind() {
		case reflect.Complex64, reflect.Complex128:
			// TODO: reimplement this
			var bs []byte
			if rawValueType.Kind() == reflect.String {
				bs = []byte(vv.String())
			} else if rawValueType.ConvertibleTo(schemas.BytesType) {
				bs = vv.Bytes()
			}

			hasAssigned = true
			if len(bs) > 0 {
				if fieldValue.CanAddr() {
					err := json.DefaultJSONHandler.Unmarshal(bs, fieldValue.Addr().Interface())
					if err != nil {
						return nil, err
					}
				} else {
					x := reflect.New(fieldType)
					err := json.DefaultJSONHandler.Unmarshal(bs, x.Interface())
					if err != nil {
						return nil, err
					}
					fieldValue.Set(x.Elem())
				}
			}
		case reflect.Slice, reflect.Array:
			switch rawValueType.Kind() {
			case reflect.Slice, reflect.Array:
				switch rawValueType.Elem().Kind() {
				case reflect.Uint8:
					if fieldType.Elem().Kind() == reflect.Uint8 {
						hasAssigned = true
						if col.SQLType.IsText() {
							x := reflect.New(fieldType)
							err := json.DefaultJSONHandler.Unmarshal(vv.Bytes(), x.Interface())
							if err != nil {
								return nil, err
							}
							fieldValue.Set(x.Elem())
						} else {
							if fieldValue.Len() > 0 {
								for i := 0; i < fieldValue.Len(); i++ {
									if i < vv.Len() {
										fieldValue.Index(i).Set(vv.Index(i))
									}
								}
							} else {
								for i := 0; i < vv.Len(); i++ {
									fieldValue.Set(reflect.Append(*fieldValue, vv.Index(i)))
								}
							}
						}
					}
				}
			}
		case reflect.String:
			if rawValueType.Kind() == reflect.String {
				hasAssigned = true
				fieldValue.SetString(vv.String())
			}
		case reflect.Bool:
			if rawValueType.Kind() == reflect.Bool {
				hasAssigned = true
				fieldValue.SetBool(vv.Bool())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch rawValueType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				hasAssigned = true
				fieldValue.SetInt(vv.Int())
			}
		case reflect.Float32, reflect.Float64:
			switch rawValueType.Kind() {
			case reflect.Float32, reflect.Float64:
				hasAssigned = true
				fieldValue.SetFloat(vv.Float())
			}
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			switch rawValueType.Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				hasAssigned = true
				fieldValue.SetUint(vv.Uint())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				hasAssigned = true
				fieldValue.SetUint(uint64(vv.Int()))
			}
		case reflect.Struct:
			if fieldType.ConvertibleTo(schemas.TimeType) {
				dbTZ := session.engine.DatabaseTZ
				if col.TimeZone != nil {
					dbTZ = col.TimeZone
				}

				if rawValueType == schemas.TimeType {
					hasAssigned = true

					t := vv.Convert(schemas.TimeType).Interface().(time.Time)

					z, _ := t.Zone()
					// set new location if database don't save timezone or give an incorrect timezone
					if len(z) == 0 || t.Year() == 0 || t.Location().String() != dbTZ.String() { // !nashtsai! HACK tmp work around for lib/pq doesn't properly time with location
						session.engine.logger.Debugf("empty zone key[%v] : %v | zone: %v | location: %+v\n", key, t, z, *t.Location())
						t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(),
							t.Minute(), t.Second(), t.Nanosecond(), dbTZ)
					}

					t = t.In(session.engine.TZLocation)
					fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
				} else if rawValueType == schemas.IntType || rawValueType == schemas.Int64Type ||
					rawValueType == schemas.Int32Type {
					hasAssigned = true

					t := time.Unix(vv.Int(), 0).In(session.engine.TZLocation)
					fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
				} else {
					if d, ok := vv.Interface().([]uint8); ok {
						hasAssigned = true
						t, err := session.byte2Time(col, d)
						if err != nil {
							session.engine.logger.Errorf("byte2Time error: %v", err)
							hasAssigned = false
						} else {
							fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
						}
					} else if d, ok := vv.Interface().(string); ok {
						hasAssigned = true
						t, err := session.str2Time(col, d)
						if err != nil {
							session.engine.logger.Errorf("byte2Time error: %v", err)
							hasAssigned = false
						} else {
							fieldValue.Set(reflect.ValueOf(t).Convert(fieldType))
						}
					} else {
						return nil, fmt.Errorf("rawValueType is %v, value is %v", rawValueType, vv.Interface())
					}
				}
			} else if nulVal, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
				// !<winxxp>! 增加支持sql.Scanner接口的结构，如sql.NullString
				hasAssigned = true
				if err := nulVal.Scan(vv.Interface()); err != nil {
					session.engine.logger.Errorf("sql.Sanner error: %v", err)
					hasAssigned = false
				}
			} else if col.SQLType.IsJson() {
				if rawValueType.Kind() == reflect.String {
					hasAssigned = true
					x := reflect.New(fieldType)
					if len([]byte(vv.String())) > 0 {
						err := json.DefaultJSONHandler.Unmarshal([]byte(vv.String()), x.Interface())
						if err != nil {
							return nil, err
						}
						fieldValue.Set(x.Elem())
					}
				} else if rawValueType.Kind() == reflect.Slice {
					hasAssigned = true
					x := reflect.New(fieldType)
					if len(vv.Bytes()) > 0 {
						err := json.DefaultJSONHandler.Unmarshal(vv.Bytes(), x.Interface())
						if err != nil {
							return nil, err
						}
						fieldValue.Set(x.Elem())
					}
				}
			} else if session.statement.UseCascade {
				table, err := session.engine.tagParser.ParseWithCache(*fieldValue)
				if err != nil {
					return nil, err
				}

				hasAssigned = true
				if len(table.PrimaryKeys) != 1 {
					return nil, errors.New("unsupported non or composited primary key cascade")
				}
				var pk = make(schemas.PK, len(table.PrimaryKeys))
				pk[0], err = asKind(vv, rawValueType)
				if err != nil {
					return nil, err
				}

				if !pk.IsZero() {
					// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
					// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
					// property to be fetched lazily
					structInter := reflect.New(fieldValue.Type())
					has, err := session.ID(pk).NoCascade().get(structInter.Interface())
					if err != nil {
						return nil, err
					}
					if has {
						fieldValue.Set(structInter.Elem())
					} else {
						return nil, errors.New("cascade obj is not exist")
					}
				}
			}
		case reflect.Ptr:
			// !nashtsai! TODO merge duplicated codes above
			switch fieldType {
			// following types case matching ptr's native type, therefore assign ptr directly
			case schemas.PtrStringType:
				if rawValueType.Kind() == reflect.String {
					x := vv.String()
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrBoolType:
				if rawValueType.Kind() == reflect.Bool {
					x := vv.Bool()
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrTimeType:
				if rawValueType == schemas.PtrTimeType {
					hasAssigned = true
					var x = rawValue.Interface().(time.Time)
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrFloat64Type:
				if rawValueType.Kind() == reflect.Float64 {
					x := vv.Float()
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrUint64Type:
				if rawValueType.Kind() == reflect.Int64 {
					var x = uint64(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrInt64Type:
				if rawValueType.Kind() == reflect.Int64 {
					x := vv.Int()
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrFloat32Type:
				if rawValueType.Kind() == reflect.Float64 {
					var x = float32(vv.Float())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrIntType:
				if rawValueType.Kind() == reflect.Int64 {
					var x = int(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrInt32Type:
				if rawValueType.Kind() == reflect.Int64 {
					var x = int32(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrInt8Type:
				if rawValueType.Kind() == reflect.Int64 {
					var x = int8(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrInt16Type:
				if rawValueType.Kind() == reflect.Int64 {
					var x = int16(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrUintType:
				if rawValueType.Kind() == reflect.Int64 {
					var x = uint(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.PtrUint32Type:
				if rawValueType.Kind() == reflect.Int64 {
					var x = uint32(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.Uint8Type:
				if rawValueType.Kind() == reflect.Int64 {
					var x = uint8(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.Uint16Type:
				if rawValueType.Kind() == reflect.Int64 {
					var x = uint16(vv.Int())
					hasAssigned = true
					fieldValue.Set(reflect.ValueOf(&x))
				}
			case schemas.Complex64Type:
				var x complex64
				if len([]byte(vv.String())) > 0 {
					err := json.DefaultJSONHandler.Unmarshal([]byte(vv.String()), &x)
					if err != nil {
						return nil, err
					}
					fieldValue.Set(reflect.ValueOf(&x))
				}
				hasAssigned = true
			case schemas.Complex128Type:
				var x complex128
				if len([]byte(vv.String())) > 0 {
					err := json.DefaultJSONHandler.Unmarshal([]byte(vv.String()), &x)
					if err != nil {
						return nil, err
					}
					fieldValue.Set(reflect.ValueOf(&x))
				}
				hasAssigned = true
			} // switch fieldType
		} // switch fieldType.Kind()

		// !nashtsai! for value can't be assigned directly fallback to convert to []byte then back to value
		if !hasAssigned {
			data, err := value2Bytes(&rawValue)
			if err != nil {
				return nil, err
			}

			if err = session.bytes2Value(col, fieldValue, data); err != nil {
				return nil, err
			}
		}
	}
	return pk, nil
}

// saveLastSQL stores executed query information
func (session *Session) saveLastSQL(sql string, args ...interface{}) {
	session.lastSQL = sql
	session.lastSQLArgs = args
	session.logSQL(sql, args...)
}

func (session *Session) logSQL(sqlStr string, sqlArgs ...interface{}) {
	if session.showSQL {
		if len(sqlArgs) > 0 {
			session.engine.logger.Infof("[SQL] %v %#v", sqlStr, sqlArgs)
		} else {
			session.engine.logger.Infof("[SQL] %v", sqlStr)
		}
	}
}

// LastSQL returns last query information
func (session *Session) LastSQL() (string, []interface{}) {
	return session.lastSQL, session.lastSQLArgs
}

// Unscoped always disable struct tag "deleted"
func (session *Session) Unscoped() *Session {
	session.statement.SetUnscoped()
	return session
}

func (session *Session) incrVersionFieldValue(fieldValue *reflect.Value) {
	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fieldValue.SetInt(fieldValue.Int() + 1)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fieldValue.SetUint(fieldValue.Uint() + 1)
	}
}

// ContextHook sets the context on this session
func (session *Session) Context(ctx context.Context) *Session {
	session.ctx = ctx
	return session
}

// PingContext test if database is ok
func (session *Session) PingContext(ctx context.Context) error {
	if session.isAutoClose {
		defer session.Close()
	}

	session.engine.logger.Infof("PING DATABASE %v", session.engine.DriverName())
	return session.DB().PingContext(ctx)
}
