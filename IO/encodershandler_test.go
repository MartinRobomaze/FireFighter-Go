package IO

import (
	"FireFighter/comm"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestEncodersHandler(t *testing.T) {
	ch, err := comm.New("/dev/ttyUSB0")
	if err != nil {
		logrus.WithError(err).Log(logrus.ErrorLevel, "error opening serial port")
	}

	h := NewEncodersHandler(ch)
	go func() {
		for {
			h.Update()
		}
	}()
	for {
		if data := h.GetData(); data != nil {
			logrus.Printf("%+v", data)
		}
	}
}
