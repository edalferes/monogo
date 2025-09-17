package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrMissingCredentials = errors.New("username and password are required")
	ErrInvalidData        = errors.New("invalid data")
	ErrUserAlreadyExists  = errors.New("user already exists")
)
