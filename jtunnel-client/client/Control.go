package client

import (
	"crypto/tls"
	"golang/jtunnel-client/db"
	"golang/jtunnel-client/tunnels"
	"golang/jtunnel-client/util"
	"golang/proto"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var controlConnections map[string]net.Conn
var tunnelsMap map[string]net.Conn

const (
	TunnelPort  = 2121
	ControlPort = 9999
)

func init() {
	controlConnections = make(map[string]net.Conn)
	tunnelsMap = make(map[string]net.Conn)
}

func (client *Client) StartControlConnection(localDb db.LocalDb, isLocal bool) {
	log.Println("Starting Control connection")
	conn, err := getControlConnection(isLocal)
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
		log.Println("Here")
		message, err := proto.ReceiveMessage(conn)
		if err != nil {
			if err.Error() == "EOF" {
				panic("Server closed control connection stopping client now")
			}
			log.Println("Error on control connection ", "Error", err.Error())
		}
		if message.MessageType == "init-request" {
			go HandleIncomingRequest(message, isLocal)
		}
		if message.MessageType == "ack-tunnel-create" {
			log.Println("Received Ack for creating tunnel from the upstream server")
			port, _ := strconv.Atoi(string(message.Data))
			tunnels.UpdateHostNameToPortMap(message.HostName, port)
		}
	}
}

func HandleIncomingRequest(message *proto.Message, isLocal bool) {
	tunnel := createNewTunnel(message, isLocal)
	log.Println("Created a new TunnelPort")
	localConn, localConnErr := createLocalConnection(tunnels.GetLocalServer(message.HostName))
	if localConnErr != nil {
		log.Printf("Could not connect to local server on port %d "+
			"Please check if server is running.\n", tunnels.GetLocalServer(message.HostName).Port)
		return
	}
	log.Println("Created Local Connection", localConn.RemoteAddr())
	/*tunnelProcessor := util.NewTeeReader(message.TunnelId, tunnel, localConn, localDb, false,
	tunnels.GetLocalServer(message.HostName))*/
	sig := make(chan bool)
	go func() {
		log.Println("Reading data form tunnel")
		//err := tunnelProcessor.TunnelToLocal()
		_, err := io.Copy(localConn, tunnel)
		if err != nil && !strings.Contains(err.Error(), "use of closed") {
			log.Println("Error reading from tunnel ", err.Error())
		}
		log.Println("Finished writing data from tunnel to local")
		//sig <- true
	}()
	//err := tunnelProcessor.LocalToTunnel()
	log.Println("Copying from Local to Tunnel")
	_, err := io.Copy(tunnel, localConn)
	err = localConn.Close()
	tunnel.Close()
	if err != nil && !strings.Contains(err.Error(), "use of closed") {
		log.Println("Error writing to tunnel ", err.Error())
	}
	log.Println("Finished Writing data to tunnel")
	//<-sig
	close(sig)
	log.Println("All Done")
	closeConnections(localConn, tunnel)
}

func getControlConnection(isLocal bool) (net.Conn, error) {
	if isLocal {
		conn, err := net.Dial("tcp", util.GetRemoteUrl(isLocal, ControlPort))
		return conn, err
	} else {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := tls.Dial("tcp", util.GetRemoteUrl(isLocal, ControlPort), conf)
		return conn, err
	}
}

func closeConnections(localConn net.Conn, tunnel net.Conn) {
	if !checkClosed(localConn) {
		err := localConn.Close()
		if err != nil && !strings.Contains(err.Error(), "use of closed") {
			log.Println("Error while closing local connection ", err.Error())
			return
		}
	}
	if !checkClosed(tunnel) {
		err := tunnel.Close()
		if err != nil && !strings.Contains(err.Error(), "use of closed") {
			log.Println("Error while closing tunnel connection ", err.Error())
			return
		}
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

func createNewTunnel(message *proto.Message, isLocal bool) net.Conn {
	conn := createTunnelConnection(isLocal)
	mutex := sync.Mutex{}
	mutex.Lock()
	tunnelsMap[message.TunnelId] = conn
	mutex.Unlock()
	proto.SendMessage(message, conn)
	return conn
}

func createTunnelConnection(isLocal bool) net.Conn {
	if isLocal {
		conn, _ := net.Dial("tcp", util.GetRemoteUrl(isLocal, TunnelPort))
		return conn
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, _ := tls.Dial("tcp", util.GetRemoteUrl(isLocal, TunnelPort), conf)
	return conn
}

func createLocalConnection(server tunnels.LocalServer) (net.Conn, error) {
	if server.Tls {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := tls.Dial("tcp", server.ServerFqdn+":"+strconv.Itoa(int(server.Port)), conf)
		return conn, err
	} else {
		conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(int(server.Port)))
		return conn, err
	}

}
