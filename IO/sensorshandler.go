package IO

import (
	"FireFighter/comm"
	"github.com/sirupsen/logrus"
	"strings"
)

type SensorsData struct {
	LightSensors    []int64
	DistanceSensors []int64
	IMU             float64
}

type SensorsHandler struct {
	CommHandler     *comm.Handler
	sensorsDataChan chan SensorsData
}

func NewSensorsHandler(commHandler *comm.Handler) *SensorsHandler {
	return &SensorsHandler{
		CommHandler:     commHandler,
		sensorsDataChan: make(chan SensorsData),
	}
}

func (s *SensorsHandler) Update() {
	sensorsMessage := comm.Message{
		MsgType: comm.Sensors,
		Data:    nil,
	}

	sensorsMessageEnc, err := s.CommHandler.EncodeMessage(sensorsMessage)
	if err != nil {
		logrus.WithError(err).Log(logrus.ErrorLevel, "error encoding sensors message")
		return
	}

	respRaw, err := s.CommHandler.WriteMessage(sensorsMessageEnc)
	if err != nil {
		logrus.WithError(err).Log(logrus.ErrorLevel, "error sending sensors message")
		return
	}

	//logrus.Println(respRaw)

	var sensorsData SensorsData

	for _, dataRaw := range strings.Split(respRaw, "\t") {
		resp, err := s.CommHandler.DecodeMessage(dataRaw)
		if err != nil {
			logrus.
				WithError(err).
				WithField("rawMessage", dataRaw).
				Log(logrus.ErrorLevel, "error decoding sensors message")
			return
		}

		switch resp.MsgType {
		case comm.LightSensors:
			val, ok := resp.Data.([]int64)
			if ok {
				sensorsData.LightSensors = val
			}
		case comm.DistanceSensors:
			val, ok := resp.Data.([]int64)
			if ok {
				sensorsData.DistanceSensors = val
			}
		case comm.IMUSensor:
			val, ok := resp.Data.(float64)
			if ok {
				sensorsData.IMU = val
			}
		default:
			logrus.
				WithField("MessageType", resp.MsgType).
				Log(logrus.ErrorLevel, "invalid sensors message type")
			return
		}
	}

	s.sensorsDataChan <- sensorsData
}

func (s *SensorsHandler) GetData() *SensorsData {
	select {
	case data := <-s.sensorsDataChan:
		return &data
	default:
		return nil
	}
}
