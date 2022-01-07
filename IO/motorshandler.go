package IO

import (
	"FireFighter/comm"
	"fmt"
	"github.com/sirupsen/logrus"
)

type MotorDirection string

type MotorData struct {
	Motor     string
	Direction MotorDirection
	Speed     uint8
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
		for _, motor := range motorsData {
			data := []interface{}{motor.Motor, motor.Direction, motor.Speed}

			msg := comm.Message{
				MsgType: comm.Motors,
				Data:    data,
			}

			encodedMsg, err := m.CommHandler.EncodeMessage(msg)
			if err != nil {
				logrus.WithError(err).
					WithField("message", fmt.Sprintf("%+v", msg)).
					Log(logrus.ErrorLevel, "error encoding motors message")

				return
			}

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

			respData, ok := msgResp.Data.([]interface{})
			if !ok {
				logrus.WithError(err).
					WithField("message", response).
					Log(logrus.ErrorLevel, "error casting motors message")
				return
			}

			if !arrEqual(respData, data) {
				logrus.WithField("request", fmt.Sprintf("%+v", msg)).
					WithField("response", fmt.Sprintf("%+v", *msgResp)).
					Log(logrus.ErrorLevel, "invalid motors message")

				return
			}
		}
	}
}

func (m *MotorsHandler) SetMotors(motors []MotorData) {
	select {
	case m.motorsDataChan <- motors:
	default:
	}
}

func arrEqual(arr1 []interface{}, arr2 []interface{}) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if fmt.Sprint(arr1[i]) != fmt.Sprint(arr2[i]) {
			return false
		}
	}

	return true
}
