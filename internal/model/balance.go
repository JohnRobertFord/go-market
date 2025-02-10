package model

type Balance struct {
	Current  uint64 `json:"current"`
	Withdraw uint64 `json:"withdrawn"`
}
