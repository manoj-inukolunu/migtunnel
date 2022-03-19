package proto

import (
	"encoding/gob"
	"fmt"
	"net"
)

func SendMessage(mess *Message, tunnel net.Conn) error {
	enc := gob.NewEncoder(tunnel)
	err := enc.Encode(&mess)
	if err != nil {
		fmt.Println("Unable to encode message ", err.Error())
	}
	return err
}

func ReceiveMessage(tunnel net.Conn) (*Message, error) {
	dec := gob.NewDecoder(tunnel)
	message := &Message{}
	err := dec.Decode(message)
	if err != nil {
		fmt.Println("Unable to read message", err.Error())
	}
	return message, err
}
