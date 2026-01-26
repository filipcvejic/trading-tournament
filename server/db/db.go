package db

import (
	"context"
	"github.com/filipcvejic/trading_tournament/db/sqlc"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DB struct {
	Pool  *pgxpool.Pool
	Query *sqlc.Queries
}

func NewDatabase(connStr string) *DB {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("database connection failed:", err)
	}

	query := sqlc.New(pool)

	database := DB{
		Pool:  pool,
		Query: query,
	}

	return &database
}

func (d *DB) WithTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	tx, err := d.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // safe even if committed

	q := d.Query.WithTx(tx)

	if err := fn(q); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
