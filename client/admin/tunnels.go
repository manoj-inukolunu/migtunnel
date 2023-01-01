package admin

import (
	"log"
	"sync"
)

var tunnelConnections = sync.Map{}

var hostNameToPortMap = sync.Map{}

func init() {
	//log.Println("Init of tunnels called")
}

func UpdateHostNameToPortMap(hostName string, localServerPort int) {
	log.Printf("Request from %s are now being routed to %d\n", hostName, localServerPort)
	hostNameToPortMap.Store(hostName, localServerPort)
}
