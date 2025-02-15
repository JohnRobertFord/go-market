package model

type Accrual struct {
	Order    string `json:"order"`
	Status   string `json:"status"`
	Accrual  uint64 `json:"accrual,omitempty"`
	Username string `json:"-"`
}
