// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/duanchi/min/v2/db/xorm/convert"
	"github.com/duanchi/min/v2/db/xorm/dialects"
	"github.com/duanchi/min/v2/db/xorm/private/json"
	"github.com/duanchi/min/v2/db/xorm/private/utils"
	"github.com/duanchi/min/v2/db/xorm/schemas"
)

func (session *Session) str2Time(col *schemas.Column, data string) (outTime time.Time, outErr error) {
	sdata := strings.TrimSpace(data)
	var x time.Time
	var err error

	var parseLoc = session.engine.DatabaseTZ
	if col.TimeZone != nil {
		parseLoc = col.TimeZone
	}

	if sdata == utils.ZeroTime0 || sdata == utils.ZeroTime1 {
	} else if !strings.ContainsAny(sdata, "- :") { // !nashtsai! has only found that mymysql driver is using this for time type column
		// time stamp
		sd, err := strconv.ParseInt(sdata, 10, 64)
		if err == nil {
			x = time.Unix(sd, 0)
			//session.engine.logger.Debugf("time(0) key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
		} else {
			//session.engine.logger.Debugf("time(0) err key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
		}
	} else if len(sdata) > 19 && strings.Contains(sdata, "-") {
		x, err = time.ParseInLocation(time.RFC3339Nano, sdata, parseLoc)
		session.engine.logger.Debugf("time(1) key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
		if err != nil {
			x, err = time.ParseInLocation("2006-01-02 15:04:05.999999999", sdata, parseLoc)
			//session.engine.logger.Debugf("time(2) key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
		}
		if err != nil {
			x, err = time.ParseInLocation("2006-01-02 15:04:05.9999999 Z07:00", sdata, parseLoc)
			//session.engine.logger.Debugf("time(3) key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
		}
	} else if len(sdata) == 19 && strings.Contains(sdata, "-") {
		x, err = time.ParseInLocation("2006-01-02 15:04:05", sdata, parseLoc)
		//session.engine.logger.Debugf("time(4) key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
	} else if len(sdata) == 10 && sdata[4] == '-' && sdata[7] == '-' {
		x, err = time.ParseInLocation("2006-01-02", sdata, parseLoc)
		//session.engine.logger.Debugf("time(5) key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
	} else if col.SQLType.Name == schemas.Time {
		if strings.Contains(sdata, " ") {
			ssd := strings.Split(sdata, " ")
			sdata = ssd[1]
		}

		sdata = strings.TrimSpace(sdata)
		if session.engine.dialect.URI().DBType == schemas.MYSQL && len(sdata) > 8 {
			sdata = sdata[len(sdata)-8:]
		}

		st := fmt.Sprintf("2006-01-02 %v", sdata)
		x, err = time.ParseInLocation("2006-01-02 15:04:05", st, parseLoc)
		//session.engine.logger.Debugf("time(6) key[%v]: %+v | sdata: [%v]\n", col.FieldName, x, sdata)
	} else {
		outErr = fmt.Errorf("unsupported time format %v", sdata)
		return
	}
	if err != nil {
		outErr = fmt.Errorf("unsupported time format %v: %v", sdata, err)
		return
	}
	outTime = x.In(session.engine.TZLocation)
	return
}

func (session *Session) byte2Time(col *schemas.Column, data []byte) (outTime time.Time, outErr error) {
	return session.str2Time(col, string(data))
}

var (
	nullFloatType = reflect.TypeOf(sql.NullFloat64{})
)

// convert a db data([]byte) to a field value
func (session *Session) bytes2Value(col *schemas.Column, fieldValue *reflect.Value, data []byte) error {
	if structConvert, ok := fieldValue.Addr().Interface().(convert.Conversion); ok {
		return structConvert.FromDB(data)
	}

	if structConvert, ok := fieldValue.Interface().(convert.Conversion); ok {
		return structConvert.FromDB(data)
	}

	var v interface{}
	key := col.Name
	fieldType := fieldValue.Type()

	switch fieldType.Kind() {
	case reflect.Complex64, reflect.Complex128:
		x := reflect.New(fieldType)
		if len(data) > 0 {
			err := json.DefaultJSONHandler.Unmarshal(data, x.Interface())
			if err != nil {
				session.engine.logger.Errorf("%v", err)
				return err
			}
			fieldValue.Set(x.Elem())
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		v = data
		t := fieldType.Elem()
		k := t.Kind()
		if col.SQLType.IsText() {
			x := reflect.New(fieldType)
			if len(data) > 0 {
				err := json.DefaultJSONHandler.Unmarshal(data, x.Interface())
				if err != nil {
					session.engine.logger.Errorf("%v", err)
					return err
				}
				fieldValue.Set(x.Elem())
			}
		} else if col.SQLType.IsBlob() {
			if k == reflect.Uint8 {
				fieldValue.Set(reflect.ValueOf(v))
			} else {
				x := reflect.New(fieldType)
				if len(data) > 0 {
					err := json.DefaultJSONHandler.Unmarshal(data, x.Interface())
					if err != nil {
						session.engine.logger.Errorf("%v", err)
						return err
					}
					fieldValue.Set(x.Elem())
				}
			}
		} else {
			return ErrUnSupportedType
		}
	case reflect.String:
		fieldValue.SetString(string(data))
	case reflect.Bool:
		v, err := asBool(data)
		if err != nil {
			return fmt.Errorf("arg %v as bool: %s", key, err.Error())
		}
		fieldValue.Set(reflect.ValueOf(v))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sdata := string(data)
		var x int64
		var err error
		// for mysql, when use bit, it returned \x01
		if col.SQLType.Name == schemas.Bit &&
			session.engine.dialect.URI().DBType == schemas.MYSQL { // !nashtsai! TODO dialect needs to provide conversion interface API
			if len(data) == 1 {
				x = int64(data[0])
			} else {
				x = 0
			}
		} else if strings.HasPrefix(sdata, "0x") {
			x, err = strconv.ParseInt(sdata, 16, 64)
		} else if strings.HasPrefix(sdata, "0") {
			x, err = strconv.ParseInt(sdata, 8, 64)
		} else if strings.EqualFold(sdata, "true") {
			x = 1
		} else if strings.EqualFold(sdata, "false") {
			x = 0
		} else {
			x, err = strconv.ParseInt(sdata, 10, 64)
		}
		if err != nil {
			return fmt.Errorf("arg %v as int: %s", key, err.Error())
		}
		fieldValue.SetInt(x)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return fmt.Errorf("arg %v as float64: %s", key, err.Error())
		}
		fieldValue.SetFloat(x)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		x, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("arg %v as int: %s", key, err.Error())
		}
		fieldValue.SetUint(x)
	//Currently only support Time type
	case reflect.Struct:
		// !<winxxp>! 增加支持sql.Scanner接口的结构，如sql.NullString
		if nulVal, ok := fieldValue.Addr().Interface().(sql.Scanner); ok {
			if err := nulVal.Scan(data); err != nil {
				return fmt.Errorf("sql.Scan(%v) failed: %s ", data, err.Error())
			}
		} else {
			if fieldType.ConvertibleTo(schemas.TimeType) {
				x, err := session.byte2Time(col, data)
				if err != nil {
					return err
				}
				v = x
				fieldValue.Set(reflect.ValueOf(v).Convert(fieldType))
			} else if session.statement.UseCascade {
				table, err := session.engine.tagParser.ParseWithCache(*fieldValue)
				if err != nil {
					return err
				}

				// TODO: current only support 1 primary key
				if len(table.PrimaryKeys) > 1 {
					return errors.New("unsupported composited primary key cascade")
				}

				var pk = make(schemas.PK, len(table.PrimaryKeys))
				rawValueType := table.ColumnType(table.PKColumns()[0].FieldName)
				pk[0], err = str2PK(string(data), rawValueType)
				if err != nil {
					return err
				}

				if !pk.IsZero() {
					// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
					// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
					// property to be fetched lazily
					structInter := reflect.New(fieldValue.Type())
					has, err := session.ID(pk).NoCascade().get(structInter.Interface())
					if err != nil {
						return err
					}
					if has {
						v = structInter.Elem().Interface()
						fieldValue.Set(reflect.ValueOf(v))
					} else {
						return errors.New("cascade obj is not exist")
					}
				}
			}
		}
	case reflect.Ptr:
		// !nashtsai! TODO merge duplicated codes above
		//typeStr := fieldType.String()
		switch fieldType.Elem().Kind() {
		// case "*string":
		case schemas.StringType.Kind():
			x := string(data)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*bool":
		case schemas.BoolType.Kind():
			d := string(data)
			v, err := strconv.ParseBool(d)
			if err != nil {
				return fmt.Errorf("arg %v as bool: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&v).Convert(fieldType))
		// case "*complex64":
		case schemas.Complex64Type.Kind():
			var x complex64
			if len(data) > 0 {
				err := json.DefaultJSONHandler.Unmarshal(data, &x)
				if err != nil {
					session.engine.logger.Errorf("%v", err)
					return err
				}
				fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
			}
		// case "*complex128":
		case schemas.Complex128Type.Kind():
			var x complex128
			if len(data) > 0 {
				err := json.DefaultJSONHandler.Unmarshal(data, &x)
				if err != nil {
					session.engine.logger.Errorf("%v", err)
					return err
				}
				fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
			}
		// case "*float64":
		case schemas.Float64Type.Kind():
			x, err := strconv.ParseFloat(string(data), 64)
			if err != nil {
				return fmt.Errorf("arg %v as float64: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*float32":
		case schemas.Float32Type.Kind():
			var x float32
			x1, err := strconv.ParseFloat(string(data), 32)
			if err != nil {
				return fmt.Errorf("arg %v as float32: %s", key, err.Error())
			}
			x = float32(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint64":
		case schemas.Uint64Type.Kind():
			var x uint64
			x, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint":
		case schemas.UintType.Kind():
			var x uint
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint32":
		case schemas.Uint32Type.Kind():
			var x uint32
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint32(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint8":
		case schemas.Uint8Type.Kind():
			var x uint8
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint8(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*uint16":
		case schemas.Uint16Type.Kind():
			var x uint16
			x1, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			x = uint16(x1)
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int64":
		case schemas.Int64Type.Kind():
			sdata := string(data)
			var x int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == schemas.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int64(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x, err = strconv.ParseInt(sdata, 16, 64)
			} else if strings.HasPrefix(sdata, "0") {
				x, err = strconv.ParseInt(sdata, 8, 64)
			} else {
				x, err = strconv.ParseInt(sdata, 10, 64)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int":
		case schemas.IntType.Kind():
			sdata := string(data)
			var x int
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == schemas.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int32":
		case schemas.Int32Type.Kind():
			sdata := string(data)
			var x int32
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == schemas.Bit &&
				session.engine.dialect.URI().DBType == schemas.MYSQL {
				if len(data) == 1 {
					x = int32(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int32(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int32(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int32(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int8":
		case schemas.Int8Type.Kind():
			sdata := string(data)
			var x int8
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == schemas.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int8(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int8(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int8(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int8(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*int16":
		case schemas.Int16Type.Kind():
			sdata := string(data)
			var x int16
			var x1 int64
			var err error
			// for mysql, when use bit, it returned \x01
			if col.SQLType.Name == schemas.Bit &&
				strings.Contains(session.engine.DriverName(), "mysql") {
				if len(data) == 1 {
					x = int16(data[0])
				} else {
					x = 0
				}
			} else if strings.HasPrefix(sdata, "0x") {
				x1, err = strconv.ParseInt(sdata, 16, 64)
				x = int16(x1)
			} else if strings.HasPrefix(sdata, "0") {
				x1, err = strconv.ParseInt(sdata, 8, 64)
				x = int16(x1)
			} else {
				x1, err = strconv.ParseInt(sdata, 10, 64)
				x = int16(x1)
			}
			if err != nil {
				return fmt.Errorf("arg %v as int: %s", key, err.Error())
			}
			fieldValue.Set(reflect.ValueOf(&x).Convert(fieldType))
		// case "*SomeStruct":
		case reflect.Struct:
			switch fieldType {
			// case "*.time.Time":
			case schemas.PtrTimeType:
				x, err := session.byte2Time(col, data)
				if err != nil {
					return err
				}
				v = x
				fieldValue.Set(reflect.ValueOf(&x))
			default:
				if session.statement.UseCascade {
					structInter := reflect.New(fieldType.Elem())
					table, err := session.engine.tagParser.ParseWithCache(structInter.Elem())
					if err != nil {
						return err
					}

					if len(table.PrimaryKeys) > 1 {
						return errors.New("unsupported composited primary key cascade")
					}

					var pk = make(schemas.PK, len(table.PrimaryKeys))
					rawValueType := table.ColumnType(table.PKColumns()[0].FieldName)
					pk[0], err = str2PK(string(data), rawValueType)
					if err != nil {
						return err
					}

					if !pk.IsZero() {
						// !nashtsai! TODO for hasOne relationship, it's preferred to use join query for eager fetch
						// however, also need to consider adding a 'lazy' attribute to xorm tag which allow hasOne
						// property to be fetched lazily
						has, err := session.ID(pk).NoCascade().get(structInter.Interface())
						if err != nil {
							return err
						}
						if has {
							v = structInter.Interface()
							fieldValue.Set(reflect.ValueOf(v))
						} else {
							return errors.New("cascade obj is not exist")
						}
					}
				} else {
					return fmt.Errorf("unsupported struct type in Scan: %s", fieldValue.Type().String())
				}
			}
		default:
			return fmt.Errorf("unsupported type in Scan: %s", fieldValue.Type().String())
		}
	default:
		return fmt.Errorf("unsupported type in Scan: %s", fieldValue.Type().String())
	}

	return nil
}

// convert a field value of a struct to interface for put into db
func (session *Session) value2Interface(col *schemas.Column, fieldValue reflect.Value) (interface{}, error) {
	if fieldValue.CanAddr() {
		if fieldConvert, ok := fieldValue.Addr().Interface().(convert.Conversion); ok {
			data, err := fieldConvert.ToDB()
			if err != nil {
				return 0, err
			}
			if col.SQLType.IsBlob() {
				return data, nil
			}
			return string(data), nil
		}
	}

	if fieldConvert, ok := fieldValue.Interface().(convert.Conversion); ok {
		data, err := fieldConvert.ToDB()
		if err != nil {
			return 0, err
		}
		if col.SQLType.IsBlob() {
			return data, nil
		}
		return string(data), nil
	}

	fieldType := fieldValue.Type()
	k := fieldType.Kind()
	if k == reflect.Ptr {
		if fieldValue.IsNil() {
			return nil, nil
		} else if !fieldValue.IsValid() {
			session.engine.logger.Warnf("the field [%s] is invalid", col.FieldName)
			return nil, nil
		} else {
			// !nashtsai! deference pointer type to instance type
			fieldValue = fieldValue.Elem()
			fieldType = fieldValue.Type()
			k = fieldType.Kind()
		}
	}

	switch k {
	case reflect.Bool:
		return fieldValue.Bool(), nil
	case reflect.String:
		return fieldValue.String(), nil
	case reflect.Struct:
		if fieldType.ConvertibleTo(schemas.TimeType) {
			t := fieldValue.Convert(schemas.TimeType).Interface().(time.Time)
			tf := dialects.FormatColumnTime(session.engine.dialect, session.engine.DatabaseTZ, col, t)
			return tf, nil
		} else if fieldType.ConvertibleTo(nullFloatType) {
			t := fieldValue.Convert(nullFloatType).Interface().(sql.NullFloat64)
			if !t.Valid {
				return nil, nil
			}
			return t.Float64, nil
		}

		if !col.SQLType.IsJson() {
			// !<winxxp>! 增加支持driver.Valuer接口的结构，如sql.NullString
			if v, ok := fieldValue.Interface().(driver.Valuer); ok {
				return v.Value()
			}

			fieldTable, err := session.engine.tagParser.ParseWithCache(fieldValue)
			if err != nil {
				return nil, err
			}
			if len(fieldTable.PrimaryKeys) == 1 {
				pkField := reflect.Indirect(fieldValue).FieldByName(fieldTable.PKColumns()[0].FieldName)
				return pkField.Interface(), nil
			}
			return 0, fmt.Errorf("no primary key for col %v", col.Name)
		}

		if col.SQLType.IsText() {
			bytes, err := json.DefaultJSONHandler.Marshal(fieldValue.Interface())
			if err != nil {
				session.engine.logger.Errorf("%v", err)
				return 0, err
			}
			return string(bytes), nil
		} else if col.SQLType.IsBlob() {
			bytes, err := json.DefaultJSONHandler.Marshal(fieldValue.Interface())
			if err != nil {
				session.engine.logger.Errorf("%v", err)
				return 0, err
			}
			return bytes, nil
		}
		return nil, fmt.Errorf("Unsupported type %v", fieldValue.Type())
	case reflect.Complex64, reflect.Complex128:
		bytes, err := json.DefaultJSONHandler.Marshal(fieldValue.Interface())
		if err != nil {
			session.engine.logger.Errorf("%v", err)
			return 0, err
		}
		return string(bytes), nil
	case reflect.Array, reflect.Slice, reflect.Map:
		if !fieldValue.IsValid() {
			return fieldValue.Interface(), nil
		}

		if col.SQLType.IsText() {
			bytes, err := json.DefaultJSONHandler.Marshal(fieldValue.Interface())
			if err != nil {
				session.engine.logger.Errorf("%v", err)
				return 0, err
			}
			return string(bytes), nil
		} else if col.SQLType.IsBlob() {
			var bytes []byte
			var err error
			if (k == reflect.Slice) &&
				(fieldValue.Type().Elem().Kind() == reflect.Uint8) {
				bytes = fieldValue.Bytes()
			} else {
				bytes, err = json.DefaultJSONHandler.Marshal(fieldValue.Interface())
				if err != nil {
					session.engine.logger.Errorf("%v", err)
					return 0, err
				}
			}
			return bytes, nil
		}
		return nil, ErrUnSupportedType
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return int64(fieldValue.Uint()), nil
	default:
		return fieldValue.Interface(), nil
	}
}
