package lepton

import (
	"github.com/MartinRobomaze/gocv"
	"github.com/sirupsen/logrus"

	"image"
	"testing"
)

func TestDriver(t *testing.T) {
	d, err := New()
	if err != nil {
		logrus.Fatal(err)
	}

	frameChan, err := d.StartStreaming()
	if err != nil {
		logrus.Fatal(err)
	}

	window := gocv.NewWindow("mopslik")

	for {
		data := <-frameChan

		mat := gocv.NewMatFromInt16Arr(120, 160, gocv.MatTypeCV16UC1, data.Data)

		data8Bit := gocv.NewMat()

		gocv.Normalize(mat, &data8Bit, 0, 255, gocv.NormMinMax)
		data8Bit.ConvertTo(&data8Bit, gocv.MatTypeCV8U)

		resizedImg := gocv.NewMat()
		gocv.Resize(data8Bit, &resizedImg, image.Point{X: 640, Y: 480}, 0, 0, gocv.InterpolationLinear)
		window.IMShow(resizedImg)
		window.WaitKey(1)
	}
}
