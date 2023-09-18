package lepton

import (
	gouvc "github.com/MartinRobomaze/go-uvc"
	"github.com/pkg/errors"
	"log"
)

const (
	VID = 0x1e4e
	PID = 0x0100
)

type Driver struct {
	UvcHandle    *gouvc.UVC
	DeviceHandle *gouvc.Device
	StreamHandle *gouvc.Stream
}

func New() (*Driver, error) {
	uvc := &gouvc.UVC{}

	err := uvc.Init()
	if err != nil {
		return nil, errors.Wrap(err, "error initializing uvc")
	}

	dev, err := uvc.FindDevice(VID, PID, "")
	if err != nil {
		return nil, errors.Wrap(err, "error finding device")
	}

	return &Driver{
		UvcHandle:    uvc,
		DeviceHandle: dev,
	}, nil
}

func (d *Driver) StartStreaming() (<-chan *gouvc.Frame, error) {
	err := d.DeviceHandle.Open()
	if err != nil {
		return nil, errors.Wrap(err, "error opening connection to device")
	}

	formatDesc := d.DeviceHandle.GetFormatDesc()
	frameDesc := formatDesc.FrameDescriptors()
	log.Println(formatDesc)
	stream, err := d.DeviceHandle.GetStream(
		gouvc.FRAME_FORMAT_UNCOMPRESSED,
		int(frameDesc[0].Width),
		int(frameDesc[0].Height),
		int(10000000/frameDesc[0].DefaultFrameInterval),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error getting stream handle")
	}

	if err = stream.Open(); err != nil {
		return nil, errors.Wrap(err, "error opening stream")
	}

	frameChan, err := stream.Start()
	if err != nil {
		return nil, errors.Wrap(err, "error starting stream")
	}

	return frameChan, nil
}

func (d *Driver) Close() error {
	err := d.StreamHandle.Stop()
	if err != nil {
		return errors.Wrap(err, "error stopping stream")
	}

	err = d.StreamHandle.Close()
	if err != nil {
		return errors.Wrap(err, "error closing stream")
	}

	err = d.DeviceHandle.Close()
	if err != nil {
		return errors.Wrap(err, "error closing device handle")
	}

	d.UvcHandle.Exit()

	return nil
}
