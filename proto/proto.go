package proto

type Message struct {
	HostName    string
	MessageId   string
	MessageType string
	Data        []byte
}

func NewMessage(hostName string, messageId string, messageType string, data []byte) *Message {
	ret := new(Message)
	ret.MessageId = messageId
	ret.HostName = hostName
	ret.Data = data
	ret.MessageType = messageType
	return ret
}
