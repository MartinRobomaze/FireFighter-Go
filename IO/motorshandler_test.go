package IO

import (
	"FireFighter/comm"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestMotorsHandler(t *testing.T) {
	ch, err := comm.New("/dev/ttyUSB0")
	if err != nil {
		logrus.WithError(err).Log(logrus.ErrorLevel, "error opening serial port")
	}

	h := NewMotorsHandler(ch)
	go func() {
		for {
			h.Update()
		}
	}()
	for {
		motors := []MotorData{
			{
				Motor:     "A",
				Direction: "F",
				Speed:     255,
			},
		}

		h.SetMotors(motors)
	}
}
