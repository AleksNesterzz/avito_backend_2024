package apierr

import "errors"

var (
	ErrNoAuth = errors.New("user not authorized")
)
