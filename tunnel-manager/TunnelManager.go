package tunnel_manager

import (
	"log"
	"net"
	"strings"
	"sync"
)

var tunnelConnections = sync.Map{}

func init() {
	log.Println("Init for Tunnel Manager called")
}

func RemoveTunnelConnection(tunnelId string) {
	if conn, ok := GetTunnelConnection(tunnelId); ok {
		err := conn.Close()
		if err != nil {
			log.Println("Tunnel connection is already closed for TunnelId=", tunnelId)
			tunnelConnections.Delete(tunnelId)
			return
		}
		tunnelConnections.Delete(tunnelId)
	}
}

func GetTunnelConnection(connectionId string) (net.Conn, bool) {
	conn, ok := tunnelConnections.Load(connectionId)
	if ok {
		return conn.(net.Conn), ok
	}
	return nil, false
}

func ListAllConnectionsAsString() string {
	var s = strings.Builder{}
	tunnelConnections.Range(func(key, value any) bool {
		s.WriteString(key.(string))
		s.WriteString("\n")
		return true
	})
	return s.String()
}

func SaveTunnelConnection(tunnelId string, conn net.Conn) {
	log.Println("Saving Control Connection for host=", tunnelId)
	tunnelConnections.Store(tunnelId, conn)
}
