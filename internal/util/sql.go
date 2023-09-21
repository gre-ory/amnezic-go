package util

import (
	"context"
	"database/sql"
)

// //////////////////////////////////////////////////
// query

func SqlQuery(ctx context.Context, tx *sql.Tx, query string, args ...any) *sql.Rows {

	stmt, err := tx.Prepare(query) // to avoid SQL injection
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		panic(err)
	}
	return rows
}

// //////////////////////////////////////////////////
// sql scan

func SqlScan(rows *sql.Rows, scanFn func(rows *sql.Rows)) {
	defer rows.Close()
	for rows.Next() {
		scanFn(rows)
	}
}

// //////////////////////////////////////////////////
// sql decode

func SqlDecode[Row any](rows *sql.Rows, decodeFn func(rows *sql.Rows) *Row) []*Row {
	defer rows.Close()

	result := make([]*Row, 0)
	for rows.Next() {
		result = append(result, decodeFn(rows))
	}
	return result
}
