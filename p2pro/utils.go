package p2pro

func CelsiusToRaw(tempC float64) int {
	return int(64 * (tempC + 273.15))
}

func RawToCelsius(tempRaw int) float64 {
	return float64(tempRaw/64) - 273.15
}
