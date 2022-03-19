package server

import (
	"errors"
	"github.com/valyala/fasthttp"
	"log"
	"net"
	"sync"
)

var tunnelConnMap = make(map[string]net.Conn)
var httpConnMap = make(map[string]*fasthttp.RequestCtx)
var subDomainToTunnelUUIDMap = make(map[string]string)
var mapMutex = sync.RWMutex{}

func AddTunnelConnection(conn net.Conn, uuid string) {
	tunnelConnMap[uuid] = conn
}
func RemoveHttpConnection(uuid string) {
	mapMutex.Lock()
	delete(httpConnMap, uuid)
	mapMutex.Unlock()
}

func GetTunnelConnection(uuid string) net.Conn {
	return tunnelConnMap[uuid]
}

func AddHttpConnection(ctx *fasthttp.RequestCtx, uuid string) {
	log.Println("Saving ctx with id ", uuid)
	mapMutex.Lock()
	httpConnMap[uuid] = ctx
	mapMutex.Unlock()
}

func GetHttpConnection(uuid string) *fasthttp.RequestCtx {
	mapMutex.Lock()
	ctx, ok := httpConnMap[uuid]
	mapMutex.Unlock()
	if !ok {
		return nil
	}
	return ctx
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
