package model

import "time"

type Withdrawal struct {
	Order       Order     `db:"order_id"`
	Sum         uint64    `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}

type WithdrawalRequest struct {
	Order string `json:"order"`
	Sum   uint64 `json:"sum"`
}
