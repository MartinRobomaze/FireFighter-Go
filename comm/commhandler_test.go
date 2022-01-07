package comm

import (
	"fmt"
	"log"
	"strings"
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
		msgEncoded, err := h.EncodeMessage(Message{
			MsgType: Sensors,
			Data:    nil,
		})
		if err != nil {
			log.Println(err.Error())
		}

		respRaw, err := h.WriteMessage(msgEncoded)
		if err != nil {
			log.Println(err.Error())
		}

		for _, dataRaw := range strings.Split(respRaw, "\t") {
			resp, err := h.DecodeMessage(dataRaw)
			if err != nil {
				log.Println(err)
			}
			log.Printf("%+v", resp)
		}
	}
}
