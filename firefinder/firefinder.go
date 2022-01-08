package firefinder

import (
	"github.com/MartinRobomaze/gocv"
	"image"
)

type FireFinder struct {
	MaxFireWidthHeightRatio float64
	DWidthHeightRatio       float64
	Resolution              []int
	ViewingAngleHorizontal  float64
	ViewingAngleVertical    float64

	tempThresholdRaw int
	prevFireBBox     image.Rectangle
	dilationKernel   gocv.Mat
}

func New(
	tempThreshold float64,
	maxFireWHRatio float64,
	DWHRatio float64,
	DilationKernelSize image.Point,
	resolution []int,
	viewAngleHor float64,
	viewAngleVer float64,
) *FireFinder {
	return &FireFinder{
		MaxFireWidthHeightRatio: maxFireWHRatio,
		DWidthHeightRatio:       DWHRatio,
		Resolution:              resolution,
		ViewingAngleHorizontal:  viewAngleHor,
		ViewingAngleVertical:    viewAngleVer,
		tempThresholdRaw:        CelsiusToRaw(tempThreshold),
		prevFireBBox:            image.Rect(1, 1, 1, 1),
		dilationKernel:          gocv.Ones(DilationKernelSize.Y, DilationKernelSize.X, gocv.MatTypeCV8U),
	}
}

func (f *FireFinder) FindFire(data gocv.Mat) (firePresent bool, fireData *FireData) {
	_, maxTemp, _, maxTempLoc := gocv.MinMaxLoc(data)
	if int(maxTemp) < f.tempThresholdRaw {
		return false, nil
	}

	gocv.Threshold(data, &data, float32(f.tempThresholdRaw), 255, gocv.ThresholdToZero)

	data8Bit := gocv.NewMat()

	gocv.Normalize(data, &data8Bit, 0, 255, gocv.NormMinMax)
	data8Bit.ConvertTo(&data8Bit, gocv.MatTypeCV8U)

	gocv.Dilate(data8Bit, &data8Bit, f.dilationKernel)

	contours := gocv.FindContours(data8Bit, gocv.RetrievalList, gocv.ChainApproxSimple)

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)

		if gocv.PointPolygonTest(contour, maxTempLoc, false) < 0 {
			continue
		}

		fireBBox := gocv.BoundingRect(contour)

		whRatio := float64(fireBBox.Dx()) / float64(fireBBox.Dy())
		prevWhRatio := float64(f.prevFireBBox.Dx()) / float64(f.prevFireBBox.Dy())

		if whRatio-prevWhRatio > f.DWidthHeightRatio {
			return false, nil
		}

		if whRatio < f.MaxFireWidthHeightRatio {
			return true, &FireData{
				FireTemperature: RawToCelsius(int(maxTemp)),
				FireLocation:    maxTempLoc,
			}
		}
	}

	return false, nil
}

func (f *FireFinder) FireCoordsToAngles(fireCoords image.Point) (angleHor float64, angleVer float64) {
	angleHor = (float64(fireCoords.X) / float64(f.Resolution[0]) * f.ViewingAngleHorizontal) -
		f.ViewingAngleHorizontal/2

	angleVer = f.ViewingAngleVertical/2 -
		(float64(fireCoords.Y) / float64(f.Resolution[1]) * f.ViewingAngleVertical)

	return
}

func CelsiusToRaw(tempC float64) int {
	return int(tempC*100 + 27315)
}

func RawToCelsius(tempRaw int) float64 {
	return float64(tempRaw-27315) / 100
}
