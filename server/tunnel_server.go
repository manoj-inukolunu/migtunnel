package server

import (
	"crypto/tls"
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

type TunnelServerConfig struct {
	ClientControlServerPort int
	ServerHttpServerPort    int
	ClientTunnelServerPort  int
	ServerAdminServerPort   int
	ServerTlsConfig         *tls.Config
}

func Start(tunnelServerConfig TunnelServerConfig) {
	useTLS := (tunnelServerConfig.ServerTlsConfig != nil)

	//Start all the servers
	go startHttpServer(tunnelServerConfig.ServerHttpServerPort)
	go admin.StartAdminServer(tunnelServerConfig.ServerAdminServerPort)
	if useTLS {
		go startTLSClientTunnelServer(tunnelServerConfig)
	} else {
		go startClientTunnelServer(tunnelServerConfig.ClientTunnelServerPort)
	}

	var listener net.Listener
	var err error
	addr := ":" + strconv.Itoa(tunnelServerConfig.ClientControlServerPort)
	if useTLS {
		log.Println("Using TLS to create client control server")
		listener, err = tls.Listen("tcp", addr, tunnelServerConfig.ServerTlsConfig)
	} else {
		listener, err = net.Listen("tcp", addr)
	}

	if err != nil {
		log.Printf("Could not create server with address=%s error=%s\n", addr, err.Error())
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Could not accept connection error=%s\n", err.Error())
		}
		log.Printf("Received Connection from %s \n", conn.RemoteAddr())
		go handleControlConnection(conn)
	}

}

func startTLSClientTunnelServer(config TunnelServerConfig) {
	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(config.ClientTunnelServerPort), config.ServerTlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	workWithListener(ln)
}

func startClientTunnelServer(port int) {
	log.Println("Starting Client Tunnel Server on port", port)
	httpListener, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
	workWithListener(httpListener)
}

func workWithListener(httpListener net.Listener) {
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
	defer conn.Close()
	message, err := proto.ReceiveMessage(conn)
	if err != nil {
		log.Println("Failed to receive message from tunnel connection", conn)
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
	log.Println("Starting http server")
	defer httpListener.Close()
	for {
		conn, err := httpListener.Accept()
		if err != nil {
			log.Println("Error ", err)
		}
		go handleIncomingHttpRequest(conn)
	}

}

func handleIncomingHttpRequest(conn net.Conn) {
	defer conn.Close()
	id := uuid.New().String()
	vhostConn, err := vhost.HTTP(conn)
	if err != nil {
		log.Println("Not a valid http connection", err)
	}
	controlConnection, ok := control_manager.GetControlConnection(vhostConn.Host())
	if !ok {
		log.Println("Control Connection not found for host=", vhostConn.Host())
		return
	}
	err = proto.SendMessage(proto.NewMessage("localhost", id, "init-request", []byte(id)), controlConnection)
	if err != nil {
		log.Println("Could not send message to client connection for host", vhostConn.Host())
		return
	}
	//wait until tunnelConnections has id
	for {
		if clientConn, ok := tunnel_manager.GetTunnelConnection(id); ok {
			log.Println("Found Connection for tunnelId=", id)
			// new connection created between client and server
			// copy data between source connection and client connection in a new go routine
			go io.Copy(clientConn, vhostConn)
			// copy data between client connection and source connections
			_, err := io.Copy(vhostConn, clientConn)
			if err != nil {
				log.Println("Failed")
				return
			}
			log.Println("Copy Done")
			tunnel_manager.RemoveTunnelConnection(id)
			break
		}

	}
}

func handleControlConnection(conn net.Conn) {
	defer conn.Close()
	for {
		message, err := proto.ReceiveMessage(conn)

		if err != nil && err == io.EOF {
			log.Printf("Connection closed  %s", err)
			return
		}

		if err != nil {
			log.Println("Unknown error occurred", err)
			return
		}

		if message.MessageType == "register" {
			log.Printf("Registering %s\n", message)
			control_manager.SaveControlConnection(message.HostName+".migtunnel.net", conn)
		}

		log.Println("Received Message ", message)
		log.Println("Received Message = " + message.MessageType + " " + message.TunnelId + " " + string(message.Data))
	}

}
