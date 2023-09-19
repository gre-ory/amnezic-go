package util_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/gre-ory/amnezic-go/internal/util"
)

func TestSqlTable(t *testing.T) {

	//
	// sql table
	//

	type TestRow struct {
		Id    int64  `sql:"id,auto-generated"`
		Name  string `sql:"name"`
		Value int64  `sql:"value,read-only"`
	}

	ErrTestRowNotFound := errors.New("test row not found")

	//
	// context
	//

	ctx := context.Background()

	//
	// logger
	//

	config := zap.NewDevelopmentConfig()
	config.Development = false
	logger, _ := config.Build()

	//
	// db mock
	//

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Info("[DEBUG] failed to mock db", zap.Error(err))
		t.Fatalf("unable to mock db: %s", err)
	}
	defer db.Close()

	//
	// insert row
	//

	t.Run("insert-row", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("INSERT INTO test (name,value) VALUES ($1,$2) RETURNING id,name,value").
			ExpectQuery().
			WithArgs("my-name", 99).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "value"}).
					AddRow(1002, "my-name", 99).
					AddRow(9999, "bad", 0),
			)
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)

		var gotRow *TestRow
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotRow = table.InsertRow(ctx, tx, &TestRow{
				Id:    1001,
				Name:  "my-name",
				Value: 99,
			})
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, &TestRow{
			Id:    1002,
			Name:  "my-name",
			Value: 99,
		}, gotRow, "wrong row")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// insert none
	//

	t.Run("insert-none", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("INSERT INTO test (name,value) VALUES ($1,$2) RETURNING id,name,value").
			ExpectQuery().
			WithArgs("my-name", 99).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "value"}),
			)
		mock.ExpectRollback()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)

		var gotRow *TestRow
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotRow = table.InsertRow(ctx, tx, &TestRow{
				Id:    1001,
				Name:  "my-name",
				Value: 99,
			})
		})

		require.Equal(t, ErrTestRowNotFound, gotErr, "wrong error")
		require.Nil(t, gotRow, "wrong row")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// select row
	//

	t.Run("select-row", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT id,name,value FROM test where id = $1 LIMIT 1").
			ExpectQuery().
			WithArgs(1002).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "value"}).
					AddRow(1002, "my-name", 99).
					AddRow(9999, "bad", 0),
			)
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)
		var gotRow *TestRow
		var gotSelectErr error
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotRow, gotSelectErr = table.SelectRow(ctx, tx, "where id = $1", 1002)
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, nil, gotSelectErr, "wrong select error")
		require.Equal(t, &TestRow{
			Id:    1002,
			Name:  "my-name",
			Value: 99,
		}, gotRow, "wrong row")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// select none
	//

	t.Run("select-none", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT id,name,value FROM test where id = $1 LIMIT 1").
			ExpectQuery().
			WithArgs(1002).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "value"}),
			)
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)
		var gotRow *TestRow
		var gotSelectErr error
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotRow, gotSelectErr = table.SelectRow(ctx, tx, "where id = $1", 1002)
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, ErrTestRowNotFound, gotSelectErr, "wrong select error")
		require.Nil(t, gotRow, "wrong row")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// exists row
	//

	t.Run("exists-row", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT EXISTS( SELECT 1 FROM test where id = $1 )").
			ExpectQuery().
			WithArgs(1002).
			WillReturnRows(
				sqlmock.NewRows([]string{"exists"}).
					AddRow(true),
			)
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)
		var gotExists bool
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotExists = table.ExistsRow(ctx, tx, "where id = $1", 1002)
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, true, gotExists, "wrong exists")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// exists none
	//

	t.Run("exists-row", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("SELECT EXISTS( SELECT 1 FROM test where id = $1 )").
			ExpectQuery().
			WithArgs(1002).
			WillReturnRows(
				sqlmock.NewRows([]string{"exists"}).
					AddRow(false),
			)
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)
		var gotExists bool
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotExists = table.ExistsRow(ctx, tx, "where id = $1", 1002)
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, false, gotExists, "wrong exists")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// update row
	//

	t.Run("update-row", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("UPDATE test SET name=$2 where id = $1 LIMIT 1 RETURNING id,name,value").
			ExpectQuery().
			WithArgs(1001, "my-name").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "value"}).
					AddRow(1002, "my-name", 99).
					AddRow(9999, "bad", 0),
			)
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)

		var gotRow *TestRow
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotRow = table.UpdateRow(ctx, tx, &TestRow{
				Id:    1001,
				Name:  "my-name",
				Value: 99,
			}, "where id = $1", 1001)
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, &TestRow{
			Id:    1002,
			Name:  "my-name",
			Value: 99,
		}, gotRow, "wrong row")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// update none
	//

	t.Run("update-none", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("UPDATE test SET name=$2 where id = $1 LIMIT 1 RETURNING id,name,value").
			ExpectQuery().
			WithArgs(1001, "my-name").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "value"}),
			)
		mock.ExpectRollback()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)

		var gotRow *TestRow
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotRow = table.UpdateRow(ctx, tx, &TestRow{
				Id:    1001,
				Name:  "my-name",
				Value: 99,
			}, "where id = $1", 1001)
		})

		require.Equal(t, ErrTestRowNotFound, gotErr, "wrong error")
		require.Nil(t, gotRow, "wrong row")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// delete row
	//

	t.Run("delete-row", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("DELETE FROM test where id = $1 LIMIT 1").
			ExpectExec().
			WithArgs(1001).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)

		var gotDeleteErr error
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotDeleteErr = table.DeleteRow(ctx, tx, "where id = $1", 1001)
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, nil, gotDeleteErr, "wrong delete error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	//
	// delete none
	//

	t.Run("delete-none", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectPrepare("DELETE FROM test where id = $1 LIMIT 1").
			ExpectExec().
			WithArgs(1001).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectCommit()

		table := util.NewSqlTable[TestRow](logger, "test", ErrTestRowNotFound)

		var gotDeleteErr error
		gotErr := util.SqlTransaction(ctx, db, func(tx *sql.Tx) {
			gotDeleteErr = table.DeleteRow(ctx, tx, "where id = $1", 1001)
		})

		require.Equal(t, nil, gotErr, "wrong error")
		require.Equal(t, ErrTestRowNotFound, gotDeleteErr, "wrong delete error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
