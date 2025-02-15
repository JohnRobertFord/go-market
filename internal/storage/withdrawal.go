package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/JohnRobertFord/go-market/internal/util"
	"github.com/jackc/pgx/v5"
)

var listAllWithdrawalQuery = `SELECT order_id, sum, processed_at FROM withdrawals WHERE username=$1 ORDER BY processed_at;`
var updateBalancesQuery = `INSERT INTO balances(username, current_balance, withdrawn)
								VALUES ($1, $2, $3)
								ON CONFLICT (username) DO UPDATE
								SET current_balance = $2, withdrawn = $3;`
var changeOrdersQuery = `INSERT INTO orders(username, order_id, status, updated_at)
								VALUES($1, $2, 'NEW', $3);`
var updateOrdersQuery = `UPDATE orders SET status=$2, accrual=$3 WHERE number=$1`
var changeWithdrawalsQuery = `INSERT INTO withdrawals(username, order_id, sum, processed_at)
								VALUES ($1, $2, $3, $4);`

type WithdrawalRepository struct {
	db *DB
}

func NewWithdrawalRepository(db *DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		db,
	}
}

func (wh *WithdrawalRepository) Withdrawals(ctx context.Context, user string) ([]model.Withdrawal, error) {
	rows, err := wh.db.pgPool.Query(ctx, listAllWithdrawalQuery, user)
	if err != nil {
		return nil, err
	}
	out, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Withdrawal])
	if err != nil {
		return nil, model.ErrInternal
	}
	if len(out) == 0 {
		return nil, model.ErrNoRows
	}
	return out, nil
}

func (wh *WithdrawalRepository) Withdraw(ctx context.Context, username string, m *model.WithdrawalRequest) error {
	tx, err := wh.db.pgPool.Begin(ctx)
	if err != nil {
		return model.ErrInternalDB
	}
	defer func() {
		if errRollBack := tx.Rollback(ctx); errRollBack != nil {
			fmt.Printf("rollback error: %v", errRollBack)
		}
	}()

	bUint, err := strconv.ParseUint(m.Order, 10, 64)
	if err != nil {
		return model.ErrData
	}
	// Luhn check
	if !util.Valid(bUint) {
		return model.ErrCheck
	}

	// get balance
	row, err := wh.db.pgPool.Query(ctx, getBalanceQuery, username)
	if err != nil {
		return model.ErrInternalDB
	}
	ex, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.Balance])
	if err != nil {
		return model.ErrInternalDB
	}
	if ex.Current < m.Sum {
		return model.ErrInsufficientFunds
	}

	ex.Current -= m.Sum
	ex.Withdraw += m.Sum

	// update balance
	_, err = tx.Exec(ctx, updateBalancesQuery, username, ex.Current, ex.Withdraw)
	if err != nil {
		return model.ErrInternalDB // error in balance query
	}

	// update orders `INSERT INTO orders(username, order_id, status, accrual, updated_at) VALUES($1, $2, 'NEW', $3, $4);`
	_, err = tx.Exec(ctx, changeOrdersQuery, username, m.Order, time.Now())
	if err != nil {
		return model.ErrInternalDB // error in order query
	}

	// update withdrawals
	_, err = tx.Exec(ctx, changeWithdrawalsQuery, username, m.Order, m.Sum, time.Now())
	if err != nil {
		return model.ErrInternalDB // error in withdrawn query
	}

	return tx.Commit(ctx)
}
