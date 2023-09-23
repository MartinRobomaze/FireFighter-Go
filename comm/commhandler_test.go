package comm

import (
	"fmt"
	"log"
	"testing"
)

func TestPug(t *testing.T) {
	h, err := New("/dev/ttyUSB0")
	if err != nil {
		log.Fatal(err.Error())
	}
	i := 0
	timeouts := 0

	defer func() {
		fmt.Println(i, timeouts)
	}()

	for {
		msgEncoded := h.EncodeMessage(Message{
			MsgType: Sensors,
			Data:    nil,
		})
		if err != nil {
			log.Println(err.Error())
			continue
		}

		respRaw, err := h.WriteMessage(msgEncoded)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println(len(respRaw))

		msg, err := h.DecodeMessage(respRaw)
		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println(len(msg.Data))
		}
	}
}
