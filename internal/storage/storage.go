package storage

import "errors"

var (
	ErrNotFound = errors.New("entity not found")
	ErrURLExists   = errors.New("url exists")
)