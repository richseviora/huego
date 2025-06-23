package client

import (
	"context"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/common"
	room2 "github.com/richseviora/huego/pkg/resources/room"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_APIClientParse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			t.Error("failed to read body", err)
			return
		}
		if string(body) != `{"children":[{"rid":"1234","rtype":"light"}],"metadata":{"name":"TEST ROOM","archetype":"bedroom"}}` {
			t.Errorf("Expected body to be %s, got %s", `{"children":[{"rid":"1234","rtype":"light"}],"metadata":{"name":"TEST ROOM","archetype":"bedroom"}}`, string(body))
		}
		_, err = w.Write([]byte(`{"data":[{"id":"123","type":"room","metadata":{"archetype":"kitchen","name":"Kitchen"}}],"errors":[]}`))
		if err != nil {
			t.Error("failed to write response", err)
		}
	}))
	defer server.Close()
	client := NewAPIClient(server.URL, "1234567890", logger.NoopLogger{})
	result, err := client.roomService.GetRoom(context.Background(), "123")
	if err != nil {
		t.Error(err)
	}
	expected := &room2.RoomData{
		ID:   "123",
		Type: "room",
		Metadata: room2.RoomMetadata{
			Name:      "Kitchen",
			Archetype: common.Kitchen,
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", result, expected)
	}
}

func Test_APIClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			t.Error("failed to read body", err)
			return
		}
		if string(body) != `{"children":[{"rid":"1234","rtype":"light"}],"metadata":{"name":"TEST ROOM","archetype":"bedroom"}}` {
			t.Errorf("Expected body to be %s, got %s", `{"children":[{"rid":"1234","rtype":"light"}],"metadata":{"name":"TEST ROOM","archetype":"bedroom"}}`, string(body))
		}
		_, err = w.Write([]byte(`{"data":[{"rid":"1234","rtype":"room"}],"errors":[]}`))
		if err != nil {
			t.Error("failed to write response", err)
		}
	}))
	defer server.Close()
	client := NewAPIClient(server.URL, "1234567890", logger.NoopLogger{})
	result, err := client.roomService.CreateRoom(context.Background(), room2.RoomCreate{
		Children: []common.Reference{
			common.Reference{
				RID:   "1234",
				RType: "light",
			},
		},
		Metadata: room2.RoomMetadata{
			Name:      "TEST ROOM",
			Archetype: common.Bedroom,
		},
	})
	if err != nil {
		t.Error(err)
	}
	expected := &common.Reference{
		RID:   "1234",
		RType: "room",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", result, expected)
	}
}
