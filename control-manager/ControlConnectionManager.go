package control_manager

import (
	"golang/proto"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var controlConnections map[string]net.Conn

func init() {
	controlConnections := make(map[string]net.Conn)
	log.Println("Init for ControlConnectionManager called")
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Checking for closed control connections")
				for hostName, conn := range controlConnections {
					err := proto.SendMessage(proto.PingMessage(), conn)
					if err != nil {
						log.Printf("Could not send ping to hostName=%s ,error=%s, will be closing connection\n", hostName, err)
						delete(controlConnections, hostName)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	initCronitorHeartbeat()
}

func initCronitorHeartbeat() {
	log.Println("Starting Cronitor Heartbeat , ticker at 60 seconds")
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, err := http.Get("https://cronitor.link/p/70f15f445cc1490ba2183404be52079c/UQZCNL")
				if err != nil {
					log.Println("Failed to send heartbeat", err.Error())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func ListAllConnectionsAsString() string {
	var s = strings.Builder{}
	for hostName, _ := range controlConnections {
		s.WriteString(hostName)
		s.WriteString("\n")
	}

	return s.String()
}

func GetControlConnection(hostName string) (net.Conn, bool) {
	conn, ok := controlConnections[(hostName)]
	if ok {
		return conn.(net.Conn), ok
	}
	return nil, false
}

func SaveControlConnection(hostName string, conn net.Conn) {
	log.Println("Saving Control Connection for host=", hostName)
	controlConnections[hostName] = conn

}
