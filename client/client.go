package main

import (
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"jtunnel-go/proto"
	"net"
)

var registered = false

func main() {

	conn, err := createControlConnection()
	if err != nil {
		fmt.Println("Unable to connect to localhost:9999", err.Error())
	}
	handleData(conn)

}

func makeRequest(data []byte) ([]byte, error) {
	server, _ := net.ResolveTCPAddr("tcp", "localhost:3131")
	client, _ := net.ResolveTCPAddr("tcp", ":")
	conn, err := net.DialTCP("tcp", client, server)
	if err != nil {
		fmt.Println("Error is ", err.Error())
		return nil, err
	} else {
		fmt.Println("Making Local Request")
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Received Error", err.Error())
			return nil, err
		}
		fmt.Println("Reading Local Response")
		result, _ := ioutil.ReadAll(conn)
		fmt.Println("Read Local Response")
		fmt.Println(string(result))
		return result, nil
	}
}

func createControlConnection() (*net.TCPConn, error) {
	server, _ := net.ResolveTCPAddr("tcp", "localhost:9999")
	client, _ := net.ResolveTCPAddr("tcp", ":")
	conn, err := net.DialTCP("tcp", client, server)
	return conn, err
}

func handleData(conn *net.TCPConn) {
	for {
		if !registered {
			err := proto.SendMessage(proto.NewMessage("localhost:8080", uuid.New().String(), "register", make([]byte, 0)), conn)
			if err != nil {
				fmt.Println("Unable to handleData ", err.Error())
			}
			registered = true
		}
		message, err := proto.ReceiveMessage(conn)
		if err != nil {
			fmt.Println("Unable to receive message", err.Error())
		}
		resp, err := makeRequest(message.Data)
		if err != nil {
			fmt.Println("Unable to make local request", err.Error())
			return
		}
		message = proto.NewMessage("localhost:8080", message.MessageId, "response", resp)
		err = proto.SendMessage(message, conn)
		if err != nil {
			fmt.Println("Unable to send response to tunnel server ", err.Error())
		}

	}

}
