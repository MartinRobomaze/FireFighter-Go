package comm

type MessageType string

const (
	Sensors         MessageType = "A"
	LightSensors                = "L"
	DistanceSensors             = "D"
	IMUSensor                   = "I"
	Encoders                    = "N"
	Motors                      = "M"
)

const (
	msgStart     string = "<"
	msgEnd              = ">"
	dataStart           = "{"
	dataEnd             = "}"
	commStart           = "~"
	valSeparator        = ","
)

const EncodersResetCommand = "R"

type Message struct {
	MsgType MessageType
	Data    interface{}
}
