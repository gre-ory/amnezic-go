package util

import (
	"context"
	"database/sql"
)

func SqlTransaction(ctx context.Context, db *sql.DB, onTransaction func(tx *sql.Tx)) (err error) {

	var tx *sql.Tx
	tx, err = db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if panicErr, ok := r.(error); ok {
				tx.Rollback()
				err = panicErr
			} else {
				panic(r)
			}
		}
	}()

	onTransaction(tx)

	err = tx.Commit()
	return
}
