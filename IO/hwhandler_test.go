package IO

//
//import (
//	"FireFighter/comm"
//	"github.com/sirupsen/logrus"
//	"testing"
//)
//
//func TestHwHandler(t *testing.T) {
//	ch, err := comm.New("/dev/ttyUSB0")
//	if err != nil {
//		logrus.WithError(err).Log(logrus.ErrorLevel, "error opening serial port")
//	}
//
//	h := New(ch)
//
//	for {
//		if data := h.SensorsHandler.GetData(); data != nil {
//			logrus.Printf("%+v", data)
//		}
//		if data := h.EncodersHandler.GetData(); data != nil {
//			logrus.Printf("%+v", data)
//		}
//		h.MotorsHandler.SetMotors([]MotorData{
//			{
//				Motor:     "A",
//				Direction: "F",
//				Speed:     255,
//			},
//		})
//	}
//}
