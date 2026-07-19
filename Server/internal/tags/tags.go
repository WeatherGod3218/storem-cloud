package tags

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func hsvToRGB(h, s, v float64) (int, int, int) { //STRAIGHT FROM ROBLOX STUDIO LOOOOLLLLL
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var (
		r float64
		g float64
		b float64
	)

	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return int((r + m) * 255), int((g + m) * 255), int((b + m) * 255)
}

func hexFromHSV(h, s, v float64) string {
	r, g, b := hsvToRGB(h, s, v)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func GenerateRandomTagHexColor() string {
	hue := rand.IntN(360)
	hex := hexFromHSV(float64(hue), 1, 0.5)
	return hex
}
