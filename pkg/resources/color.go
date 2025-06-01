package resources

import "math"

const (
	sRGBGamma     = 2.4
	sRGBThreshold = 0.04045

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

func roundToNearest(value float64, nearest int) float64 {
	return math.Round(value/float64(nearest)) * float64(nearest)
}

func KelvinToMirekRounded(kelvin int32) int32 {
	mirek := KelvinToMirek(float64(kelvin))
	return int32(math.Round(mirek))
}

func MirekToKelvinRounded(mirek int32) int32 {
	kelvin := MirekToKelvin(float64(mirek))
	return int32(roundToNearest(kelvin, 100))
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

// RGBToXY converts RGB color values to XY coordinates and brightness
// Takes RGB values in range [0, 255] and returns XY coordinates and brightness [0, 100]
func RGBToXY(color RGBColor) (XYCoord, float64) {
	// Convert RGB to [0,1] range
	r := float64(color.R) / 255.0
	g := float64(color.G) / 255.0
	b := float64(color.B) / 255.0

	// Apply gamma correction
	if r > sRGBThreshold {
		r = math.Pow((r+0.055)/1.055, sRGBGamma)
	} else {
		r = r / 12.92
	}
	if g > sRGBThreshold {
		g = math.Pow((g+0.055)/1.055, sRGBGamma)
	} else {
		g = g / 12.92
	}
	if b > sRGBThreshold {
		b = math.Pow((b+0.055)/1.055, sRGBGamma)
	} else {
		b = b / 12.92
	}

	// Convert to XYZ color space
	X := r*0.4124 + g*0.3576 + b*0.1805
	Y := r*0.2126 + g*0.7152 + b*0.0722
	Z := r*0.0193 + g*0.1192 + b*0.9505

	// Calculate XY coordinates
	sum := X + Y + Z
	var x, y float64
	if sum > 0 {
		x = X / sum
		y = Y / sum
	}

	// Calculate brightness (0-100)
	brightness := Y * 100

	return XYCoord{X: x, Y: y}, brightness
}

// XYToRGB converts XY color coordinates to RGB values
// Returns red, green, blue values in the range [0, 255]
func XYToRGB(color XYCoord, brightness float64) RGBColor {
	// Convert brightness from [0,100] to [0,1]
	brightness = math.Max(0, math.Min(brightness/100.0, 1.0))

	// Calculate XY to XYZ conversion
	Y := brightness
	X := (Y / color.Y) * color.X
	Z := (Y / color.Y) * (1 - color.X - color.Y)

	// Convert XYZ to RGB
	r := X*3.2406 + Y*-1.5372 + Z*-0.4986
	g := X*-0.9689 + Y*1.8758 + Z*0.0415
	b := X*0.0557 + Y*-0.2040 + Z*1.0570

	// Apply gamma correction and convert to 0-255 range
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	if r > 0.0031308 {
		r = 1.055*math.Pow(r, 1/2.4) - 0.055
	} else {
		r = 12.92 * r
	}
	if g > 0.0031308 {
		g = 1.055*math.Pow(g, 1/2.4) - 0.055
	} else {
		g = 12.92 * g
	}
	if b > 0.0031308 {
		b = 1.055*math.Pow(b, 1/2.4) - 0.055
	} else {
		b = 12.92 * b
	}

	rC := int(math.Round(r * 255))
	gC := int(math.Round(g * 255))
	bC := int(math.Round(b * 255))
	return RGBColor{rC, gC, bC}
}
