package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/inconshreveable/go-vhost"
	"golang/admin"
	"golang/control-manager"
	"golang/proto"
	"golang/tunnel-manager"
	"io"
	"log"
	"net"
	"strconv"
)

func Start(controlServerPort int, httpServerPort int, tunnelServerPort int, adminServerPort int) {
	addr := ":" + strconv.Itoa(controlServerPort)

	go startHttpServer(httpServerPort)

	go admin.StartAdminServer(adminServerPort)

	go startClientTunnelServer(tunnelServerPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Could not create server with address=%s error=%s\n", addr, err.Error())
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Could not accept connection error=%s\n", err.Error())
			panic(err)
		}
		fmt.Printf("Received Connection from %s \n", conn.RemoteAddr())
		go handleControlConnection(conn)
	}

}

func startClientTunnelServer(port int) {
	fmt.Println("Starting Client Tunnel Server on port", port)
	httpListener, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	for {
		conn, err := httpListener.Accept()
		if err != nil {
			log.Println("Failed to accept tunnel connection ", conn, err.Error())
		} else {
			go handleClientTunnelServerConnection(conn)
		}
	}
}

func handleClientTunnelServerConnection(conn net.Conn) {
	message, err := proto.ReceiveMessage(conn)
	if err != nil {
		fmt.Println("Failed to receive message from tunnel connection", conn)
		err := conn.Close()
		if err != nil {
			return
		}
		return
	} else {
		if message.MessageType == "init-request" {
			log.Printf("Createing a new Tunnel %s\n", message)
			tunnel_manager.SaveTunnelConnection(message.TunnelId, conn)
			//go handleTunnelConnection(message, conn)
		} else {
			log.Println("Initial message from tunnel connection should be of type `init-request` found message ",
				message)
			//TODO : Check again , close the connection here?
			err := conn.Close()
			if err != nil {
				return
			}
			return
		}
	}

}

func startHttpServer(port int) {
	httpListener, _ := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	fmt.Println("Starting http server")
	for {
		conn, err := httpListener.Accept()
		if err != nil {
			fmt.Println("Error ", err)
		}
		go handleIncomingHttpRequest(conn)
	}

}

func handleIncomingHttpRequest(conn net.Conn) {
	id := uuid.New().String()
	vhostConn, err := vhost.HTTP(conn)
	if err != nil {
		fmt.Println("Not a valid http connection", err)
	}
	controlConnection, ok := control_manager.GetControlConnection(vhostConn.Host())
	if !ok {
		log.Println("Control Connection not found for host=", vhostConn.Host())
		return
	}
	err = proto.SendMessage(proto.NewMessage("localhost", id, "init-request", []byte(id)), controlConnection)
	if err != nil {
		fmt.Println("Could not send message to client connection for host", vhostConn.Host())
		return
	}
	//wait until tunnelConnections has id
	for {
		if clientConn, ok := tunnel_manager.GetTunnelConnection(id); ok {
			fmt.Println("Found Connection for tunnelId=", id)
			// new connection created between client and server
			// copy data between source connection and client connection in a new go routine
			go io.Copy(clientConn, vhostConn)
			// copy data between client connection and source connections
			_, err := io.Copy(vhostConn, clientConn)
			if err != nil {
				fmt.Println("Failed")
				return
			}
			fmt.Println("Copy Done")
			tunnel_manager.RemoveTunnelConnection(id)
			break
		}

	}
}

func handleControlConnection(conn net.Conn) {
	for {
		message, err := proto.ReceiveMessage(conn)

		if message.MessageType == "register" {
			log.Printf("Registering %s\n", message)
			control_manager.SaveControlConnection(message.HostName, conn)
		}

		if err != nil && err == io.EOF {
			log.Printf("Connection closed  %s", err)
			return
		}
		if err != nil {
			log.Println("Unknown error occurred", err)
			conn.Close()
			return
		}
		log.Println("Received Message ", message)
		log.Println("Received Message = " + message.MessageType + " " + message.TunnelId + " " + string(message.Data))
	}

}
