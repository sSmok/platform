package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/auth/internal/client/db"
	"github.com/sSmok/auth/internal/client/db/prettier"
)

type key string

// TxKey - ключ, который мы передаем в контекст для определения наличия транзакции
const TxKey key = "tx"

type pg struct {
	pool *pgxpool.Pool
}

// NewPG создает инстанс для работы с БД через переданное подключение pgxpool.Pool
func NewPG(pool *pgxpool.Pool) db.DB {
	return &pg{pool: pool}
}

func (pg *pg) Ping(ctx context.Context) error {
	err := pg.pool.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (pg *pg) ScanOneContext(ctx context.Context, dest interface{}, query db.Query, args ...interface{}) error {
	row, err := pg.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	err = pgxscan.ScanOne(dest, row)
	if err != nil {
		return err
	}

	return nil
}

func (pg *pg) ScanAllContext(ctx context.Context, dest interface{}, query db.Query, args ...interface{}) error {
	rows, err := pg.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	err = pgxscan.ScanAll(dest, rows)
	if err != nil {
		return err
	}

	return nil
}

func (pg *pg) ExecContext(ctx context.Context, query db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, query, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query.QueryRaw, args...)
	}

	return pg.pool.Exec(ctx, query.QueryRaw, args...)
}

func (pg *pg) QueryContext(ctx context.Context, query db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, query, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query.QueryRaw, args...)
	}

	return pg.pool.Query(ctx, query.QueryRaw, args...)
}

func (pg *pg) QueryRowContext(ctx context.Context, query db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, query, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query.QueryRaw, args...)
	}

	return pg.pool.QueryRow(ctx, query.QueryRaw, args...)
}

func (pg *pg) Close() {
	pg.pool.Close()
}

func (pg *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return pg.pool.BeginTx(ctx, txOptions)
}

// MakeContextTransaction - добавляет в контекст транзакцию
func MakeContextTransaction(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func logQuery(ctx context.Context, query db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(query.QueryRaw, prettier.DollarPlaceholder, args...)
	log.Println(ctx, fmt.Sprintf("method called: %s", query.Name), fmt.Sprintf("query: %s", prettyQuery))
}
