package behavior_script

import (
	"context"
	"github.com/richseviora/huego/pkg/resources/common"
)

type SchemaReference struct {
	Ref string `json:"$ref"`
}

type Metadata struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// Data - the canonical name is in the Metadata field. 'Motion Sensor' is an example
type Data struct {
	ID                  string          `json:"id"`
	Type                string          `json:"type"`
	Description         string          `json:"description"`
	ConfigurationSchema SchemaReference `json:"configuration_schema"`
	TriggerSchema       SchemaReference `json:"trigger_schema,omitempty"`
	StateSchema         SchemaReference `json:"state_schema,omitempty"`
	Version             string          `json:"version"`
	Metadata            Metadata        `json:"metadata"`
	SupportedFeatures   []string        `json:"supported_features"`
	MaxNumberInstances  *int            `json:"max_number_instances,omitempty"`
}

func (d Data) Identity() string {
	return d.ID
}

var (
	_ common.Identable = &Data{}
)

type Service interface {
	GetAllBehaviorScripts(ctx context.Context) (*common.ResourceList[Data], error)
	GetBehaviorScript(ctx context.Context, id string) (*Data, error)
}
