package pkg

import (
	"context"
	"net/http"
)

type ClientInterface interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
