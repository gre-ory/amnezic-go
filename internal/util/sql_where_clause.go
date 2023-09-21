package util

import (
	"fmt"
	"strings"
)

type SqlWhereClause interface {
	IsEmpty() bool
	WithCondition(condition string, args ...any) SqlWhereClause
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

type sqlWhereClause struct {
	conditions  []string
	args        []any
	placeHolder int
}

func (wc *sqlWhereClause) IsEmpty() bool {
	return wc == nil || len(wc.conditions) == 0
}

func (wc *sqlWhereClause) WithCondition(condition string, args ...any) SqlWhereClause {
	if wc != nil {
		count := strings.Count(condition, "%s")
		if count != len(args) {
			panic(fmt.Sprintf("mismatch between number of placeholders '%%s' (%d) and number of arg (%d)! ( condition: '%s', args: %#v )", count, len(args), condition, args))
		}
		wc.conditions = append(wc.conditions, condition)
		wc.args = append(wc.args, args...)
	}
	return wc
}

func (wc *sqlWhereClause) Generate(placeHolder int) (string, []any) {
	if wc.IsEmpty() {
		return "", wc.args
	}
	whereClause := "WHERE " + strings.Join(wc.conditions, " AND ")
	placeHolders := make([]any, 0, len(wc.args))
	for range wc.args {
		placeHolder++
		placeHolders = append(placeHolders, fmt.Sprintf("$%d", placeHolder))
	}
	whereClause = fmt.Sprintf(whereClause, placeHolders...)
	return whereClause, wc.args
}
