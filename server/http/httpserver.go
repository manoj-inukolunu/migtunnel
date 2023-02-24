package http

import (
	"github.com/google/uuid"
	"github.com/inconshreveable/go-vhost"
	"golang/proto"
	"golang/server/control"
	tunnel_manager "golang/tunnel-manager"
	"io"
	"log"
	"net"
	"strconv"
)

type Server struct {
	// This is the communication channel between tunnel server and the http server
	TunnelChannel  chan string
	Port           int
	ControlManager control.Server
}

func (s *Server) Start() {
	httpListener, _ := net.Listen("tcp", "localhost:"+strconv.Itoa(s.Port))
	log.Println("Starting client server")
	for {
		conn, err := httpListener.Accept()
		if err != nil {
			log.Println("Error ", err)
		}
		log.Println("Received connection from ", conn)
		go s.handleIncomingHttpRequest(conn)
	}
}

func (s *Server) handleIncomingHttpRequest(conn net.Conn) {
	id := uuid.New().String()
	vhostConn, err := vhost.HTTP(conn)
	if err != nil {
		log.Println("Not a valid client connection", err)
	}
	log.Println("Converted from conn to vhostConn ", vhostConn)
	s.ControlManager.SendMessage(*proto.NewMessage(vhostConn.Host(), id, "init-request", []byte(id)))
	//wait until tunnelConnections has id
	select {
	case <-s.TunnelChannel:
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
