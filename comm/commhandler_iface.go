package comm

type HandlerIface interface {
	EncodeMessage(msg Message) (encodedMsg []byte)
	DecodeMessage(encodedMsg []byte) (msg Message, err error)
	WriteMessage(encodedMsg []byte) (response []byte, err error)
}
