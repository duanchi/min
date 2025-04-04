// Copyright 2019 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statements

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/duanchi/min/v2/db/xorm/core"
	"github.com/duanchi/min/v2/db/xorm/dialects"
	"github.com/duanchi/min/v2/db/xorm/private/utils"
	"github.com/duanchi/min/v2/db/xorm/schemas"
	"github.com/xormplus/builder"
)

func (statement *Statement) genSelectSql(dialect dialects.Dialect, rownumber string) string {

	var sql = statement.RawSQL
	var orderBys = statement.OrderStr
	pLimitN := statement.LimitN

	if dialect.URI().DBType != schemas.MSSQL && dialect.URI().DBType != schemas.ORACLE {
		if statement.Start > 0 {
			sql = fmt.Sprintf("%v LIMIT %v OFFSET %v", sql, statement.LimitN, statement.Start)
			if pLimitN != nil {
				sql = fmt.Sprintf("%v LIMIT %v OFFSET %v", sql, *pLimitN, statement.Start)
			} else {
				sql = fmt.Sprintf("%v LIMIT 0 OFFSET %v", sql, *pLimitN)
			}
		} else if pLimitN != nil {
			sql = fmt.Sprintf("%v LIMIT %v", sql, statement.LimitN)
		}
	} else if dialect.URI().DBType == schemas.ORACLE {
		if statement.Start != 0 || pLimitN != nil {
			sql = fmt.Sprintf("SELECT aat.* FROM (SELECT at.*,ROWNUM %v FROM (%v) at WHERE ROWNUM <= %d) aat WHERE %v > %d",
				rownumber, sql, statement.Start+*pLimitN, rownumber, statement.Start)
		}
	} else {
		keepSelect := false
		var fullQuery string
		if statement.Start > 0 {
			fullQuery = fmt.Sprintf("SELECT sq.* FROM (SELECT ROW_NUMBER() OVER (ORDER BY %v) AS %v,", orderBys, rownumber)
		} else if pLimitN != nil {
			fullQuery = fmt.Sprintf("SELECT TOP %d", *pLimitN)
		} else {
			keepSelect = true
		}

		if !keepSelect {
			expr := `^\s*SELECT\s*`
			reg, err := regexp.Compile(expr)
			if err != nil {
				fmt.Println(err)
			}
			sql = strings.ToUpper(sql)
			if reg.MatchString(sql) {
				str := reg.FindAllString(sql, -1)
				fullQuery = fmt.Sprintf("%v %v", fullQuery, sql[len(str[0]):])
			}
		}

		if statement.Start > 0 {
			// T-SQL offset starts with 1, not like MySQL with 0;
			if pLimitN != nil {
				fullQuery = fmt.Sprintf("%v) AS sq WHERE %v BETWEEN %d AND %d", fullQuery, rownumber, statement.Start+1, statement.Start+*pLimitN)
			} else {
				fullQuery = fmt.Sprintf("%v) AS sq WHERE %v >= %d", fullQuery, rownumber, statement.Start+1)
			}
		} else {
			fullQuery = fmt.Sprintf("%v ORDER BY %v", fullQuery, orderBys)
		}

		if keepSelect {
			if len(orderBys) > 0 {
				sql = fmt.Sprintf("%v ORDER BY %v", sql, orderBys)
			}
		} else {
			sql = fullQuery
		}
	}

	return sql
}

