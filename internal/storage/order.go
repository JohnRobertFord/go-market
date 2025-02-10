package storage

import (
	"context"
	"errors"
	"time"

	"github.com/JohnRobertFord/go-market/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	createOrderQuery = `INSERT INTO orders("username", "order_id", "status", "created_at", "updated_at") VALUES($1, $2, $3, $4, $4)`
	listOrdersQuery  = `SELECT username, order_id, status, created_at, updated_at FROM orders WHERE username=$1 ORDER BY updated_at DESC`
	searchQuery      = `SELECT "username", "order_id", "status", "created_at", "updated_at" FROM orders WHERE order_id=$1`
)

type OrderRepository struct {
	db *DB
}

func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}

func (o *OrderRepository) Create(ctx context.Context, order *model.Order) error {

	row, err := o.db.pgPool.Query(ctx, searchQuery, order.OrderID)
	if err != nil {
		return err
	}
	ex, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.Order])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = o.db.pgPool.Exec(ctx, createOrderQuery, order.Username, order.OrderID, order.Status, time.Now().Format(time.RFC3339))
			if e, ok := err.(*pgconn.PgError); ok {
				if e.Code == "23505" {
					return model.ErrInternalDB
				} else {
					return model.ErrInternal
				}
			}
			return nil
		} else {
			return model.ErrInternal
		}
	}
	if ex.Username == order.Username {
		return model.ErrDataExist
	} else {
		return model.ErrDataConflict
	}

}

// статус заказа:
// - `NEW` — заказ загружен в систему, но не попал в обработку;
// - `PROCESSING` — вознаграждение за заказ рассчитывается;
// - `INVALID` — система расчёта вознаграждений отказала в расчёте;
// - `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.
func (o *OrderRepository) List(ctx context.Context, user string) ([]model.Order, error) {

	rows, err := o.db.pgPool.Query(ctx, listOrdersQuery, user)
	if err != nil {
		return nil, err
	}
	out, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return out, nil
}
