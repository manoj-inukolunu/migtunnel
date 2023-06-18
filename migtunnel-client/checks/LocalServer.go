package checks

import (
	"golang/migtunnel-client/tunnels"
	"log"
	"net"
	"time"
)

func CheckLocalServerPorts() {
	registeredTunnels := tunnels.GetRegisteredTunnels()

	log.Println("Admin Server running on %s ")

	registeredTunnels.Range(func(key, value any) bool {
		if rawConnect("localhost", value.(string)) {
			log.Println("SUCCESS Connection to server on port %s successful", value)
		} else {
			log.Println("FAIL could not connect to server on port %s", value)
		}

		return true
	})

}

func rawConnect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		log.Println("Could not connect to port %s %s", port, err.Error())
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return true
}
