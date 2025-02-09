package model

import "errors"

var (
	ErrInternal     = errors.New("Internal error")
	ErrInternalDB   = errors.New("DB error")
	ErrDataConflict = errors.New("Data conflict")
	ErrNoRows       = errors.New("No rows in result set")
	ErrValidate     = errors.New("Username or password is not correct")
	ErrData         = errors.New("Bad input data")
	ErrDataExist    = errors.New("Data exists")
	ErrCheck        = errors.New("Check error")
)
