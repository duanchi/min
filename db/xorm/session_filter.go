package xorm

import (
	"reflect"
	"strings"

	"github.com/duanchi/min/v2/server/httpserver/context"
	"github.com/duanchi/min/v2/util"
)

const FILTER_MODE_EQUAL = "="
const FILTER_MODE_LESS = "<"
const FILTER_MODE_LESS_OR_EQUAL = "<="
const FILTER_MODE_GREATER_OR_EQUAL = ">="
const FILTER_MODE_GREATER = ">"
const FILTER_MODE_NOT_EQUAL = "!="
const FILTER_MODE_IN = "IN"
const FILTER_MODE_LIKE = "LIKE"

type FilterOptions struct {
	QueryKey   string
	Column     string
	FilterMode string
}

func (session *Session) FilterWithOption(filterOptions FilterOptions, value any, defaults ...any) *Session {
	defaultValue := ""
	if len(defaults) > 0 {
		defaultValue = util.ToString(defaults[0])
	}
	switch value.(type) {
	case *context.Context:
		v := value.(*context.Context).Query(filterOptions.QueryKey, defaultValue)
		if v == "" {
			return session
		}
		return session.applyFilter(filterOptions, v)
	default:
		return session.applyFilter(filterOptions, value)
	}
}

func (session *Session) Filter(column string, value any, defaults ...any) *Session {
	return session.FilterWithOption(FilterOptions{
		QueryKey:   column,
		Column:     column,
		FilterMode: FILTER_MODE_EQUAL,
	}, value, defaults...)
}

func (session *Session) FilterWithContext(ctx *context.Context, column ...string) *Session {
	for _, c := range column {
		session = session.FilterWithOption(FilterOptions{
			QueryKey:   c,
			Column:     c,
			FilterMode: FILTER_MODE_EQUAL,
		}, ctx)
	}
	return session
}

func (session *Session) FilterWithContextDefaults(ctx *context.Context, columnWithDefaults ...any) *Session {
	t := reflect.TypeOf(columnWithDefaults[0])

	if t.Kind() == reflect.Map {
		for k, v := range columnWithDefaults[0].(map[string]any) {
			session = session.FilterWithOption(FilterOptions{
				QueryKey:   k,
				Column:     k,
				FilterMode: FILTER_MODE_EQUAL,
			}, ctx, v)
		}
	} else {
		if len(columnWithDefaults) == 1 {
			session = session.FilterWithOption(FilterOptions{
				QueryKey:   columnWithDefaults[0].(string),
				Column:     columnWithDefaults[0].(string),
				FilterMode: FILTER_MODE_EQUAL,
			}, ctx)
		} else if len(columnWithDefaults) > 1 {
			session = session.FilterWithOption(FilterOptions{
				QueryKey:   columnWithDefaults[0].(string),
				Column:     columnWithDefaults[0].(string),
				FilterMode: FILTER_MODE_EQUAL,
			}, ctx, columnWithDefaults[1])
		}
	}
	return session
}

func (session *Session) FilterWithRange(rangeKeys [2]string, column string, value any) *Session {
	startFilter := FilterOptions{
		QueryKey:   rangeKeys[0],
		Column:     column,
		FilterMode: FILTER_MODE_GREATER_OR_EQUAL,
	}
	endFilter := FilterOptions{
		QueryKey:   rangeKeys[1],
		Column:     column,
		FilterMode: FILTER_MODE_LESS,
	}
	switch value.(type) {
	case *context.Context:
		startValue := value.(*context.Context).Query(startFilter.QueryKey)
		endValue := value.(*context.Context).Query(endFilter.QueryKey)
		if startValue != "" {
			session = session.FilterWithOption(startFilter, startValue)
		}
		if endValue != "" {
			session = session.FilterWithOption(endFilter, endValue)
		}
	default:
		t := reflect.TypeOf(value)
		v := reflect.ValueOf(value)
		if t != nil && (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) {
			if v.Len() >= 2 {
				session = session.FilterWithOption(endFilter, v.Index(1).Interface())
			}
			if v.Len() >= 1 {
				session = session.FilterWithOption(startFilter, v.Index(0).Interface())
			}
		}
	}
	return session
}

func (session *Session) applyFilter(filterOptions FilterOptions, value any) *Session {
	switch filterOptions.FilterMode {
	case FILTER_MODE_GREATER, FILTER_MODE_GREATER_OR_EQUAL, FILTER_MODE_LESS, FILTER_MODE_LESS_OR_EQUAL, FILTER_MODE_EQUAL, FILTER_MODE_NOT_EQUAL:
		session.And("`"+filterOptions.Column+"` "+filterOptions.FilterMode+" ?", value)
	case FILTER_MODE_LIKE:
		session.And("`"+filterOptions.Column+"` LIKE ?", "%"+util.ToString(value)+"%")
	case FILTER_MODE_IN:
		switch value.(type) {
		case string:
			values := strings.Split(value.(string), ",")
			for i := range values {
				values[i] = strings.TrimSpace(values[i])
			}
			session.In(filterOptions.Column, values)
		default:
			t := reflect.TypeOf(value)
			if t != nil && (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) {
				session.In(filterOptions.Column, value)
			}
		}
	}
	return session
}
