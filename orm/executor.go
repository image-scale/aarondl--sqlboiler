package orm

import (
	"context"
	"database/sql"
)

type Executor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type ContextExecutor interface {
	Executor
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Transactor interface {
	Executor
	Commit() error
	Rollback() error
}

type ContextTransactor interface {
	ContextExecutor
	Commit() error
	Rollback() error
}

type Beginner interface {
	Begin() (*sql.Tx, error)
}

type ContextBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func Begin() (Transactor, error) {
	b, ok := currentDB.(Beginner)
	if !ok {
		panic("orm: global database does not support Begin")
	}
	tx, err := b.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	b, ok := currentDB.(ContextBeginner)
	if !ok {
		panic("orm: global database does not support BeginTx")
	}
	return b.BeginTx(ctx, opts)
}
