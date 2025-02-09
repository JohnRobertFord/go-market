package model

import "time"

type Order struct {
	Username  string
	OrderID   uint64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
