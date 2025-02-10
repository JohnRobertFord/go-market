package model

import "errors"

var (
	ErrInternal     = errors.New("internal error")
	ErrInternalDB   = errors.New("db error")
	ErrDataConflict = errors.New("data conflict")
	ErrNoRows       = errors.New("no rows in result set")
	ErrValidate     = errors.New("username or password is not correct")
	ErrData         = errors.New("bad input data")
	ErrDataExist    = errors.New("data exists")
	ErrCheck        = errors.New("check error")
)
