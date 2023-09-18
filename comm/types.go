package comm

type MessageType byte

const (
	Sensors       MessageType = 'S'
	Encoders                  = 'E'
	EncodersReset             = 'R'
	Motors                    = 'M'
	BrakeMotors               = 'B'
)

const (
	msgStart byte = '<'
	msgEnd        = '>'
)

type Message struct {
	MsgType MessageType
	Data    []byte
}
