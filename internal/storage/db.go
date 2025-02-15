package storage

import (
	"context"
	"log"
	"sync"

	"github.com/JohnRobertFord/go-market/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pgPool *pgxpool.Pool
	cfg    *config.Config
}

var (
	pgInstance           *DB
	pgOnce               sync.Once
	createUserTableQuery = `CREATE TABLE IF NOT EXISTS users(
		"id"		 int generated always as identity,
		"username"	 varchar(20) UNIQUE NOT NULL,
		"hash" 		 varchar(100) NOT NULL,
		"created_at" TIMESTAMP NOT NULL
		);`
	createOrderTableQuery = `CREATE TABLE IF NOT EXISTS orders(
		"id"		 int generated always as identity,
		"username"	 varchar(20) NOT NULL,
		"order_id" 	 varchar(50) NOT NULL,
		"status"	 varchar(20) NOT NULL,
		"accrual"	 NUMERIC,
		"created_at" TIMESTAMP NOT NULL,
		"updated_at" TIMESTAMP NOT NULL
		);`
	createBalanceTableQuery = `CREATE TABLE IF NOT EXISTS balances(
		"id"		 int generated always as identity,
		"username"	 varchar(20) UNIQUE NOT NULL,
		"current"	 NUMERIC NOT NULL,
		"withdrawn"	 NUMERIC NOT NULL
		);`
	createWithdrawalTableQuery = `CREATE TABLE IF NOT EXISTS withdrawals(
		"id"		 int generated always as identity,
		"username"	 varchar(20) NOT NULL,
		"order_id"	 varchar(50) NOT NULL,
		"sum"		 NUMERIC NOT NULL,
		"processed_at" TIMESTAMP NOT NULL
		);`
)

// number, sum, processed_at
func (p *DB) Close() {
	p.pgPool.Close()
}

func NewStorage(ctx context.Context, c *config.Config) (*DB, error) {

	pgOnce.Do(func() {
		dbPool, err := pgxpool.New(ctx, c.DatabaseURI)
		if err != nil {
			log.Printf("unable to create connection pool: %s", err)
			return
		}
		pgInstance = &DB{dbPool, c}
	})

	err := pgInstance.pgPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	_, err = pgInstance.pgPool.Exec(ctx, createUserTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = pgInstance.pgPool.Exec(ctx, createOrderTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = pgInstance.pgPool.Exec(ctx, createBalanceTableQuery)
	if err != nil {
		return nil, err
	}
	_, err = pgInstance.pgPool.Exec(ctx, createWithdrawalTableQuery)
	if err != nil {
		return nil, err
	}
	return pgInstance, nil
}