func (statement *Statement) GenQuerySQL(sqlOrArgs ...interface{}) (string, []interface{}, error) {
	if len(sqlOrArgs) > 0 {
		return statement.ConvertSQLOrArgs(sqlOrArgs...)
	}

	if statement.RawSQL != "" {
		var dialect = statement.dialect
		rownumber := "xorm" + utils.NewShortUUID().String()
		sql := statement.genSelectSql(dialect, rownumber)

		params := statement.RawParams
		i := len(params)

		//		var result []map[string]interface{}
		//		var err error
		if i == 1 {
			vv := reflect.ValueOf(params[0])
			if vv.Kind() != reflect.Ptr || vv.Elem().Kind() != reflect.Map {
				return sql, params, nil
			} else {
				sqlStr1, param, _ := core.MapToSlice(sql, params[0])
				return sqlStr1, param, nil
			}
		} else {
			return sql, params, nil
		}
		//		return session.statement.RawSQL, session.statement.RawParams, nil
	}

	if len(statement.TableName()) <= 0 {
		return "", nil, ErrTableNotFound
	}

	var columnStr = statement.ColumnStr()
	if len(statement.SelectStr) > 0 {
		columnStr = statement.SelectStr
	} else {
		if statement.JoinStr == "" {
			if columnStr == "" {
				if statement.GroupByStr != "" {
					columnStr = statement.quoteColumnStr(statement.GroupByStr)
				} else {
					columnStr = statement.genColumnStr()
				}
			}
		} else {
			if columnStr == "" {
				if statement.GroupByStr != "" {
					columnStr = statement.quoteColumnStr(statement.GroupByStr)
				} else {
					columnStr = "*"
				}
			}
		}
		if columnStr == "" {
			columnStr = "*"
		}
	}

	if err := statement.ProcessIDParam(); err != nil {
		return "", nil, err
	}

	sqlStr, condArgs, err := statement.genSelectSQL(columnStr, true, true)
	if err != nil {
		return "", nil, err
	}
	args := append(statement.joinArgs, condArgs...)

	// for mssql and use limit
	qs := strings.Count(sqlStr, "?")
	if len(args)*2 == qs {
		args = append(args, args...)
	}

	return sqlStr, args, nil
}

func (statement *Statement) GenSumSQL(bean interface{}, columns ...string) (string, []interface{}, error) {
	if statement.RawSQL != "" {
		return statement.GenRawSQL(), statement.RawParams, nil
	}

	statement.SetRefBean(bean)

	var sumStrs = make([]string, 0, len(columns))
	for _, colName := range columns {
		if !strings.Contains(colName, " ") && !strings.Contains(colName, "(") {
			colName = statement.quote(colName)
		} else {
			colName = statement.ReplaceQuote(colName)
		}
		sumStrs = append(sumStrs, fmt.Sprintf("COALESCE(sum(%s),0)", colName))
	}
	sumSelect := strings.Join(sumStrs, ", ")

	if err := statement.mergeConds(bean); err != nil {
		return "", nil, err
	}

	sqlStr, condArgs, err := statement.genSelectSQL(sumSelect, true, true)
	if err != nil {
		return "", nil, err
	}

	return sqlStr, append(statement.joinArgs, condArgs...), nil
}

func (statement *Statement) GenGetSQL(bean interface{}) (string, []interface{}, error) {
	v := rValue(bean)
	isStruct := v.Kind() == reflect.Struct
	if isStruct {
		statement.SetRefBean(bean)
	}

	var columnStr = statement.ColumnStr()
	if len(statement.SelectStr) > 0 {
		columnStr = statement.SelectStr
	} else {
		// TODO: always generate column names, not use * even if join
		if len(statement.JoinStr) == 0 {
			if len(columnStr) == 0 {
				if len(statement.GroupByStr) > 0 {
					columnStr = statement.quoteColumnStr(statement.GroupByStr)
				} else {
					columnStr = statement.genColumnStr()
				}
			}
		} else {
			if len(columnStr) == 0 {
				if len(statement.GroupByStr) > 0 {
					columnStr = statement.quoteColumnStr(statement.GroupByStr)
				}
			}
		}
	}

	if len(columnStr) == 0 {
		columnStr = "*"
	}

	if isStruct {
		if err := statement.mergeConds(bean); err != nil {
			return "", nil, err
		}
	} else {
		if err := statement.ProcessIDParam(); err != nil {
			return "", nil, err
		}
	}

	sqlStr, condArgs, err := statement.genSelectSQL(columnStr, true, true)
	if err != nil {
		return "", nil, err
	}

	return sqlStr, append(statement.joinArgs, condArgs...), nil
}

// GenCountSQL generates the SQL for counting
func (statement *Statement) GenCountSQL(beans ...interface{}) (string, []interface{}, error) {
	if statement.RawSQL != "" {
		return statement.GenRawSQL(), statement.RawParams, nil
	}

	var condArgs []interface{}
	var err error
	if len(beans) > 0 {
		statement.SetRefBean(beans[0])
		if err := statement.mergeConds(beans[0]); err != nil {
			return "", nil, err
		}
	}

	var selectSQL = statement.SelectStr
	if len(selectSQL) <= 0 {
		if statement.IsDistinct {
			selectSQL = fmt.Sprintf("count(DISTINCT %s)", statement.ColumnStr())
		} else if statement.ColumnStr() != "" {
			selectSQL = fmt.Sprintf("count(%s)", statement.ColumnStr())
		} else {
			selectSQL = "count(*)"
		}
	}
	sqlStr, condArgs, err := statement.genSelectSQL(selectSQL, false, false)
	if err != nil {
		return "", nil, err
	}

	return sqlStr, append(statement.joinArgs, condArgs...), nil
}

