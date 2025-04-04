// Copyright 2016 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"database/sql"
	"reflect"

	"github.com/duanchi/min/v2/db/xorm/core"
	"github.com/xormplus/builder"
	// "github.com/duanchi/min/v2/db/xorm/private/statements"
)

func (session *Session) queryPreprocess(sqlStr *string, paramStr ...interface{}) {
	for _, filter := range session.engine.dialect.Filters() {
		*sqlStr = filter.Do(*sqlStr)
	}

	session.lastSQL = *sqlStr
	session.lastSQLArgs = paramStr
}

func (session *Session) queryRows(sqlStr string, args ...interface{}) (*core.Rows, error) {
	defer session.resetStatement()
	if session.statement.LastError != nil {
		return nil, session.statement.LastError
	}

	session.queryPreprocess(&sqlStr, args...)

	session.lastSQL = sqlStr
	session.lastSQLArgs = args
	if session.showSQL {
		if len(args) > 0 {
			session.engine.logger.Infof("[SQL][%p] %v %#v", session, sqlStr, args)
		} else {
			session.engine.logger.Infof("[SQL][%p] %v", session, sqlStr)
		}
	}

	if session.isAutoCommit {
		var db *core.DB
		if session.sessionType == groupSession {
			db = session.engine.engineGroup.Subordinate().DB()
		} else {
			db = session.DB()
		}

		if session.prepareStmt {
			// don't clear stmt since session will cache them
			stmt, err := session.doPrepare(db, sqlStr)
			if err != nil {
				return nil, err
			}

			rows, err := stmt.QueryContext(session.ctx, args...)
			if err != nil {
				return nil, err
			}
			return rows, nil
		}

		rows, err := db.QueryContext(session.ctx, sqlStr, args...)
		if err != nil {
			return nil, err
		}
		return rows, nil
	}

	rows, err := session.tx.QueryContext(session.ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (session *Session) queryRow(sqlStr string, args ...interface{}) *core.Row {
	return core.NewRow(session.queryRows(sqlStr, args...))
}

func value2Bytes(rawValue *reflect.Value) ([]byte, error) {
	str, err := value2String(rawValue)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

func value2Value(rawValue *reflect.Value) (Value, error) {
	str, err := value2String(rawValue)
	if err != nil {
		return nil, err
	}
	return Value(str), nil
}

func row2map(rows *core.Rows, fields []string) (resultsMap map[string][]byte, err error) {
	result := make(map[string][]byte)
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			result[key] = []byte{}
			continue
		}

		if data, err := value2Bytes(&rawValue); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}
	}
	return result, nil
}

func row2mapValue(rows *core.Rows, fields []string) (resultsMap map[string]Value, err error) {
	result := make(map[string]Value)
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		if rawValue.Interface() == nil {
			result[key] = nil
			continue
		}

		if data, err := value2Value(&rawValue); err == nil {
			result[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}
	}
	return result, nil
}

func row2Record(rows *core.Rows, fields []string) (record Record, err error) {
	record = make(Record)
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		if rawValue.Interface() == nil {
			record[key] = nil
			continue
		}

		if data, err := value2Value(&rawValue); err == nil {
			record[key] = data
		} else {
			return nil, err // !nashtsai! REVIEW, should return err or just error log?
		}
	}
	return record, nil
}

func rows2maps(rows *core.Rows) (resultsSlice []map[string][]byte, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2map(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func rows2mapsValue(rows *core.Rows) (resultsSlice []map[string]Value, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2mapValue(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func rows2Result(rows *core.Rows) (result Result, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		r, err := row2Record(rows, fields)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

func (session *Session) queryBytes(sqlStr string, args ...interface{}) ([]map[string][]byte, error) {
	rows, err := session.queryRows(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows2maps(rows)
}

func (session *Session) queryValue(sqlStr string, args ...interface{}) ([]map[string]Value, error) {
	rows, err := session.queryRows(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows2mapsValue(rows)
}

func (session *Session) queryResult(sqlStr string, args ...interface{}) (Result, error) {
	rows, err := session.queryRows(sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows2Result(rows)
}

func (session *Session) exec(sqlStr string, args ...interface{}) (sql.Result, error) {
	defer session.resetStatement()

	session.queryPreprocess(&sqlStr, args...)

	session.lastSQL = sqlStr
	session.lastSQLArgs = args

	if session.showSQL {
		if len(args) > 0 {
			session.engine.logger.Infof("[SQL][%p] %v %#v", session, sqlStr, args)
		} else {
			session.engine.logger.Infof("[SQL][%p] %v", session, sqlStr)
		}
	}

	if !session.isAutoCommit {
		return session.tx.ExecContext(session.ctx, sqlStr, args...)
	}

	if session.prepareStmt {
		stmt, err := session.doPrepare(session.DB(), sqlStr)
		if err != nil {
			return nil, err
		}

		res, err := stmt.ExecContext(session.ctx, args...)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return session.DB().ExecContext(session.ctx, sqlStr, args...)
}

func convertSQLOrArgs(sqlOrArgs ...interface{}) (string, []interface{}, error) {
	switch sqlOrArgs[0].(type) {
	case string:
		return sqlOrArgs[0].(string), sqlOrArgs[1:], nil
	case *builder.Builder:
		return sqlOrArgs[0].(*builder.Builder).ToSQL()
	case builder.Builder:
		bd := sqlOrArgs[0].(builder.Builder)
		return bd.ToSQL()
	}

	return "", nil, ErrUnSupportedType
}

// Exec raw sql
func (session *Session) Exec(sqlOrArgs ...interface{}) (sql.Result, error) {
	if session.isAutoClose {
		defer session.Close()
	}

	if len(sqlOrArgs) == 0 {
		return nil, ErrUnSupportedType
	}

	sqlStr, args, err := convertSQLOrArgs(sqlOrArgs...)
	if err != nil {
		return nil, err
	}

	return session.exec(sqlStr, args...)
}
