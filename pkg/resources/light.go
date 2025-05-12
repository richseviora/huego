package resources

import (
	"context"
	"fmt"

	"github.com/richseviora/huego/pkg"
)

type Reference struct {
	RID   string `json:"rid"`
	RType string `json:"rtype"`
}

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

type Dimming struct {
	Brightness float64 `json:"brightness"`
}

type XYCoord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ColorGamut struct {
	Red   XYCoord `json:"red"`
	Blue  XYCoord `json:"blue"`
	Green XYCoord `json:"green"`
}

type Color struct {
	XY        XYCoord    `json:"xy"`
	Gamut     ColorGamut `json:"gamut"`
	GamutType string     `json:"gamut_type"`
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

type PowerUp struct {
	Preset     string `json:"preset"`
	Configured bool   `json:"configured"`
	On         struct {
		Mode string  `json:"mode"`
		On   LightOn `json:"on"`
	} `json:"on"`
	Dimming struct {
		Mode    string  `json:"mode"`
		Dimming Dimming `json:"dimming"`
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
	Owner     Reference            `json:"owner"`
	On        LightOn              `json:"on"`
	Dimming   DimmingInfo          `json:"dimming"`
	ColorTemp ColorTemperatureInfo `json:"color_temperature"`
	Color     Color                `json:"color"`
	Type      string               `json:"type"`
}

type ResourceList[T any] struct {
	Data   []T      `json:"data"`
	Errors []string `json:"errors"`
}

// LightService handles light-related API operations
type LightService struct {
	client *pkg.APIClient
}

// NewLightService creates a new LightService instance
func NewLightService(client *pkg.APIClient) *LightService {
	return &LightService{
		client: client,
	}
}

// GetLight retrieves a single light by its ID
func (s *LightService) GetLight(ctx context.Context, id string) (*Light, error) {
	return pkg.Get[Light](ctx, fmt.Sprintf("/clip/v2/resource/light/%s", id), s.client)
}

// GetAllLights retrieves all available lights
func (s *LightService) GetAllLights(ctx context.Context) (*ResourceList[Light], error) {
	return pkg.Get[ResourceList[Light]](ctx, "/clip/v2/resource/light", s.client)
}
