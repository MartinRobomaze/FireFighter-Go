package utils

func ValMap(value float64, iStart float64, iStop float64, oStart float64, oStop float64) int {
	return int(oStart + (oStop-oStart)*((value-iStart)/(iStop-iStart)))
}
