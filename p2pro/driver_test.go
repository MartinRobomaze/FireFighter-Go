package p2pro

import (
	"github.com/MartinRobomaze/gocv"
	"github.com/sirupsen/logrus"
	"image"
	"testing"
)

func TestDriver(t *testing.T) {
	d, err := New("/dev/video0")
	if err != nil {
		logrus.Fatal(err)
	}

	frameChan, err := d.StartStreaming()
	if err != nil {
		logrus.Fatal(err)
	}

	window := gocv.NewWindow("mopslik")

	for {
		frame := <-frameChan
		mat16Bit := gocv.NewMatFromInt16Arr(VideoHeight, VideoWidth, gocv.MatTypeCV16UC1, frame)
		mat8Bit := gocv.NewMat()
		
		gocv.Normalize(mat16Bit, &mat8Bit, 0, 255, gocv.NormMinMax)
		mat8Bit.ConvertTo(&mat8Bit, gocv.MatTypeCV8U)

		resizedImg := gocv.NewMat()
		gocv.Resize(mat16Bit, &resizedImg, image.Point{X: 640, Y: 480}, 0, 0, gocv.InterpolationLinear)
		window.IMShow(resizedImg)
		window.WaitKey(1)
	}
}
