package app

import "errors"

var (
	ErrForbidden = errors.New("user don't have permission")
)
