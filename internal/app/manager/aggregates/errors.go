package aggregates

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)
