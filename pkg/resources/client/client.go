package client

import "errors"

var ErrUnauthorized = errors.New("unauthorized")

var ErrServiceUnavailable = errors.New("service unavailable")
