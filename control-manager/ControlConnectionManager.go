package control_manager

import (
	"fmt"
	"golang/proto"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var controlConnections = sync.Map{}

func init() {
	fmt.Println("Init for ControlConnectionManager called")
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Checking for closed control connections")
				controlConnections.Range(func(hostName, connection any) bool {
					conn := connection.(net.Conn)
					err := proto.SendMessage(proto.PingMessage(), conn)
					if err != nil {
						fmt.Printf("Could not send ping to hostName=%s ,error=%s, will be closing connection\n", hostName, err)
						controlConnections.Delete(hostName)
					}
					return true
				})
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func ListAllConnectionsAsString() string {
	var s = strings.Builder{}
	controlConnections.Range(func(key, value any) bool {
		s.WriteString(key.(string))
		s.WriteString("\n")
		return true
	})
	return s.String()
}

func GetControlConnection(hostName string) (net.Conn, bool) {
	conn, ok := controlConnections.Load(hostName)
	if ok {
		return conn.(net.Conn), ok
	}
	return nil, false
}

func SaveControlConnection(hostName string, conn net.Conn) {
	log.Println("Saving Control Connection for host=", hostName)
	controlConnections.Store(hostName, conn)

}
