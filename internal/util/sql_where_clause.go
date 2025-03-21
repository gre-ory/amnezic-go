package util

import (
	"fmt"
	"strings"
)

type SqlWhereClause interface {
	IsEmpty() bool
	WithCondition(condition string, args ...any) SqlWhereClause
	WithRandomOrder() SqlWhereClause
	WithLimit(limit int) SqlWhereClause
	Generate(placeHolder int) (string, []any)
}

func NewSqlWhereClause() SqlWhereClause {
	return &sqlWhereClause{
		conditions: make([]string, 0),
		args:       make([]any, 0),
	}
}

func NewSqlCondition(condition string, args ...any) SqlWhereClause {
	wc := &sqlWhereClause{
		conditions:  make([]string, 0),
		args:        make([]any, 0),
		placeHolder: 0,
	}
	return wc.WithCondition(condition, args...)
}

var NoWhereClause = NewSqlWhereClause()

type sqlWhereClause struct {
	conditions  []string
	args        []any
	placeHolder int
	orderBy     string
	limit       int
}

func (wc *sqlWhereClause) IsEmpty() bool {
	return wc == nil || len(wc.conditions) == 0
}

const NumberedPlaceHolder = "$_"

func (wc *sqlWhereClause) WithCondition(condition string, args ...any) SqlWhereClause {
	if wc != nil {
		count := strings.Count(condition, NumberedPlaceHolder)
		if count != len(args) {
			panic(fmt.Sprintf("mismatch between number of placeholders '%s' (%d) and number of arguments (%d)! ( condition: '%s', args: %#v )", NumberedPlaceHolder, count, len(args), condition, args))
		}
		wc.conditions = append(wc.conditions, condition)
		wc.args = append(wc.args, args...)
	}
	return wc
}

func (wc *sqlWhereClause) WithRandomOrder() SqlWhereClause {
	wc.orderBy = "RANDOM()"
	return wc
}

func (wc *sqlWhereClause) WithLimit(limit int) SqlWhereClause {
	wc.limit = limit
	return wc
}

func (wc *sqlWhereClause) Generate(placeHolder int) (string, []any) {
	var whereClause string
	if !wc.IsEmpty() {
		whereClause = "WHERE " + strings.Join(wc.conditions, " AND ")
		for range wc.args {
			placeHolder++
			whereClause = strings.Replace(whereClause, NumberedPlaceHolder, fmt.Sprintf("$%d", placeHolder), 1)
		}
	}
	if wc.orderBy != "" {
		if whereClause != "" {
			whereClause += " "
		}
		whereClause += fmt.Sprintf("ORDER BY %s", wc.orderBy)
	}
	if wc.limit != 0 {
		if whereClause != "" {
			whereClause += " "
		}
		whereClause += fmt.Sprintf("LIMIT %d", wc.limit)
	}
	return whereClause, wc.args
}
