package IO

import (
	"FireFighter/comm"
	"fmt"
	"github.com/sirupsen/logrus"
)

type EncodersHandler struct {
	CommHandler      *comm.Handler
	encodersDataChan chan []float64
	encoderResetChan chan bool
}

func NewEncodersHandler(commHandler *comm.Handler) *EncodersHandler {
	return &EncodersHandler{
		CommHandler:      commHandler,
		encodersDataChan: make(chan []float64, 10),
		encoderResetChan: make(chan bool),
	}
}

func (e *EncodersHandler) Update() {
	select {
	case <-e.encoderResetChan:
		msg := comm.Message{
			MsgType: comm.EncodersReset,
			Data:    nil,
		}

		msgEncoded := e.CommHandler.EncodeMessage(msg)

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

		if msgResp.MsgType != comm.EncodersReset {
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

		msgEncoded := e.CommHandler.EncodeMessage(msg)

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

		encodersData := make([]float64, 4)
		var respRawIdx int

		for i := 0; i < 4; i++ {
			var encoderDataDeg int32
			// IMU data is encoded in 4 byts (32 bits)
			for j := 3; j >= 0; j-- {
				encoderDataDeg = encoderDataDeg | (int32(msgResp.Data[respRawIdx]) << (j * 8))
				respRawIdx++
			}

			encodersData[i] = float64(encoderDataDeg) / 360
		}

		select {
		case e.encodersDataChan <- encodersData:
		default:
		}
	}
}

func (e *EncodersHandler) GetData() []float64 {
	select {
	case data := <-e.encodersDataChan:
		for len(e.encodersDataChan) > 0 {
			<-e.encodersDataChan
		}

		return data
	default:
		return nil
	}
}
