package main

import (
	"fmt"
	"golang/client/admin"
	"golang/proto"
	"io"
	"net"
	"strconv"
	"sync"
)

var ControlConnections map[string]net.Conn
var tunnels map[string]net.Conn

func main() {
	ControlConnections = make(map[string]net.Conn)
	tunnels = make(map[string]net.Conn)
	fmt.Println("Starting Admin Server on ", 1234)
	go admin.StartServer(1234)
	startControlConnection()

}

func createNewTunnel(message *proto.Message) net.Conn {
	conn, _ := net.Dial("tcp", "localhost:2121")
	mutex := sync.Mutex{}
	mutex.Lock()
	tunnels[message.TunnelId] = conn
	mutex.Unlock()
	proto.SendMessage(message, conn)
	return conn
}

func createLocalConnection() net.Conn {
	conn, _ := net.Dial("tcp", "localhost:3131")
	return conn
}

func startControlConnection() {
	fmt.Println("Starting Control connection")
	conn, _ := net.Dial("tcp", "localhost:9999")
	mutex := sync.Mutex{}
	mutex.Lock()
	ControlConnections["data"] = conn
	admin.SaveControlConnection(conn)
	mutex.Unlock()

	for {
		message, _ := proto.ReceiveMessage(conn)
		fmt.Println("Received Message", message)
		if message.MessageType == "init-request" {
			tunnel := createNewTunnel(message)
			fmt.Println("Created a new Tunnel", message)
			localConn := createLocalConnection()
			fmt.Println("Created Local Connection", localConn.RemoteAddr())
			go io.Copy(localConn, tunnel)
			fmt.Println("Writing data to local Connection")
			io.Copy(tunnel, localConn)
			fmt.Println("Finished Writing data to tunnel")
			tunnel.Close()
		}
		if message.MessageType == "ack-tunnel-create" {
			fmt.Println("Received Ack for creating tunnel from the upstream server")
			port, _ := strconv.Atoi(string(message.Data))
			admin.UpdateHostNameToPortMap(message.HostName, port)
		}
	}
}
