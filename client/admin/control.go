package admin

import (
	"errors"
	"golang/proto"
	"log"
	"net"
	"strconv"
	"sync"
)

// TODO : Check if we need a map here at all.
// I think one connection should be enough
var controlConnections = sync.Map{}

const controlConnectionKey = "Main"

func SaveControlConnection(conn net.Conn) {
	log.Println("Saving Control Connection")
	controlConnections.Store(controlConnectionKey, conn)
}

func CreateControlConnection(controlServer string, port int) error {
	conn, err := net.Dial("tcp", controlServer+":"+strconv.Itoa(port))
	controlConnections.Store(controlConnectionKey, conn)
	return err
}

func GetControlConnection() (net.Conn, bool) {
	if data, ok := controlConnections.Load(controlConnectionKey); ok {
		return data.(net.Conn), true
	}
	log.Println("Control Connection not found in the map")
	return nil, false
}

func CreateTunnel(request TunnelCreateRequest) error {
	if conn, ok := GetControlConnection(); ok {
		proto.SendMessage(proto.NewMessage(request.HostName, "Random", "register", []byte("asdf")), conn)
		return nil
	}
	return errors.New("Control Connection not found , will not be creating tunnel")
}
