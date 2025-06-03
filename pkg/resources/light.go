package resources

import (
	"context"
	"fmt"
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
	XY XYCoord `json:"xy"`
}

type ColorInfo struct {
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

type ResourceError struct {
	Description string `json:"description"`
}

type ResourceList[T any] struct {
	Data   []T             `json:"data"`
	Errors []ResourceError `json:"errors"`
}

// LightService handles light-related API operations
type LightService struct {
	client *APIClient
}

type ResourceUpdateResponse struct {
	Errors []struct {
		Description string `json:"description"`
	} `json:"errors"`
	Data []Reference `json:"data"`
}

// NewLightService creates a new LightService instance
func NewLightService(client *APIClient) *LightService {
	return &LightService{
		client: client,
	}
}

// GetLight retrieves a single light by its ID
func (s *LightService) GetLight(ctx context.Context, id string) (*Light, error) {
	result, err := Get[ResourceList[Light]](ctx, fmt.Sprintf("/clip/v2/resource/light/%s", id), s.client)
	if result == nil || len(result.Data) == 0 {
		return nil, fmt.Errorf("light not found")
	}
	room, err := FirstOrError(result)
	if err != nil {
		return nil, fmt.Errorf("light not found")
	}
	if room.ID != id {
		return nil, fmt.Errorf("light not found")
	}
	return room, nil
}

// GetAllLights retrieves all available lights
func (s *LightService) GetAllLights(ctx context.Context) (*ResourceList[Light], error) {
	return Get[ResourceList[Light]](ctx, "/clip/v2/resource/light", s.client)
}

func (s *LightService) UpdateLight(ctx context.Context, update LightUpdate) error {
	result, err := Put[ResourceUpdateResponse](ctx, "/clip/v2/resource/light/"+update.ID, update, s.client)
	if err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("failed to update light: %v", result.Errors)
	}
	return nil
}
