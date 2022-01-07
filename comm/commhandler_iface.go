package comm

type HandlerIface interface {
	EncodeMessage(msg Message) (encodedMsg string)
	DecodeMessage(encodedMsg string) (msg Message, err error)
	WriteMessage(encodedMsg string) (response string, err error)
}
