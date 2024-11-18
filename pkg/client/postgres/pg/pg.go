package pgv1

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	pgClient "github.com/kms-qwe/platform_common/pkg/client/postgres"
	"github.com/kms-qwe/platform_common/pkg/client/postgres/prettier"
)

type key string

const (
	// TxKey key for transaction in context
	TxKey key = "tx"
)

// MakeContextTx make context with TxKey
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

type pg struct {
	dbc *pgxpool.Pool
}

// NewDB creates new pgClietn.DB
func NewDB(dbc *pgxpool.Pool) pgClient.DB {
	return &pg{
		dbc: dbc,
	}
}

// ScanOneContext scan one row result of QueryContext with pgxscan
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q pgClient.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext scan result of QueryContext with pgxscan
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q pgClient.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// ExecContext wrapper of pgpool Exec
func (p *pg) ExecContext(ctx context.Context, q pgClient.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext wrapper of pgpool Query
func (p *pg) QueryContext(ctx context.Context, q pgClient.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext wrapper of pgpool QueryRow
func (p *pg) QueryRowContext(ctx context.Context, q pgClient.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// Ping pings db
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close db connection
func (p *pg) Close() {
	p.dbc.Close()
}

// BeginTx begins transactions
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

func logQuery(ctx context.Context, q pgClient.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
