package light

import (
	"context"
	"github.com/richseviora/huego/pkg/resources/color"
	"github.com/richseviora/huego/pkg/resources/common"
)

type LightMetadata struct {
	Name      string `json:"name"`
	Archetype string `json:"archetype"`
	Function  string `json:"function"`
}

type LightOn struct {
	On bool `json:"on"`
}

type DimmingInfo struct {
	Brightness  float64 `json:"brightness"`
	MinDimLevel float64 `json:"min_dim_level"`
}

type ColorGamut struct {
	Red   color.XYCoord `json:"red"`
	Blue  color.XYCoord `json:"blue"`
	Green color.XYCoord `json:"green"`
}

type Color struct {
	XY color.XYCoord `json:"xy"`
}

type ColorInfo struct {
	XY        color.XYCoord `json:"xy"`
	Gamut     ColorGamut    `json:"gamut"`
	GamutType string        `json:"gamut_type"`
}

type ColorTemperature struct {
	Mirek int `json:"mirek"`
}

type ColorTemperatureInfo struct {
	Mirek       int  `json:"mirek"`
	MirekValid  bool `json:"mirek_valid"`
	MirekSchema struct {
		MirekMinimum int `json:"mirek_minimum"`
		MirekMaximum int `json:"mirek_maximum"`
	}
}

type LightUpdate struct {
	ID       string `json:-`
	Metadata *struct {
		Name     *string `json:"name"`
		Function *string `json:"function"`
	} `json:"metadata"`
}

type PowerUp struct {
	Preset     string `json:"preset"`
	Configured bool   `json:"configured"`
	On         struct {
		Mode string  `json:"mode"`
		On   LightOn `json:"on"`
	} `json:"on"`
	Dimming struct {
		Mode    string         `json:"mode"`
		Dimming common.Dimming `json:"dimming"`
	}
	Color struct {
		Mode      string               `json:"mode"`
		ColorTemp ColorTemperatureInfo `json:"color_temperature"`
	}
}

// Light represents the light resource data
type Light struct {
	ID        string               `json:"id"`
	IDv1      string               `json:"idv1"`
	Metadata  LightMetadata        `json:"metadata"`
	Owner     common.Reference     `json:"owner"`
	On        LightOn              `json:"on"`
	Dimming   DimmingInfo          `json:"dimming"`
	ColorTemp ColorTemperatureInfo `json:"color_temperature"`
	Color     Color                `json:"color"`
	Type      string               `json:"type"`
}

func (l Light) Identity() string {
	return l.ID
}

var (
	_ common.Identable = &Light{}
)

type LightService interface {
	GetLight(ctx context.Context, id string) (*Light, error)
	GetAllLights(ctx context.Context) (*common.ResourceList[Light], error)
	UpdateLight(ctx context.Context, update LightUpdate) error
}
