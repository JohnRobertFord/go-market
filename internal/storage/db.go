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
	pgInstance       *DB
	pgOnce           sync.Once
	createTableQuery = `CREATE TABLE IF NOT EXISTS users(
		"id"		 int generated always as identity,
		"username"	 varchar(20) UNIQUE NOT NULL,
		"hash" 		 varchar(100) NOT NULL,
		"created_at" DATE NOT NULL
		);`
)

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

	_, err = pgInstance.pgPool.Exec(ctx, createTableQuery)

	return pgInstance, err
}
