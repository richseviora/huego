package color

import (
	"fmt"
	"testing"
)

func TestMirekToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		mirek    float64
		expected float64
	}{
		{"minimum mirek", 153, 6500},
		{"maximum mirek", 500, 2000},
		{"below minimum mirek", 100, 6500},
		{"above maximum mirek", 600, 2000},
		{"middle range", 250, 4000},
		{"exact conversion 1", 200, 5000},
		{"exact conversion 2", 400, 2500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MirekToKelvin(tt.mirek)
			if result != tt.expected {
				t.Errorf("MirekToKelvin(%v) = %v, want %v", tt.mirek, result, tt.expected)
			}
		})
	}
}

func Test_KelvinToMirekRoundedAndBack(t *testing.T) {
	tests := []struct {
		name   string
		kelvin int32
	}{
		{"cold", 6500},
		{"warm", 2200},
		{"mid", 2500},
	}
	for i := range 40 {
		value := (i + 22) * 100
		tests = append(tests, struct {
			name   string
			kelvin int32
		}{fmt.Sprintf("value %v", value), int32(value)})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mirek := KelvinToMirekRounded(tt.kelvin)
			result := MirekToKelvinRounded(mirek)
			if result != tt.kelvin {
				t.Errorf("KelvinToMirekAndBack(%v = %v = %v)", tt.kelvin, mirek, result)
			}
		})
	}
}
