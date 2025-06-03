package resources

import (
	"context"
	"fmt"
)

type Reference struct {
	RID   string `json:"rid"`
	RType string `json:"rtype"`
}

type Dimming struct {
	Brightness float64 `json:"brightness"`
}

type XYCoord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
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

type Identable interface {
	Identity() string
}

type ResourcePathable interface {
	CollectionPath() string
	ResourcePath(id string) string
}

func GetSingularResource[T Identable](id string, path string, ctx context.Context, client *APIClient, resourceName string) (*T, error) {
	result, err := Get[ResourceList[T]](ctx, path, client)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 {
		if result.Errors != nil && len(result.Errors) > 0 {
			return nil, fmt.Errorf(result.Errors[0].Description)
		}
		return nil, fmt.Errorf("resource ID %s of type %s not found", id, resourceName)
	}
	resource, err := FirstOrError[T](result)
	if err != nil {
		return nil, fmt.Errorf("resource ID %s of type %s not found", id, resourceName)
	}
	if (*resource).Identity() != id {
		return nil, fmt.Errorf("resource ID %s of type %s not matched", id, resourceName)
	}
	return resource, nil
}

func FirstOrError[T any](list *ResourceList[T]) (*T, error) {
	if list == nil || len(list.Data) == 0 {
		return nil, fmt.Errorf("resource not found")
	}
	return &list.Data[0], nil
}

func CreateResource[T any](path string, ctx context.Context, create T, client *APIClient, resourceName string) (*Reference, error) {
	result, err := Post[ResourceUpdateResponse](ctx, path, create, client)
	if err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("failed to create resource %s: %v", resourceName, result.Errors)
	}
	return &result.Data[0], nil
}

func UpdateResource[T any](path string, ctx context.Context, create T, client *APIClient, resourceName string) (*Reference, error) {
	result, err := Put[ResourceUpdateResponse](ctx, path, create, client)
	if err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("failed to update resource %s: %v", resourceName, result.Errors)
	}
	return &result.Data[0], nil
}
