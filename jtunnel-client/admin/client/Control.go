package client

import (
	"crypto/tls"
	"github.com/google/uuid"
	"golang/jtunnel-client/admin/tunnels"
	"golang/proto"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var controlConnections map[string]net.Conn
var tunnelsMap map[string]net.Conn

func init() {
	controlConnections = make(map[string]net.Conn)
	tunnelsMap = make(map[string]net.Conn)
}

func (client *Client) StartControlConnection() {
	log.Println("Starting Control connection")
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", uuid.New().String()+".migtunnel.net:9999", conf)
	if err != nil {
		log.Println("Failed to establish control connection ", "Error", err.Error())
		panic(err)
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	controlConnections["data"] = conn
	tunnels.SaveControlConnection(conn)
	mutex.Unlock()

	for {
		message, err := proto.ReceiveMessage(conn)
		if err != nil {
			if err.Error() == "EOF" {
				panic("Server closed control connection stopping client now")
			}
			log.Println("Error on control connection ", "Error", err.Error())
		}
		if message.MessageType == "init-request" {
			tunnel := createNewTunnel(message)
			log.Println("Created a new Tunnel")
			localConn := createLocalConnection(tunnels.GetPortForHostName(message.HostName))
			log.Println("Created Local Connection", localConn.RemoteAddr())
			go func() {
				_, err := io.Copy(localConn, tunnel)
				if err != nil {
					closeConnections(localConn, tunnel)
				}
			}()
			log.Println("Writing data to local Connection")
			_, err := io.Copy(tunnel, localConn)
			if err != nil {
				closeConnections(localConn, tunnel)
			}

			log.Println("Finished Writing data to tunnel")
			closeConnections(localConn, tunnel)
		}
		if message.MessageType == "ack-tunnel-create" {
			log.Println("Received Ack for creating tunnel from the upstream server")
			port, _ := strconv.Atoi(string(message.Data))
			tunnels.UpdateHostNameToPortMap(message.HostName, port)
		}
	}
}

func closeConnections(localConn net.Conn, tunnel net.Conn) {
	if !checkClosed(localConn) {
		localConn.Close()
	}
	if !checkClosed(tunnel) {
		tunnel.Close()
	}
}

func checkClosed(conn net.Conn) bool {
	one := make([]byte, 1)
	conn.SetReadDeadline(time.Now())
	if _, err := conn.Read(one); err == io.EOF {
		log.Println("Detected closed Local connection")
		conn.Close()
		return true
	}
	return false
}

func createNewTunnel(message *proto.Message) net.Conn {
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	conn, _ := tls.Dial("tcp", uuid.New().String()+".migtunnel.net:2121", conf)
	mutex := sync.Mutex{}
	mutex.Lock()
	tunnelsMap[message.TunnelId] = conn
	mutex.Unlock()
	proto.SendMessage(message, conn)
	return conn
}

func createLocalConnection(port int16) net.Conn {
	conn, _ := net.Dial("tcp", "localhost:"+strconv.Itoa(int(port)))
	return conn
}
