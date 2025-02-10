package storage

import "context"

// var getBalanceQuery = `SELECT current, withdraw FROM balance WHERE username=$1`

type BalanceRepository struct {
	db *DB
}

func NewBalanceRepository(db *DB) *BalanceRepository {
	return &BalanceRepository{
		db,
	}
}

func (b *BalanceRepository) Balance(ctx context.Context) error {
	// row, err := b.db.pgPool.Query(ctx, getBalanceQuery)
	// if err != nil {
	// 	return err
	// }
	return nil
}
