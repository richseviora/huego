package behavior_instance

import (
	"context"
	"github.com/richseviora/huego/pkg/resources/common"
)

type State struct {
	// The only observed value is "device"
	SourceType string `json:"source_type"`
	// Should match the model ID of the device.
	ModelID string `json:"model_id"`
}
type DaylightSensitivity struct {
	DarkThreshold int `json:"dark_threshold"`
	Offset        int `json:"offset"`
}
type Settings struct {
	DaylightSensitivity DaylightSensitivity `json:"daylight_sensitivity"`
}
type Action struct {
	// The reference here must be a Scene.
	Recall common.Reference `json:"recall"`
}
type RecallSingle struct {
	Action Action `json:"action"`
}
type OnMotion struct {
	RecallSingle []RecallSingle `json:"recall_single"`
}
type After struct {
	Minutes int `json:"minutes"`
}
type OnNoMotion struct {
	After        After          `json:"after"`
	RecallSingle []RecallSingle `json:"recall_single"`
}
type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}
type StartTime struct {
	Time Time   `json:"time"`
	Type string `json:"type"`
}
type TimeSlots struct {
	OnMotion   OnMotion   `json:"on_motion"`
	OnNoMotion OnNoMotion `json:"on_no_motion"`
	StartTime  StartTime  `json:"start_time"`
}
type When struct {
	Timeslots []TimeSlots `json:"timeslots"`
}

// Where The references here should be rooms.
type Where struct {
	Group common.Reference `json:"group"`
}
type Configuration struct {
	Settings Settings `json:"settings"`
	// The reference here should be a Device ID for a sensor.
	Source common.Reference `json:"source"`
	When   When             `json:"when"`
	Where  []Where          `json:"where"`
}
type Dependees struct {
	Target common.Reference `json:"target"`
	// The only observed value so far is "critical"
	Level string `json:"level"`
	// the only observed value so far is "ResourceDependee"
	Type string `json:"type"`
}
type Metadata struct {
	Name string `json:"name"`
}
type Data struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	// This should refer to a valid BehaviorScript ID
	ScriptID      string        `json:"script_id"`
	Enabled       bool          `json:"enabled"`
	State         State         `json:"state"`
	Configuration Configuration `json:"configuration"`
	/*
		 Dependees - there should be a dependee for each of:
		- the room
		- the scene
		- the motion sensor device ID
		- the motion sensor motion ID
		- the motion sensor light_level ID
	*/
	Dependees []Dependees `json:"dependees"`
	Status    string      `json:"status"`
	LastError string      `json:"last_error"`
	Metadata  Metadata    `json:"metadata"`
}

type CreateRequest struct {
	ScriptID      string        `json:"script_id"`
	Enabled       bool          `json:"enabled"`
	Configuration Configuration `json:"configuration"`
	Metadata      *Metadata     `json:"metadata"`
}

type UpdateRequest struct {
	Configuration *Configuration `json:"configuration"`
	Enabled       *bool          `json:"enabled"`
	Metadata      *Metadata      `json:"metadata"`
}

var (
	_ common.Identable = &Data{}
)

func (d Data) Identity() string {
	return d.ID
}

type Service interface {
	GetAllBehaviorInstances(ctx context.Context) (*common.ResourceList[Data], error)
	GetBehaviorInstance(ctx context.Context, id string) (*Data, error)
	UpdateBehaviorInstance(ctx context.Context, id string, update UpdateRequest) (*common.Reference, error)
	CreateBehaviorInstance(ctx context.Context, create CreateRequest) (*common.Reference, error)
	DeleteBehaviorInstance(ctx context.Context, id string) error
}
