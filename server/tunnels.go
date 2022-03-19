package server

import (
	"errors"
	"log"
	"net"
	"sync"
)

var tunnelConnMap = make(map[string]net.Conn)
var httpConnMap = make(map[string]bool)
var responseByteMap = make(map[string][]byte)
var subDomainToTunnelUUIDMap = make(map[string]string)
var mapMutex = sync.RWMutex{}

func AddTunnelConnection(conn net.Conn, uuid string) {
	tunnelConnMap[uuid] = conn
}
func RemoveHttpConnection(uuid string) {
	mapMutex.Lock()
	delete(httpConnMap, uuid)
	delete(responseByteMap, uuid)
	mapMutex.Unlock()
}

func GetTunnelConnection(uuid string) net.Conn {
	return tunnelConnMap[uuid]
}

func GetHttpConn(uuid string) bool {
	mapMutex.Lock()
	_, ok := httpConnMap[uuid]
	mapMutex.Unlock()
	return ok
}

func AddHttpConnection(uuid string) {
	log.Println("Saving ctx with id ", uuid)
	mapMutex.Lock()
	httpConnMap[uuid] = true
	mapMutex.Unlock()
}

func AddRespHttpData(uuid string, data []byte) {
	mapMutex.Lock()
	responseByteMap[uuid] = data
	mapMutex.Unlock()
}

func GetHttpConnection(uuid string) []byte {
	mapMutex.Lock()
	bytes, ok := responseByteMap[uuid]
	mapMutex.Unlock()
	if !ok {
		return nil
	}
	return bytes
}

func GetTunnelConnectionFromHostName(hostname string) (net.Conn, error) {
	if tunnelUUID, ok := subDomainToTunnelUUIDMap[hostname]; ok {
		return GetTunnelConnection(tunnelUUID), nil
	}
	return nil, errors.New("host name is not registered")
}

func AddUUIDSubDomainMap(subdomain string, uuid string) {
	subDomainToTunnelUUIDMap[subdomain] = uuid
}
