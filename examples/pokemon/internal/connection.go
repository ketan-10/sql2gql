package internal

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/elgris/sqrl"
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/ketan-10/sql2gql/examples/pokemon/internal/context_manager"
)

type IDb interface {
	Select(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error
	Get(ctx context.Context, dest interface{}, sqlizer sqrl.Sqlizer) error
	Exec(ctx context.Context, sqlizer sqrl.Sqlizer) (sql.Result, error)
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
	sqlxDB.SetMaxOpenConns(50)      // no more than 50 open at once
	sqlxDB.SetMaxIdleConns(25)      // keep up to 25 idle
	sqlxDB.SetConnMaxLifetime(30 * time.Minute)
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

	err = db.DB.GetContext(ctx, dest, query, args...)

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

	err = db.DB.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Exec(ctx context.Context, sqlizer sqrl.Sqlizer) (sql.Result, error) {

	query, args, err := sqlizer.ToSql()
	if err != nil {
		return nil, err
	}

	var res sql.Result
	res, err = db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return res, err
}
