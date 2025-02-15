package model

type Balance struct {
	Current  uint64 `json:"current" db:"current"`
	Withdraw uint64 `json:"withdrawn" db:"withdrawn"`
}
