package IO

import (
	"FireFighter/comm"
	"github.com/sirupsen/logrus"
)

type SensorsData struct {
	LightSensors    []uint8
	DistanceSensors []uint16
	IMU             int32
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

	msgDecoded, err := s.CommHandler.DecodeMessage(respRaw)
	if err != nil {
		logrus.WithError(err).Log(logrus.ErrorLevel, "error decoding sensors message")
		return
	}

	//logrus.Println(respRaw)

	var sensorsData SensorsData
	var respRawIdx int

	// Data from each light sensor is in 1 byte(8 bits).
	sensorsData.LightSensors = make([]uint8, 8)
	for i := 0; i < 8; i++ {
		sensorsData.LightSensors[i] = msgDecoded.Data[respRawIdx]
		respRawIdx++
	}

	// Data from each distance sensor is in 2 bytes(16 bits).
	sensorsData.DistanceSensors = make([]uint16, 8)
	for i := 0; i < 8; i++ {
		sensorsData.DistanceSensors[i] = uint16(msgDecoded.Data[respRawIdx])<<8 | uint16(msgDecoded.Data[respRawIdx+1])
		respRawIdx += 2
	}

	// IMU data is encoded in 4 byts (32 bits)
	for i := 3; i >= 0; i-- {
		sensorsData.IMU = sensorsData.IMU | (int32(msgDecoded.Data[respRawIdx]) << (i * 8))
		respRawIdx++
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
