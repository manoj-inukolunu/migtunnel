package http

import (
	"github.com/google/uuid"
	"github.com/inconshreveable/go-vhost"
	"golang/common"
	"golang/proto"
	"golang/server/control"
	"golang/server/tunnel-manager"
	"golang/util"
	"io"
	"log"
	"net"
	"strconv"
)

type Server struct {
	Port           int
	ControlManager control.ControlManager
	TunnelManager  tunnel_manager.TunnelManager
}

func (s *Server) Start() {
	httpListener, _ := net.Listen("tcp", "localhost:"+strconv.Itoa(s.Port))
	log.Println("Starting Http Server on port=", s.Port)
	for {
		conn, err := httpListener.Accept()
		id := uuid.New().String()
		if err != nil {
			util.LogWithPrefix(id, "Error "+err.Error())
			continue
		}
		util.LogWithPrefix(id, "Received connection from "+conn.RemoteAddr().String())
		go s.handleIncomingHttpRequest(conn, id)
	}
}

func (s *Server) handleIncomingHttpRequest(conn net.Conn, id string) {
	vhostConn, err := vhost.HTTP(conn)
	if err != nil {
		util.LogWithPrefix(id, "Not a valid client connection"+err.Error())
	}
	util.LogWithPrefix(id, "Converted from conn to vhostConn "+vhostConn.RemoteAddr().String())
	s.ControlManager.SendMessage(*proto.NewMessage(vhostConn.Host(), id, "init-request", []byte(id)))
	//wait until tunnelConnections has id
	sig := make(chan bool)
	defer close(sig)
	s.TunnelManager.HttpServerChannels[id] = sig
	select {
	case <-sig:
		if clientConn, ok := s.TunnelManager.GetTunnelConnection(id); ok {
			util.LogWithPrefix(id, "Found Connection for tunnelId="+id)
			// new connection created between client and server
			// copy data between source connection and client connection in a new go routine
			//sigC := make(chan bool)
			serverToClient := common.TeeTunnel{
				Src: vhostConn,
				Dst: clientConn,
			}
			clientToServer := common.TeeTunnel{Src: clientConn, Dst: vhostConn}
			go func() {
				err := serverToClient.CopySrcToDest()
				if err != nil {
					if err != io.EOF {
						util.LogWithPrefix(id, "Failed to copy data from http to client "+err.Error())
					}
				}
				util.LogWithPrefix(id, "Finished copying from server to client")
				/*sigC <- true*/
			}()
			// copy data between client connection and source connections
			err := clientToServer.CopySrcToDest()
			util.LogWithPrefix(id, "Finished copying from client to server")
			/*<-sigC*/
			if err != nil {
				if err != io.EOF {
					util.LogWithPrefix(id, "Failed "+err.Error())
				}
				s.TunnelManager.RemoveTunnelConnection(id)
				break
			}
			util.LogWithPrefix(id, "Copy Done")
			// close client tunnel connection
			clientConnError := clientConn.Close()
			if clientConnError != nil && clientConnError != io.EOF {
				util.LogWithPrefix(id, "Failed to close Client Tunnel Connection "+clientConnError.Error())
			}
			//close http connection
			vHostConnError := vhostConn.Close()
			if vHostConnError != nil && vHostConnError != io.EOF {
				util.LogWithPrefix(id, "Failed to close vhost Http Connection "+vHostConnError.Error())
			}
			s.cleanUp(id)
			s.TunnelManager.RemoveTunnelConnection(id)
			return
		}
	}
}

func (s *Server) cleanUp(id string) {
	util.LogWithPrefix(id, "Removing channel for tunnel id = "+id)
	delete(s.TunnelManager.HttpServerChannels, id)
}
