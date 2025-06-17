package motion

import (
	"context"
	"github.com/richseviora/huego/pkg/resources/common"
	"time"
)

type MotionReport struct {
	Changed time.Time `json:"changed"`
	Motion  bool      `json:"motion"`
}
type Motion struct {
	Motion       bool         `json:"motion"`
	MotionValid  bool         `json:"motion_valid"`
	MotionReport MotionReport `json:"motion_report"`
}
type Sensitivity struct {
	Status         string `json:"status"`
	Sensitivity    int    `json:"sensitivity"`
	SensitivityMax int    `json:"sensitivity_max"`
}
type SensitivityUpdate struct {
	Sensitivity int `json:"sensitivity"`
}
type Data struct {
	ID          string           `json:"id"`
	IDV1        string           `json:"id_v1"`
	Owner       common.Reference `json:"owner"`
	Enabled     bool             `json:"enabled"`
	Motion      Motion           `json:"motion"`
	Sensitivity Sensitivity      `json:"sensitivity"`
	Type        string           `json:"type"`
}

type UpdateRequest struct {
	Enabled     bool               `json:"enabled"`
	Sensitivity *SensitivityUpdate `json:"sensitivity,omitempty"`
}

var (
	_ common.Identable = &Data{}
)

func (d Data) Identity() string {
	return d.ID
}

type Service interface {
	GetAllMotion(ctx context.Context) (*common.ResourceList[Data], error)
	GetMotion(ctx context.Context, id string) (*Data, error)
	UpdateMotion(ctx context.Context, id string, update UpdateRequest) (*common.Reference, error)
}
