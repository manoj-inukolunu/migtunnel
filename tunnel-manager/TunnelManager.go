package tunnel_manager

import (
	"log"
	"net"
	"strings"
)

type TunnelManager struct {
	TunnelConnections  map[string]net.Conn
	HttpServerChannels map[string]chan bool
}

func (t *TunnelManager) RemoveTunnelConnection(tunnelId string) {
	if conn, ok := t.GetTunnelConnection(tunnelId); ok {
		err := conn.Close()
		if err != nil {
			log.Println("Tunnel connection is already closed for TunnelId=", tunnelId)
			delete(t.TunnelConnections, tunnelId)
			return
		}
		delete(t.TunnelConnections, tunnelId)
	}
}

func (t *TunnelManager) GetTunnelConnection(connectionId string) (net.Conn, bool) {
	conn, ok := t.TunnelConnections[connectionId]
	if ok {
		return conn.(net.Conn), ok
	}
	return nil, false
}

func (t *TunnelManager) ListAllConnectionsAsString() string {
	var s = strings.Builder{}
	for key, _ := range t.TunnelConnections {
		s.WriteString(key)
		s.WriteString("\n")
	}
	return s.String()
}

func (t *TunnelManager) SaveTunnelConnection(tunnelId string, conn net.Conn) {
	log.Println("Saving Tunnel Connection for host=", tunnelId)
	t.TunnelConnections[tunnelId] = conn
}
