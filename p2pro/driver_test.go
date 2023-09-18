package p2pro

import (
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
	"log"
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
		mat := <-frameChan

		mat16Bit := gocv.NewMat()
		mat.ConvertTo(&mat16Bit, gocv.MatTypeCV16U)

		log.Println(mat16Bit.Type())

		data8Bit := gocv.NewMat()

		gocv.Normalize(mat16Bit, &data8Bit, 0, 255, gocv.NormMinMax)
		data8Bit.ConvertTo(&data8Bit, gocv.MatTypeCV8U)

		resizedImg := gocv.NewMat()
		gocv.Resize(data8Bit, &resizedImg, image.Point{X: 640, Y: 480}, 0, 0, gocv.InterpolationLinear)
		window.IMShow(resizedImg)
		window.WaitKey(1)
	}
}
