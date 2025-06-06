package color

import (
	"math"

	"github.com/nuqz/col2xy"
)

const (
	maxKelvin = 6500
	minKelvin = 2000
	maxMirek  = 500
	minMirek  = 153
)

type RGBColor struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type XYCoord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func KelvinToMirekRounded(kelvin int32) int32 {
	mirek := KelvinToMirek(float64(kelvin))
	return int32(math.Round(mirek))
}

func MirekToKelvinRounded(mirek int32) int32 {
	kelvin := MirekToKelvin(float64(mirek))
	return int32(roundToNearest(kelvin, 100))
}

func roundToNearest(value float64, nearest int) float64 {
	return math.Round(value/float64(nearest)) * float64(nearest)
}

// KelvinToMirek converts color temperature from Kelvin to mirek value
// Mirek = 1,000,000 / color temperature (Kelvin)
func KelvinToMirek(kelvin float64) float64 {
	// Clamp Kelvin value to valid range
	kelvin = math.Max(minKelvin, math.Min(maxKelvin, kelvin))

	// Convert to mirek
	mirek := 1000000 / kelvin

	// Clamp mirek value to valid range
	return math.Max(minMirek, math.Min(maxMirek, mirek))
}

// MirekToKelvin converts color temperature from mirek to Kelvin value
// Kelvin = 1,000,000 / mirek
func MirekToKelvin(mirek float64) float64 {
	// Clamp mirek value to valid range
	mirek = math.Max(minMirek, math.Min(maxMirek, mirek))

	// Convert to Kelvin
	kelvin := 1000000 / mirek

	// Clamp Kelvin value to valid range
	return math.Max(minKelvin, math.Min(maxKelvin, kelvin))
}

func RGBtoXY2(c RGBColor) XYCoord {
	resultX, resultY := col2xy.RGB2XY(byte(c.R), byte(c.B), byte(c.G))
	return XYCoord{
		X: resultX,
		Y: resultY,
	}
}
