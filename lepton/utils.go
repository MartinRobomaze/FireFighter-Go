package lepton

func CelsiusToRaw(tempC float64) int {
	return int(tempC*100 + 27315)
}

func RawToCelsius(tempRaw int) float64 {
	return float64(tempRaw-27315) / 100
}
