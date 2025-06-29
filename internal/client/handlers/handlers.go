package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/client"
	"github.com/richseviora/huego/pkg/resources/common"
	"io"
	"net/http"
)

func GetSingularResource[T common.Identable](id string, path string, ctx context.Context, client common.RequestProcessor, resourceName string) (*T, error) {
	result, err := Get[common.ResourceList[T]](ctx, path, client)
	if err != nil || result == nil {
		return nil, err
	}
	if len(result.Data) == 0 {
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

func FirstOrError[T any](list *common.ResourceList[T]) (*T, error) {
	if list == nil || len(list.Data) == 0 {
		return nil, fmt.Errorf("resource not found")
	}
	return &list.Data[0], nil
}

func CreateResource[T any](path string, ctx context.Context, create T, client common.RequestProcessor, resourceName string) (*common.Reference, error) {
	result, err := Post[common.ResourceUpdateResponse](ctx, path, create, client)
	if err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		client.Logger().Error("Failed to Create Resource", map[string]interface{}{
			"resourceName": resourceName,
			"errors":       result.Errors,
		})
		return nil, fmt.Errorf("failed to create resource %s: %v", resourceName, result.Errors)
	}
	return &result.Data[0], nil
}

func UpdateResource[T any](path string, ctx context.Context, create T, client common.RequestProcessor, resourceName string) (*common.Reference, error) {
	result, err := Put[common.ResourceUpdateResponse](ctx, path, create, client)
	if err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("failed to update resource %s: %v", resourceName, result.Errors)
	}
	return &result.Data[0], nil
}

func decodeResponse[T any](body []byte, url string, l logger.Logger) (*T, error) {
	var result T
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&result); err != nil {
		l.Error("Failed to Decode Response", map[string]interface{}{
			"url":   url,
			"error": err,
		})
		return nil, client.ErrBadResponse
	}
	return &result, nil
}

func Get[T any](ctx context.Context, path string, c common.RequestProcessor) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL()+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return decodeResponse[T](bodyBytes, req.URL.String(), c.Logger())
}

// Delete performs a DELETE request for the specified resource
func Delete(ctx context.Context, path string, c common.RequestProcessor) error {
	req, err := http.NewRequest(http.MethodDelete, c.BaseURL()+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("delete request failed with status: %s", resp.Status)
	}

	return nil
}

// Post performs a POST request and unmarshals the response into the provided type
func Post[T any](ctx context.Context, path string, body interface{}, c common.RequestProcessor) (*T, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	c.Logger().Trace("Sending POST request", map[string]interface{}{
		"path": c.BaseURL() + path,
		"body": body,
	})
	req, err := http.NewRequest(http.MethodPost, c.BaseURL()+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	c.Logger().Trace("Received POST Response", map[string]interface{}{
		"body":   string(bodyBytes),
		"status": resp.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return decodeResponse[T](bodyBytes, req.URL.String(), c.Logger())
}

func Put[T any](ctx context.Context, path string, body interface{}, c common.RequestProcessor) (*T, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	url := c.BaseURL() + path
	c.Logger().Trace("Sending PUT request to %s with body: %s\n", map[string]interface{}{
		"url":  url,
		"body": body,
	})
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	c.Logger().Trace("Received Response Body", map[string]interface{}{
		"body": string(bodyBytes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return decodeResponse[T](bodyBytes, req.URL.String(), c.Logger())
}