func (statement *Statement) genSelectSQL(columnStr string, needLimit, needOrderBy bool) (string, []interface{}, error) {
	var (
		distinct                  string
		dialect                   = statement.dialect
		quote                     = statement.quote
		fromStr                   = " FROM "
		top, mssqlCondi, whereStr string
	)
	if statement.IsDistinct && !strings.HasPrefix(columnStr, "count") {
		distinct = "DISTINCT "
	}

	condSQL, condArgs, err := statement.GenCondSQL(statement.cond)
	if err != nil {
		return "", nil, err
	}
	if len(condSQL) > 0 {
		whereStr = " WHERE " + condSQL
	}

	if dialect.URI().DBType == schemas.MSSQL && strings.Contains(statement.TableName(), "..") {
		fromStr += statement.TableName()
	} else {
		fromStr += quote(statement.TableName())
	}

	if statement.TableAlias != "" {
		if dialect.URI().DBType == schemas.ORACLE {
			fromStr += " " + quote(statement.TableAlias)
		} else {
			fromStr += " AS " + quote(statement.TableAlias)
		}
	}
	if statement.JoinStr != "" {
		fromStr = fmt.Sprintf("%v %v", fromStr, statement.JoinStr)
	}

	pLimitN := statement.LimitN
	if dialect.URI().DBType == schemas.MSSQL {
		if pLimitN != nil {
			LimitNValue := *pLimitN
			top = fmt.Sprintf("TOP %d ", LimitNValue)
		}
		if statement.Start > 0 {
			var column string
			if len(statement.RefTable.PKColumns()) == 0 {
				for _, index := range statement.RefTable.Indexes {
					if len(index.Cols) == 1 {
						column = index.Cols[0]
						break
					}
				}
				if len(column) == 0 {
					column = statement.RefTable.ColumnsSeq()[0]
				}
			} else {
				column = statement.RefTable.PKColumns()[0].Name
			}
			if statement.needTableName() {
				if len(statement.TableAlias) > 0 {
					column = statement.TableAlias + "." + column
				} else {
					column = statement.TableName() + "." + column
				}
			}

			var orderStr string
			if needOrderBy && len(statement.OrderStr) > 0 {
				orderStr = " ORDER BY " + statement.OrderStr
			}

			var groupStr string
			if len(statement.GroupByStr) > 0 {
				groupStr = " GROUP BY " + statement.GroupByStr
			}
			mssqlCondi = fmt.Sprintf("(%s NOT IN (SELECT TOP %d %s%s%s%s%s))",
				column, statement.Start, column, fromStr, whereStr, orderStr, groupStr)
		}
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, "SELECT %v%v%v%v%v", distinct, top, columnStr, fromStr, whereStr)
	if len(mssqlCondi) > 0 {
		if len(whereStr) > 0 {
			fmt.Fprint(&buf, " AND ", mssqlCondi)
		} else {
			fmt.Fprint(&buf, " WHERE ", mssqlCondi)
		}
	}

	if statement.GroupByStr != "" {
		fmt.Fprint(&buf, " GROUP BY ", statement.GroupByStr)
	}
	if statement.HavingStr != "" {
		fmt.Fprint(&buf, " ", statement.HavingStr)
	}
	if needOrderBy && statement.OrderStr != "" {
		fmt.Fprint(&buf, " ORDER BY ", statement.OrderStr)
	}
	if needLimit {
		if dialect.URI().DBType != schemas.MSSQL && dialect.URI().DBType != schemas.ORACLE {
			if statement.Start > 0 {
				if pLimitN != nil {
					fmt.Fprintf(&buf, " LIMIT %v OFFSET %v", *pLimitN, statement.Start)
				} else {
					fmt.Fprintf(&buf, "LIMIT 0 OFFSET %v", statement.Start)
				}
			} else if pLimitN != nil {
				fmt.Fprint(&buf, " LIMIT ", *pLimitN)
			}
		} else if dialect.URI().DBType == schemas.ORACLE {
			if statement.Start != 0 || pLimitN != nil {
				oldString := buf.String()
				buf.Reset()
				rawColStr := columnStr
				if rawColStr == "*" {
					rawColStr = "at.*"
				}
				fmt.Fprintf(&buf, "SELECT %v FROM (SELECT %v,ROWNUM RN FROM (%v) at WHERE ROWNUM <= %d) aat WHERE RN > %d",
					columnStr, rawColStr, oldString, statement.Start+*pLimitN, statement.Start)
			}
		}
	}
	if statement.IsForUpdate {
		return dialect.ForUpdateSQL(buf.String()), condArgs, nil
	}

	return buf.String(), condArgs, nil
}

