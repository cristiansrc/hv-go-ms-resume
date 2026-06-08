package entity

import "errors"

var (
	// ErrNotFound is a sentinel error for when an entity is not found.
	ErrNotFound = errors.New("not found")
)
