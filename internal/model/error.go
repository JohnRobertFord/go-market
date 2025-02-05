package model

import "errors"

var (
	ErrDataConflict = errors.New("Data conflict")
	ErrNoRows       = errors.New("No rows in result set")
	ErrValidate     = errors.New("Username or password is not correct")
)
