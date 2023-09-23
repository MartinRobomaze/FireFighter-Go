package p2pro

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	dev "github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"reflect"
)

const (
	rawVideoWidth  = 256
	rawVideoHeight = 384

	VideoWidth  = rawVideoWidth
	VideoHeight = rawVideoHeight / 2
)

type Driver struct {
	Camera          *dev.Device
	closeStreamChan chan bool
}

func New(videoDevice string) (*Driver, error) {
	device, err := dev.Open(
		videoDevice,
		dev.WithPixFormat(
			v4l2.PixFormat{
				Width:       uint32(rawVideoWidth),
				Height:      uint32(rawVideoHeight),
				PixelFormat: v4l2.PixelFmtYUYV,
				Field:       v4l2.FieldNone,
			}),
		dev.WithFPS(25),
	)

	if err != nil {
		return nil, err
	}

	return &Driver{
		Camera: device,
	}, nil
}

func (d *Driver) StartStreaming() (frameChan chan []int16, err error) {
	err = d.Camera.Start(context.Background())
	if err != nil {
		return nil, err
	}

	frameChan = make(chan []int16)
	go d.stream(frameChan)

	return frameChan, nil
}

func (d *Driver) Close() error {
	err := d.Camera.Stop()
	if err != nil {
		return err
	}

	return d.Camera.Close()
}

func (d *Driver) stream(dataChan chan []int16) {
	for {
		frame := <-d.Camera.GetOutput()

		thermalData := fromBuffer(frame[len(frame)/2:], reflect.TypeOf(int16(0))).([]int16)

		select {
		case dataChan <- thermalData:
		default:
			logrus.Info("frame dropped")
		}
	}
}

func fromBuffer(data []byte, dataType reflect.Type) interface{} {
	// Calculate the element size for the given data type
	elemSize := dataType.Size()

	// Calculate the number of elements in the buffer
	numElements := len(data) / int(elemSize)

	// Create an empty slice with the specified type
	sliceType := reflect.SliceOf(dataType)
	slice := reflect.MakeSlice(sliceType, numElements, numElements)

	// Iterate over the buffer and convert bytes to the desired type
	for i := 0; i < numElements; i++ {
		offset := i * int(elemSize)
		elemData := data[offset : offset+int(elemSize)]

		// Use binary.Read to convert bytes to the desired type
		elem := reflect.New(dataType).Elem()
		err := binary.Read(bytes.NewReader(elemData), binary.LittleEndian, elem.Addr().Interface())
		if err != nil {
			panic(err)
		}

		slice.Index(i).Set(elem)
	}

	return slice.Interface()
}
