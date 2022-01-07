package comm

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"go.bug.st/serial"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Port        serial.Port
	portName    string
	portOptions *serial.Mode
	reader      *bufio.Scanner
}

func New(portName string) (*Handler, error) {
	port, err := serial.Open(portName, &serial.Mode{BaudRate: 230400})
	if err != nil {
		return nil, errors.Wrap(err, "error creating CommHandler")
	}

	return &Handler{
		Port:        port,
		portName:    portName,
		portOptions: &serial.Mode{BaudRate: 230400},
		reader:      bufio.NewScanner(port),
	}, nil
}

func (h *Handler) EncodeMessage(msg Message) (encodedMsg string, err error) {
	var dataStr string

	if msg.Data != nil {
		data, ok := msg.Data.([]interface{})
		if !ok {
			return "", errors.New("error converting data type")
		}

		for _, val := range data {
			dataStr += fmt.Sprint(val) + ","
		}
		dataStr = dataStr[:len(dataStr)-1]
	}

	encodedMsg = msgStart + string(msg.MsgType) + dataStart + dataStr + dataEnd + msgEnd

	return encodedMsg, nil
}

func (h *Handler) DecodeMessage(encodedMsg string) (msg *Message, err error) {
	encodedMsg = strings.TrimSuffix(encodedMsg, "\n")

	if strings.HasPrefix(encodedMsg, msgStart) && strings.HasSuffix(encodedMsg, msgEnd) {
		dataRaw := encodedMsg[strings.Index(encodedMsg, dataStart)+1 : strings.Index(encodedMsg, dataEnd)]

		dataStr := strings.Split(dataRaw, valSeparator)

		switch MessageType(encodedMsg[1]) {
		case LightSensors:
			data := make([]int64, len(dataStr))

			for i, strVal := range dataStr {
				data[i], err = strconv.ParseInt(strVal, 10, 64)
				if err != nil {
					return nil, errors.Wrap(err, "error decoding message")
				}
			}

			return &Message{
				MsgType: LightSensors,
				Data:    data,
			}, nil
		case DistanceSensors:
			data := make([]int64, len(dataStr))

			for i, strVal := range dataStr {
				data[i], err = strconv.ParseInt(strVal, 10, 64)
				if err != nil {
					return nil, errors.Wrap(err, "error decoding message")
				}
			}

			return &Message{
				MsgType: DistanceSensors,
				Data:    data,
			}, nil
		case IMUSensor:
			data, err := strconv.ParseFloat(dataStr[0], 64)
			if err != nil {
				return nil, errors.Wrap(err, "error decoding message")
			}

			return &Message{
				MsgType: IMUSensor,
				Data:    data,
			}, nil
		case Motors:
			data := make([]interface{}, len(dataStr))

			data[0] = dataStr[0]
			data[1] = dataStr[1]

			speed, err := strconv.ParseInt(dataStr[2], 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "error decoding message")
			}

			data[2] = speed

			return &Message{
				MsgType: Motors,
				Data:    data,
			}, nil
		case Encoders:
			if len(dataStr) == 1 {
				data := dataStr[0]

				return &Message{
					MsgType: Encoders,
					Data:    data,
				}, nil
			}

			data := make([]float64, len(dataStr))

			for i, strVal := range dataStr {
				data[i], err = strconv.ParseFloat(strVal, 64)
				if err != nil {
					return nil, errors.Wrap(err, "error decoding message")
				}
			}

			return &Message{
				MsgType: Encoders,
				Data:    data,
			}, nil
		default:
			return nil, errors.Errorf("Invalid message type %s", string(encodedMsg[1]))
		}
	}

	return nil, errors.New("Invalid message")
}

func (h *Handler) WriteMessage(encodedMsg string) (response string, err error) {
	//h.Port.ResetOutputBuffer()
	//h.Port.ResetInputBuffer()

	written, err := h.Port.Write([]byte(encodedMsg))
	if err != nil {
		return "", errors.Wrap(err, "error writing message")
	}

	if written != len(encodedMsg) {
		return "", errors.Errorf("Expected to write %d bytes, wrote %d", len(encodedMsg), written)
	}

	response, err = h.readln(50 * time.Millisecond)
	if err != nil {
		return "", errors.Wrap(err, "error reading response")
	}

	if strings.HasPrefix(response, commStart) {
		return response[1:], nil
	} else {
		return "", errors.Errorf("Invalid message %s", response)
	}
}

func (h *Handler) readln(timeout time.Duration) (string, error) {
	dataChan := make(chan string)
	errChan := make(chan error)

	go func() {
		reader := bufio.NewReader(h.Port)

		line, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
		} else {
			dataChan <- line
		}

		close(dataChan)
		close(errChan)
	}()

	select {
	case line := <-dataChan:
		return line, nil
	case err := <-errChan:
		return "", err
	case <-time.After(timeout):
		err := h.Port.Close()
		if err != nil {
			return "", errors.Wrap(err, "Error closing port on timeout")
		}

		h.Port, err = serial.Open(h.portName, h.portOptions)
		if err != nil {
			return "", errors.Wrap(err, "Error opening port on timeout")
		}

		return "", errors.New("Timeout")
	}
}
