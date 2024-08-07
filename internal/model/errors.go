package model

import "errors"

var (
	ErrNoRecord           = errors.New("model: no matching record found")
	ErrInvalidCredentials = errors.New("model: invalid crendentials")
	ErrDuplicateEmail     = errors.New("model: duplicate email")
)
