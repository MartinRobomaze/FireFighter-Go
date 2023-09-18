package p2pro

import (
	"gocv.io/x/gocv"
)

const (
	VID = 0x0bda
	PID = 0x5830
)

type Driver struct {
	Camera          *gocv.VideoCapture
	closeStreamChan chan bool
}

func New() (*Driver, error) {
	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return nil, err
	}

	return &Driver{
		Camera: cam,
	}, nil
}

func (d *Driver) StartStreaming() (frameChan chan *gocv.Mat, err error) {
	frameChan = make(chan *gocv.Mat)
	d.closeStreamChan = make(chan bool)

	go d.stream(frameChan)

	return frameChan, nil
}

func (d *Driver) Close() error {
	d.closeStreamChan <- true

	return d.Camera.Close()
}

func (d *Driver) stream(frameChan chan *gocv.Mat) {
	for {
		select {
		case <-d.closeStreamChan:
			return
		default:
			mat := gocv.NewMatWithSize(256, 384, gocv.MatTypeCV16U)
			d.Camera.Read(&mat)
			thermalData := mat.RowRange(mat.Rows()/2, mat.Rows())
			frameChan <- &thermalData
		}
	}
}
