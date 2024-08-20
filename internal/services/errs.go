package services

import "errors"

var (
	ErrKeyExists          = errors.New("key already exists")
	ErrKeyNotFound        = errors.New("key not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
