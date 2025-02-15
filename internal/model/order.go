package model

import "time"

type Order struct {
	Username  string    `json:"-" db:"username,omitempty"`
	OrderID   string    `json:"number" db:"order_id"`
	Status    string    `json:"status" db:"status"`
	Accrual   uint64    `json:"accrual,omitempty" db:"accrual,omitempty"`
	CreatedAt time.Time `json:"uploaded_at" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}
