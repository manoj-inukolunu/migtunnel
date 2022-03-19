package server

import (
	"fmt"
	"github.com/google/uuid"
	"jtunnel-go/proto"
	"log"
	"net"
)

func TunnelServer(tunnelServerPort int) {
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("Failed to start listener on 9999", err.Error())
		return
	}
	fmt.Println("Starting listener on 9999")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection ", err.Error())
			return
		}
		go handleTunnelConnection(conn)
	}
}

func handleTunnelConnection(conn net.Conn) {

	for {
		message, err := proto.ReceiveMessage(conn)
		if err != nil {
			fmt.Println("Unable to read message from client", err.Error())
			break
		}

		if message.MessageType == "register" {
			fmt.Println("Register request Received")
			id := uuid.New().String()
			AddUUIDSubDomainMap(message.HostName, id)
			AddTunnelConnection(conn, id)
		} else {
			//fmt.Println("Message Received From Client is ", message)
			ctx := GetHttpConn(message.MessageId)
			if ctx {
				AddRespHttpData(message.MessageId, message.Data)
			} else {
				log.Println("No ctx Found to write to")
			}

		}

	}
}
