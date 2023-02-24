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
	controlManager := control_manager.ControlConnectionManager{ControlConnections: make(map[string]net.Conn)}
	//Start all the servers
	go startHttpServer(tunnelServerConfig.ServerHttpServerPort, controlManager)
	go admin.StartAdminServer(tunnelServerConfig.ServerAdminServerPort, controlManager)
	if useTLS {
		go startTLSClientTunnelServer(tunnelServerConfig, controlManager)
	} else {
		go startClientTunnelServer(tunnelServerConfig.ClientTunnelServerPort, controlManager)
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
		go handleControlConnection(conn, controlManager)
	}

}

func startTLSClientTunnelServer(config TunnelServerConfig, manager control_manager.ControlConnectionManager) {
	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(config.ClientTunnelServerPort), config.ServerTlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	workWithListener(ln)
}

func startClientTunnelServer(port int, manager control_manager.ControlConnectionManager) {
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

func startHttpServer(port int, controlManager control_manager.ControlConnectionManager) {
	httpListener, _ := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	log.Println("Starting client server")
	for {
		conn, err := httpListener.Accept()
		if err != nil {
			log.Println("Error ", err)
		}
		log.Println("Received connection from ", conn)
		go handleIncomingHttpRequest(conn, controlManager)
	}

}

func handleIncomingHttpRequest(conn net.Conn, manager control_manager.ControlConnectionManager) {
	id := uuid.New().String()
	vhostConn, err := vhost.HTTP(conn)
	if err != nil {
		log.Println("Not a valid client connection", err)
	}
	log.Println("Converted from conn to vhostConn ", vhostConn)
	controlConnection, ok := manager.GetControlConnection(vhostConn.Host())
	if !ok {
		log.Println("Control Connection not found for host=", vhostConn.Host())
		return
	}
	err = proto.SendMessage(proto.NewMessage(vhostConn.Host(), id, "init-request", []byte(id)), controlConnection)
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
			go func() {
				_, err := io.Copy(clientConn, vhostConn)
				if err != nil {
					log.Println("Failed to copy data from http to client ", err.Error())
				}
			}()
			// copy data between client connection and source connections
			_, err := io.Copy(vhostConn, clientConn)
			if err != nil {
				log.Println("Failed ", err.Error())
				tunnel_manager.RemoveTunnelConnection(id)
				return
			}
			log.Println("Copy Done")
			tunnel_manager.RemoveTunnelConnection(id)
			// close client tunnel connection
			clientConnError := clientConn.Close()
			if clientConnError != nil {
				log.Println("Failed to close Client Tunnel Connection ", clientConnError.Error())
			}
			//close http connection
			vHostConnError := vhostConn.Close()
			if vHostConnError != nil {
				log.Println("Failed to close Http Connection ", vHostConnError.Error())
			}
			//close http connection
			connError := conn.Close()
			if connError != nil {
				log.Println("Failed to close Http Connection ", connError.Error())
			}
			break
		}

	}
}

func handleControlConnection(conn net.Conn, manager control_manager.ControlConnectionManager) {
	for {
		message, err := proto.ReceiveMessage(conn)

		if err != nil && err == io.EOF {
			log.Printf("Connection closed  %s", err)
			return
		}

		if err != nil {
			log.Println("Unknown error occurred", err)
			conn.Close()
			return
		}

		if message.MessageType == "register" {
			log.Printf("Registering %s\n", message)
			manager.SaveControlConnection(message.HostName+".migtunnel.net", conn)
		}

		log.Println("Received Message ", message)
		log.Println("Received Message = " + message.MessageType + " " + message.TunnelId + " " + string(message.Data))
	}

}
