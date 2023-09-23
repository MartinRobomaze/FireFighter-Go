package IO

import (
	"FireFighter/comm"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
)

type MotorDirection string

const (
	Forward  MotorDirection = "F"
	Backward                = "B"
)

type MotorData struct {
	Motor     string
	Direction MotorDirection
	Speed     int
}

type MotorsHandler struct {
	CommHandler    *comm.Handler
	motorsDataChan chan []MotorData
}

func NewMotorsHandler(commHandler *comm.Handler) *MotorsHandler {
	return &MotorsHandler{
		CommHandler:    commHandler,
		motorsDataChan: make(chan []MotorData),
	}
}

func (m *MotorsHandler) Update() {
	select {
	case motorsData := <-m.motorsDataChan:
		sort.Slice(motorsData, func(i, j int) bool {
			return motorsData[i].Motor < motorsData[j].Motor
		})

		// 2 bytes per 1 motor.
		dataEnc := make([]byte, 8)
		var dataEncIdx int
		for _, motor := range motorsData {
			var motorSpeed int16
			if motor.Direction == Forward {
				motorSpeed = int16(motor.Speed)
			} else if motor.Direction == Backward {
				motorSpeed = int16(-motor.Speed)
			}

			dataEnc[dataEncIdx] = byte((motorSpeed >> 8) & 0xFF)
			dataEncIdx++
			dataEnc[dataEncIdx] = byte(motorSpeed & 0xFF)
			dataEncIdx++
		}

		msg := comm.Message{
			MsgType: comm.Motors,
			Data:    dataEnc,
		}

		encodedMsg := m.CommHandler.EncodeMessage(msg)

		response, err := m.CommHandler.WriteMessage(encodedMsg)
		if err != nil {
			logrus.WithError(err).
				Log(logrus.ErrorLevel, "error sending motors message")

			return
		}

		msgResp, err := m.CommHandler.DecodeMessage(response)
		if err != nil {
			logrus.WithError(err).
				WithField("message", response).
				Log(logrus.ErrorLevel, "error decoding motors message")

			return
		}

		if msgResp.MsgType != comm.Motors {
			logrus.WithField("request", fmt.Sprintf("%+v", msg)).
				WithField("response", fmt.Sprintf("%+v", *msgResp)).
				Log(logrus.ErrorLevel, "invalid motors message")

			return
		}
	default:
		return
	}
}

func (m *MotorsHandler) SetMotors(motors []MotorData) {
	m.motorsDataChan <- motors
}
