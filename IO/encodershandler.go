package IO

import (
	"FireFighter/comm"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)

type EncodersHandler struct {
	CommHandler      *comm.Handler
	encodersDataChan chan []float64
	encoderResetChan chan bool
}

func NewEncodersHandler(commHandler *comm.Handler) *EncodersHandler {
	return &EncodersHandler{
		CommHandler:      commHandler,
		encodersDataChan: make(chan []float64),
		encoderResetChan: make(chan bool),
	}
}

func (e *EncodersHandler) Update() {
	select {
	case <-e.encoderResetChan:
		msg := comm.Message{
			MsgType: comm.Encoders,
			Data:    comm.EncodersResetCommand,
		}

		msgEncoded, err := e.CommHandler.EncodeMessage(msg)
		if err != nil {
			logrus.WithError(err).
				WithField("message", fmt.Sprintf("%+v", msg)).
				Log(logrus.ErrorLevel, "error encoding encoders message")

			return
		}

		respRaw, err := e.CommHandler.WriteMessage(msgEncoded)
		if err != nil {
			logrus.WithError(err).
				WithField("message", msgEncoded).
				Log(logrus.ErrorLevel, "error sending encoders message")

			return
		}

		msgResp, err := e.CommHandler.DecodeMessage(respRaw)
		if err != nil {
			logrus.WithError(err).
				WithField("response", respRaw).
				Log(logrus.ErrorLevel, "error decoding encoders response")

			return
		}

		if !reflect.DeepEqual(msgEncoded, msg) {
			logrus.WithField("request", fmt.Sprintf("%+v", msg)).
				WithField("response", fmt.Sprintf("%+v", msgResp)).
				Log(logrus.ErrorLevel, "invalid encoders message")

			return
		}
	default:
		msg := comm.Message{
			MsgType: comm.Encoders,
			Data:    nil,
		}

		msgEncoded, err := e.CommHandler.EncodeMessage(msg)
		if err != nil {
			logrus.WithError(err).
				WithField("message", fmt.Sprintf("%+v", msg)).
				Log(logrus.ErrorLevel, "error encoding encoders message")

			return
		}

		respRaw, err := e.CommHandler.WriteMessage(msgEncoded)
		if err != nil {
			logrus.WithError(err).
				WithField("message", msgEncoded).
				Log(logrus.ErrorLevel, "error sending encoders message")

			return
		}

		msgResp, err := e.CommHandler.DecodeMessage(respRaw)
		if err != nil {
			logrus.WithError(err).
				WithField("response", respRaw).
				Log(logrus.ErrorLevel, "error decoding encoders response")

			return
		}

		data, ok := msgResp.Data.([]float64)
		if ok {
			e.encodersDataChan <- data
		} else {
			logrus.WithField("response", fmt.Sprintf("%+v", msgResp)).
				Log(logrus.ErrorLevel, "error casting encoders data")

			return
		}
	}
}

func (e *EncodersHandler) GetData() []float64 {
	select {
	case data := <-e.encodersDataChan:
		return data
	default:
		return nil
	}
}
