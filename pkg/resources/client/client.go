package client

import "errors"

var ErrUnauthorized = errors.New("unauthorized")

var ErrServiceUnavailable = errors.New("service unavailable")

// ErrBadResponse This error is returned whenever the bridge unexpectedly returns a HTML body response.
var ErrBadResponse = errors.New("bad response")

var ErrNotFound = errors.New("not found")
