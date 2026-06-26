package model

import "errors"

var (
	ErrValidate = errors.New("validation failed")
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)
