package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/elgris/sqrl"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/ketan-10/sql2gql/examples/pokemon/internal/context_manager"
)

type IDb interface {
	Select(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error
	Get(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error
	Exec(ctx context.Context, sqlizer sqrl.Sqlizer) (sql.Result, error)
	BeginTxx(ctx context.Context) (*sqlx.Tx, error)
}

type DBOptions struct {
}

type DB struct {
	*DBOptions
	DB *sqlx.DB
}

var NewDB = wire.NewSet(
	wire.Struct(new(DBOptions), "*"),
	OpenConnection,
	wire.Bind(new(IDb), new(*DB)),
)

func OpenConnection(ctx context.Context, options *DBOptions) *DB {

	connection, err := context_manager.GetConnectionContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sqlxDB, err := sqlx.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{DBOptions: options, DB: sqlxDB}
}

func (db *DB) Get(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}

	if tx := context_manager.GetTransactionContext(ctx); tx != nil {
		err = tx.GetContext(ctx, dest, query, args...)
	} else {
		err = db.DB.GetContext(ctx, dest, query, args...)
	}

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Select(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error {
	query, args, err := sqlizer.ToSql()
	if err != nil {
		return err
	}

	if tx := context_manager.GetTransactionContext(ctx); tx != nil {
		err = tx.SelectContext(ctx, dest, query, args...)
	} else {
		err = db.DB.SelectContext(ctx, dest, query, args...)
	}

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Exec(ctx context.Context, sqlizer sqrl.Sqlizer) (sql.Result, error) {

	// // If prod don't attempt exec, just write to file
	// if getProdContext(ctx) {
	// 	err := db.FileGen.Write(sqlizer)
	// 	return nil, err
	// }

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, err
	}

	var res sql.Result
	if tx := context_manager.GetTransactionContext(ctx); tx != nil {
		res, err = tx.ExecContext(ctx, query, args...)
	} else {
		res, err = db.DB.ExecContext(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}

	// // write to file if query was successful
	// err = db.FileGen.Write(sqlizer)

	return res, err
}

func (db *DB) BeginTxx(ctx context.Context) (*sqlx.Tx, error) {
	return db.DB.BeginTxx(ctx, nil)
}

// Transaction

func WrapInTransaction(ctx context.Context, db IDb, f func(ctx context.Context) error) error {

	// if transaction already exists
	tx := context_manager.GetTransactionContext(ctx)
	if tx != nil {
		return f(ctx)
	}

	tx, err := db.BeginTxx(ctx)
	if err != nil {
		return err
	}

	// handle traction
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Rollback due to panic")
			panic(r)
		}
		tx.Commit()

	}()

	newContext := context_manager.WithTransaction(ctx, tx)
	return f(newContext)
}
