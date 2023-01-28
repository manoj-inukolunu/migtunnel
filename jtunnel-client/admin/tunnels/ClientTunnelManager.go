package tunnels

import (
	"errors"
	"go.uber.org/zap"
	"golang/jtunnel-client/admin/data"
	"golang/proto"
	"log"
	"net"
	"sync"
)

var controlConnections = sync.Map{}

var registeredTunnels = sync.Map{}

const controlConnectionKey = "Main"

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

func GetRegisteredTunnels() sync.Map {
	return registeredTunnels
}

func CheckStatusOfLocalServers() {
	registeredTunnels.Range(func(key, value any) bool {
		return true
	})
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

func RegisterTunnel(request data.TunnelData) error {
	if conn, ok := GetControlConnection(); ok {
		proto.SendMessage(proto.NewMessage(request.HostName, "Random", "register", []byte("asdf")), conn)
		registeredTunnels.Store(request.HostName+".migtunnel.net", request.LocalServerPort)
		return nil
	}
	return errors.New("Control Connection not found , will not be creating tunnel")
}

func GetPortForHostName(hostName string) int16 {
	port, _ := registeredTunnels.Load(hostName)
	return port.(int16)
}

func UpdateHostNameToPortMap(hostName string, localServerPort int) {
	log.Printf("Request from %s are now being routed to %d\n", hostName, localServerPort)
	registeredTunnels.Store(hostName, localServerPort)
}
