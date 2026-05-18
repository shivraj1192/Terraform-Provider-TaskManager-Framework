package errs

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrNoContent    = errors.New("no content")
	ErrNotSupported = errors.New("not supported")
)
