package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/jackc/pgx/v5"
)

var getUnprocessedOrders = "SELECT order_id, status, accrual, uploaded_at, username FROM orders WHERE status IN ('NEW','PROCESSING') ORDER BY uploaded_at"

type AccuralRepository struct {
	db *DB
}

func NewAccuralRepository(db *DB) *AccuralRepository {
	return &AccuralRepository{
		db,
	}
}

func (a *AccuralRepository) GetUnprocessedOrders(ctx context.Context) ([]model.Order, error) {
	row, err := a.db.pgPool.Query(ctx, getUnprocessedOrders)
	if err != nil {
		return nil, err
	}

	out, err := pgx.CollectRows(row, pgx.RowToStructByName[model.Order])
	fmt.Println(err)
	if err != nil {
		return nil, model.ErrInternal
	}

	return out, nil
}
func (a *AccuralRepository) UpdateOrder(ctx context.Context, order model.Order) error {
	tx, err := a.db.pgPool.Begin(ctx)
	if err != nil {
		return model.ErrInternalDB
	}
	defer func() {
		if errRollBack := tx.Rollback(ctx); errRollBack != nil {
			fmt.Printf("rollback error: %v", errRollBack)
		}
	}()

	// get balance
	row, err := a.db.pgPool.Query(ctx, getBalanceQuery, order.Username)
	if err != nil {
		return model.ErrInternalDB
	}
	balance, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.Balance])
	if err != nil {
		return model.ErrInternalDB
	}

	balance.Current += order.Accrual

	// update balance
	_, err = tx.Exec(ctx, updateBalancesQuery, order.Username, balance.Current, balance.Withdraw)
	if err != nil {
		return model.ErrInternalDB // error in balance query
	}

	// update orders `INSERT INTO orders(username, order_id, status, accrual, updated_at) VALUES($1, $2, 'NEW', $3, $4);`
	// update orders set status=$2, accrual=$3, updated_at=$4 where number=$1
	_, err = tx.Exec(ctx, changeOrdersQuery, order.Username, order.Status, order.Accrual, time.Now())
	if err != nil {
		return model.ErrInternalDB // error in order query
	}

	return tx.Commit(ctx)
}