func (statement *Statement) GenExistSQL(bean ...interface{}) (string, []interface{}, error) {
	if statement.RawSQL != "" {
		return statement.GenRawSQL(), statement.RawParams, nil
	}

	var sqlStr string
	var args []interface{}
	var joinStr string
	var err error
	if len(bean) == 0 {
		tableName := statement.TableName()
		if len(tableName) <= 0 {
			return "", nil, ErrTableNotFound
		}

		tableName = statement.quote(tableName)
		if len(statement.JoinStr) > 0 {
			joinStr = statement.JoinStr
		}

		if statement.Conds().IsValid() {
			condSQL, condArgs, err := statement.GenCondSQL(statement.Conds())
			if err != nil {
				return "", nil, err
			}

			if statement.dialect.URI().DBType == schemas.MSSQL {
				sqlStr = fmt.Sprintf("SELECT TOP 1 * FROM %s %s WHERE %s", tableName, joinStr, condSQL)
			} else if statement.dialect.URI().DBType == schemas.ORACLE {
				sqlStr = fmt.Sprintf("SELECT * FROM %s WHERE (%s) %s AND ROWNUM=1", tableName, joinStr, condSQL)
			} else {
				sqlStr = fmt.Sprintf("SELECT * FROM %s %s WHERE %s LIMIT 1", tableName, joinStr, condSQL)
			}
			args = condArgs
		} else {
			if statement.dialect.URI().DBType == schemas.MSSQL {
				sqlStr = fmt.Sprintf("SELECT TOP 1 * FROM %s %s", tableName, joinStr)
			} else if statement.dialect.URI().DBType == schemas.ORACLE {
				sqlStr = fmt.Sprintf("SELECT * FROM  %s %s WHERE ROWNUM=1", tableName, joinStr)
			} else {
				sqlStr = fmt.Sprintf("SELECT * FROM %s %s LIMIT 1", tableName, joinStr)
			}
			args = []interface{}{}
		}
	} else {
		beanValue := reflect.ValueOf(bean[0])
		if beanValue.Kind() != reflect.Ptr {
			return "", nil, errors.New("needs a pointer")
		}

		if beanValue.Elem().Kind() == reflect.Struct {
			if err := statement.SetRefBean(bean[0]); err != nil {
				return "", nil, err
			}
		}

		if len(statement.TableName()) <= 0 {
			return "", nil, ErrTableNotFound
		}
		statement.Limit(1)
		sqlStr, args, err = statement.GenGetSQL(bean[0])
		if err != nil {
			return "", nil, err
		}
	}

	return sqlStr, args, nil
}

func (statement *Statement) GenFindSQL(autoCond builder.Cond) (string, []interface{}, error) {
	if statement.RawSQL != "" {
		return statement.GenRawSQL(), statement.RawParams, nil
	}

	var sqlStr string
	var args []interface{}
	var err error

	if len(statement.TableName()) <= 0 {
		return "", nil, ErrTableNotFound
	}

	var columnStr = statement.ColumnStr()
	if len(statement.SelectStr) > 0 {
		columnStr = statement.SelectStr
	} else {
		if statement.JoinStr == "" {
			if columnStr == "" {
				if statement.GroupByStr != "" {
					columnStr = statement.quoteColumnStr(statement.GroupByStr)
				} else {
					columnStr = statement.genColumnStr()
				}
			}
		} else {
			if columnStr == "" {
				if statement.GroupByStr != "" {
					columnStr = statement.quoteColumnStr(statement.GroupByStr)
				} else {
					columnStr = "*"
				}
			}
		}
		if columnStr == "" {
			columnStr = "*"
		}
	}

	statement.cond = statement.cond.And(autoCond)

	sqlStr, condArgs, err := statement.genSelectSQL(columnStr, true, true)
	if err != nil {
		return "", nil, err
	}
	args = append(statement.joinArgs, condArgs...)
	// for mssql and use limit
	qs := strings.Count(sqlStr, "?")
	if len(args)*2 == qs {
		args = append(args, args...)
	}

	return sqlStr, args, nil
}
