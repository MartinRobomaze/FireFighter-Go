package IO

import (
	"FireFighter/comm"
	"github.com/sirupsen/logrus"
)

type SensorsData struct {
	LightSensors    []uint8
	DistanceSensors []uint16
	IMU             float64
}

type SensorsHandler struct {
	CommHandler     *comm.Handler
	sensorsDataChan chan SensorsData
}

func NewSensorsHandler(commHandler *comm.Handler) *SensorsHandler {
	return &SensorsHandler{
		CommHandler:     commHandler,
		sensorsDataChan: make(chan SensorsData, 10),
	}
}

func (s *SensorsHandler) Update() {
	sensorsMessage := comm.Message{
		MsgType: comm.Sensors,
		Data:    nil,
	}

	sensorsMessageEnc := s.CommHandler.EncodeMessage(sensorsMessage)

	respRaw, err := s.CommHandler.WriteMessage(sensorsMessageEnc)
	if err != nil {
		logrus.WithError(err).Log(logrus.ErrorLevel, "error sending sensors message")
		return
	}

	//logrus.Println(respRaw)

	var sensorsData SensorsData

	sensorsData.LightSensors = make([]uint8, 8)

	for i := 0; i < 8; i++ {
		sensorsData.LightSensors[i] = respRaw[i]
	}

	select {
	case s.sensorsDataChan <- sensorsData:
	default:
	}
}

func (s *SensorsHandler) GetData() *SensorsData {
	select {
	case data := <-s.sensorsDataChan:
		for len(s.sensorsDataChan) > 0 {
			<-s.sensorsDataChan
		}

		return &data
	default:
		return nil
	}
}
