package admin

import (
	"fmt"
	"sync"
)

var tunnelConnections = sync.Map{}

var hostNameToPortMap = sync.Map{}

func init() {
	//fmt.Println("Init of tunnels called")
}

func UpdateHostNameToPortMap(hostName string, localServerPort int) {
	fmt.Printf("Request from %s are now being routed to %d\n", hostName, localServerPort)
	hostNameToPortMap.Store(hostName, localServerPort)
}
