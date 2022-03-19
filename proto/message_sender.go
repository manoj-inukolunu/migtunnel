package proto

import (
	"encoding/json"
	"fmt"
	"net"
)

func SendMessage(mess *Message, tunnel net.Conn) error {
	enc := json.NewEncoder(tunnel)
	err := enc.Encode(&mess)
	if err != nil {
		fmt.Println("Unable to encode message ", err.Error())
	}
	return err
}

func ReceiveMessage(tunnel net.Conn) (*Message, error) {
	dec := json.NewDecoder(tunnel)
	message := &Message{}
	err := dec.Decode(message)
	if err != nil {
		fmt.Println("Unable to read message", err.Error())
	}
	return message, err
}
