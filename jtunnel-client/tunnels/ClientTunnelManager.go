package tunnels

import (
	"errors"
	"golang/jtunnel-client/data"
	"golang/proto"
	"log"
	"net"
	"sync"
)

var controlConnections = sync.Map{}

var registeredTunnels = sync.Map{}

const controlConnectionKey = "Main"

type LocalServer struct {
	ServerFqdn string
	Port       int16
	Tls        bool
}

func GetRegisteredTunnels() sync.Map {
	return registeredTunnels
}

func SaveControlConnection(conn net.Conn) {
	log.Println("Saving Control Connection")
	controlConnections.Store(controlConnectionKey, conn)
}

func GetControlConnection() (net.Conn, bool) {
	if data, ok := controlConnections.Load(controlConnectionKey); ok {
		return data.(net.Conn), true
	}
	log.Println("Control Connection not found in the map")
	return nil, false
}

func RegisterTunnel(request data.TunnelCreateRequest) error {
	if conn, ok := GetControlConnection(); ok {
		proto.SendMessage(proto.NewMessage(request.HostName, "Random", "register", []byte("asdf")), conn)
		registeredTunnels.Store(request.HostName+".migtunnel.net", LocalServer{
			ServerFqdn: request.TlsServerFQDN,
			Port:       request.LocalServerPort,
			Tls:        request.Tls,
		})
		return nil
	}
	return errors.New("Control Connection not found , will not be creating tunnel")
}

func GetLocalServer(hostName string) LocalServer {
	port, _ := registeredTunnels.Load(hostName)
	return port.(LocalServer)
}

func UpdateHostNameToPortMap(hostName string, localServerPort int) {
	log.Printf("TunnelData from %s are now being routed to %d\n", hostName, localServerPort)
	registeredTunnels.Store(hostName, localServerPort)
}
