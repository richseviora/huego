package common

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestArea struct {
	Name Area `json:"name"`
}

func TestArea_MarshalJSON(t *testing.T) {

	testCases := []struct {
		name     string
		a        Area
		expected string
	}{
		{"converts kitchen", Kitchen, `{"name":"kitchen"}`},
		{"converts living_room", LivingRoom, `{"name":"living_room"}`},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(&TestArea{
				Name: tt.a,
			})
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("json.Marshal() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

func TestArea_ParseArea(t *testing.T) {
	testCases := []struct {
		name     string
		expected Area
		input    string
	}{
		{"converts kitchen", Kitchen, "kitchen"},
		{"converts living_room", LivingRoom, "living_room"},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			area, err := ParseArea(tt.input)
			if err != nil {
				t.Errorf("ParseArea() error = %v", err)
				return
			}
			if area != tt.expected {
				t.Errorf("ParseArea() = %v, want %v", area, tt.expected)
			}
		})
	}
}

func TestArea_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		expected Area
		input    string
	}{
		{"converts kitchen", Kitchen, `{"name":"kitchen"}`},
		{"converts living_room", LivingRoom, `{"name":"living_room"}`},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var result TestArea
			err := json.Unmarshal([]byte(tt.input), &result)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			expected := TestArea{
				Name: tt.expected,
			}
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("json.Marshal() = %+v, want %+v", result, expected)
			}
		})
	}
}
