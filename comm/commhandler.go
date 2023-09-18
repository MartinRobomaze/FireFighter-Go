package comm

import (
	"bufio"
	"github.com/pkg/errors"
	"go.bug.st/serial"
	"time"
)

type Handler struct {
	Port        serial.Port
	portName    string
	portOptions *serial.Mode
	reader      *bufio.Scanner
}

func New(portName string) (*Handler, error) {
	port, err := serial.Open(portName, &serial.Mode{BaudRate: 1000000})
	if err != nil {
		return nil, errors.Wrap(err, "error creating CommHandler")
	}

	return &Handler{
		Port:        port,
		portName:    portName,
		portOptions: &serial.Mode{BaudRate: 1000000},
		reader:      bufio.NewScanner(port),
	}, nil
}

func (h *Handler) EncodeMessage(msg Message) (encodedMsg []byte) {
	encodedMsg = make([]byte, len(msg.Data)+4)

	encodedMsg[0] = msgStart
	encodedMsg[1] = byte(msg.MsgType)

	copy(encodedMsg[2:], msg.Data)

	encodedMsg[len(encodedMsg)-2] = msgEnd
	encodedMsg[len(encodedMsg)-1] = '\n'

	return encodedMsg
}

func (h *Handler) DecodeMessage(encodedMsg []byte) (msg *Message, err error) {
	if encodedMsg[0] != msgStart || encodedMsg[len(encodedMsg)-2] != msgEnd {
		return nil, errors.New("invalid message")
	}

	msg = &Message{
		MsgType: MessageType(encodedMsg[1]),
		Data:    encodedMsg[2 : len(encodedMsg)-3],
	}

	return msg, nil
}

func (h *Handler) WriteMessage(encodedMsg []byte) (response []byte, err error) {
	err = h.Port.ResetOutputBuffer()
	if err != nil {
		return nil, err
	}

	err = h.Port.ResetInputBuffer()
	if err != nil {
		return nil, err
	}

	written, err := h.Port.Write(encodedMsg)
	if err != nil {
		return nil, errors.Wrap(err, "error writing message")
	}

	if written != len(encodedMsg) {
		return nil, errors.Errorf("Expected to write %d bytes, wrote %d", len(encodedMsg), written)
	}

	response, err = h.readln(50 * time.Millisecond)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response")
	}

	return response[1:], nil
}

func (h *Handler) Close() error {
	return h.Port.Close()
}

func (h *Handler) readln(timeout time.Duration) ([]byte, error) {
	dataChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		reader := bufio.NewReader(h.Port)

		line, err := reader.ReadBytes('\n')
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
		return nil, err
	case <-time.After(timeout):
		err := h.Port.Close()
		if err != nil {
			return nil, errors.Wrap(err, "Error closing port on timeout")
		}

		h.Port, err = serial.Open(h.portName, h.portOptions)
		if err != nil {
			return nil, errors.Wrap(err, "Error opening port on timeout")
		}

		return nil, errors.New("Timeout")
	}
}
