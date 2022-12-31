package main

import (
	"crypto/tls"
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
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	conn, _ := tls.Dial("tcp", "manoj.lc-algorithms.com:2121", conf)
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
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "manoj.lc-algorithms.com:9999", conf)
	if err != nil {
		fmt.Println("Failed to establish control connection ", err.Error())
		return
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	ControlConnections["data"] = conn
	admin.SaveControlConnection(conn)
	mutex.Unlock()

	for {
		message, err := proto.ReceiveMessage(conn)
		fmt.Println("Received Message", message)
		if err != nil {
			if err.Error() == "EOF" {
				panic("Server closed control connection stopping client now")
			}
			fmt.Println("Error on control connection ", err.Error())
		}
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
