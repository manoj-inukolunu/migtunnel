package control

import (
	"crypto/tls"
	"errors"
	"golang/proto"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ControlManager struct {
	ControlConnections map[string]net.Conn
	ControlServerPort  int
	ServerTlsConfig    *tls.Config
	UseTLS             bool
}

func (a *ControlManager) CheckConnections() {
	log.Println("Init for ControlServer called")
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				//log.Println("Checking for closed control connections")
				for hostName, conn := range a.ControlConnections {
					err := proto.SendMessage(proto.PingMessage(), conn)
					if err != nil {
						log.Printf("Could not send ping to hostName=%s ,error=%s, will be closing connection\n", hostName, err)
						delete(a.ControlConnections, hostName)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func (c *ControlManager) Start() {
	log.Println("Starting Control  server on port=", c.ControlServerPort)
	var listener net.Listener
	var err error
	addr := ":" + strconv.Itoa(c.ControlServerPort)
	if c.UseTLS {
		log.Println("Using TLS to create client control server")
		listener, err = tls.Listen("tcp", addr, c.ServerTlsConfig)
	} else {
		listener, err = net.Listen("tcp", addr)
	}

	if err != nil {
		log.Printf("Could not create server with address=%s error=%s\n", addr, err.Error())
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Could not accept connection error=%s\n", err.Error())
		}
		log.Printf("Received Connection from %s \n", conn.RemoteAddr())
		go c.handleControlConnection(conn)
	}
}

func (c *ControlManager) handleControlConnection(conn net.Conn) {
	for {
		message, err := proto.ReceiveMessage(conn)

		if err != nil && err == io.EOF {
			log.Printf("Connection closed  %s", err)
			return
		}

		if err != nil {
			log.Println("Unknown error occurred", err)
			conn.Close()
			return
		}

		if message.MessageType == "register" {
			log.Printf("Registering %s\n", message)
			c.saveControlConnection(message.HostName+".migtunnel.net", conn)
		}

		log.Println("Received Message ", message)
		log.Println("Received Message = " + message.MessageType + " " + message.TunnelId + " " + string(message.Data))
	}

}

func (a *ControlManager) InitCronitorHeartbeat() {
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

func (a *ControlManager) ListAllConnectionsAsString() string {
	var s = strings.Builder{}
	for hostName, _ := range a.ControlConnections {
		s.WriteString(hostName)
		s.WriteString("\n")
	}
	return s.String()
}

func (a *ControlManager) GetControlConnection(hostName string) (net.Conn, bool) {
	conn, ok := a.ControlConnections[(hostName)]
	if ok {
		return conn.(net.Conn), ok
	}
	return nil, false
}

func (a *ControlManager) SendMessage(message proto.Message) error {
	controlConnection, ok := a.GetControlConnection(message.HostName)
	if !ok {
		log.Println("Control Connection not found for host=", message.HostName)
		return errors.New("Control Connection not found for host=" + message.HostName)
	}
	log.Println("Found control connection ")
	err := proto.SendMessage(&message, controlConnection)
	if err != nil {
		log.Println("Could not send message to client connection for host", message.HostName)
		return errors.New("Could not send message to client connection for host = " + message.HostName)
	}
	return nil
}

func (a *ControlManager) saveControlConnection(hostName string, conn net.Conn) {
	log.Println("Saving Control Connection for host=", hostName)
	a.ControlConnections[hostName] = conn

}
