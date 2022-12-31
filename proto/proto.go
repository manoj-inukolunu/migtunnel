package proto

type Message struct {
	HostName    string
	TunnelId    string
	MessageType string
	Data        []byte
}

func PingMessage() *Message {
	return NewMessage("ping", "ping", "ping", []byte("ping"))
}

func NewMessage(hostName string, tunnelId string, messageType string, data []byte) *Message {
	ret := new(Message)
	ret.TunnelId = tunnelId
	ret.HostName = hostName
	ret.Data = data
	ret.MessageType = messageType
	return ret
}
