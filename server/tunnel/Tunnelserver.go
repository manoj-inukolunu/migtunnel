package tunnel

import (
	"crypto/tls"
	"golang/proto"
	tunnelmanager "golang/server/tunnel-manager"
	"log"
	"net"
	"strconv"
)

type Server struct {
	Port          int
	UseTls        bool
	TlsConfig     *tls.Config
	TunnelManager tunnelmanager.TunnelManager
}

func (s *Server) Start() {
	if s.UseTls {
		go s.startTLSClientTunnelServer()
	} else {
		go s.startClientTunnelServer()
	}
}

func (s *Server) startTLSClientTunnelServer() {
	log.Println("Starting tls client tunnel server on port=", s.Port)
	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(s.Port), s.TlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	s.workWithListener(ln)
}

func (s *Server) startClientTunnelServer() {
	log.Println("Starting Non TLS Client Tunnel Server on port=", s.Port)
	httpListener, _ := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	s.workWithListener(httpListener)
}

func (s *Server) workWithListener(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept tunnel connection ", conn, err.Error())
		} else {
			go s.handleClientTunnelServerConnection(conn)
		}
	}
}

func (s *Server) handleClientTunnelServerConnection(conn net.Conn) {
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
			s.TunnelManager.SaveTunnelConnection(message.TunnelId, conn)
			// notify the corresponding http server go routine that there we have the connection
			s.TunnelManager.HttpServerChannels[message.TunnelId] <- true
			return
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
