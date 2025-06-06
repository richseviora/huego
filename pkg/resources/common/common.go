package common

import (
	"context"
	"net/http"
)

type Reference struct {
	RID   string `json:"rid"`
	RType string `json:"rtype"`
}

type Dimming struct {
	Brightness float64 `json:"brightness"`
}

type ResourceError struct {
	Description string `json:"description"`
}

type ResourceList[T any] struct {
	Data   []T             `json:"data"`
	Errors []ResourceError `json:"errors"`
}

type ResourceUpdateResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data []Reference `json:"data"`
}

type RequestProcessor interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
	BaseURL() string
}

type Identable interface {
	Identity() string
}

type ResourcePathable interface {
	CollectionPath() string
	ResourcePath(id string) string
}
