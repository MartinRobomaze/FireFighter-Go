package firefinder

import (
	"FireFighter/lepton"
	"github.com/MartinRobomaze/gocv"
	"github.com/sirupsen/logrus"
	"image"
	"log"
	"testing"
)

func TestFireFinder(t *testing.T) {
	d, err := lepton.New()
	if err != nil {
		logrus.Fatal(err)
	}

	frameChan, err := d.StartStreaming()
	if err != nil {
		logrus.Fatal(err)
	}

	ff := New(40, 1.5, 0.3, image.Point{X: 5, Y: 5}, []int{160, 120}, 57, 42.5)

	for {
		data := <-frameChan

		mat := gocv.NewMatFromInt16Arr(120, 160, gocv.MatTypeCV16UC1, data.Data)
		isFire, fireData := ff.FindFire(mat)
		if isFire {
			log.Printf("Fire temp: %g\tX: %d\tY: %d", fireData.FireTemperature, fireData.FireLocation.X, fireData.FireLocation.Y)
		}
	}
}
