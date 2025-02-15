package storage

import (
	"context"
	"fmt"

	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/jackc/pgx/v5"
)

var getBalanceQuery = `SELECT current, withdrawn FROM balances WHERE username=$1`

type BalanceRepository struct {
	db *DB
}

func NewBalanceRepository(db *DB) *BalanceRepository {
	return &BalanceRepository{
		db,
	}
}

func (b *BalanceRepository) Balance(ctx context.Context, user string) (*model.Balance, error) {

	row, err := b.db.pgPool.Query(ctx, getBalanceQuery, user)
	if err != nil {
		return nil, err
	}

	ex, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.Balance])
	fmt.Println(err)

	if err != nil {
		return nil, model.ErrInternal
	}

	return &ex, nil
}
