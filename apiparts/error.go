package epicrm_apiparts

import (
	"errors"
)

// Inspired by https://pkg.go.dev/io/fs#ErrNotExist
var (
	ErrNotFound    = errors.New("not found")
	ErrServerError = errors.New("internal server error")
)
