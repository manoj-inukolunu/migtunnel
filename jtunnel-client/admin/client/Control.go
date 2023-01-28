package client

import (
	"go.uber.org/zap"
	"golang/jtunnel-client/admin/tunnels"
	"golang/proto"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

var controlConnections map[string]net.Conn
var tunnelsMap map[string]net.Conn

func init() {
	controlConnections = make(map[string]net.Conn)
	tunnelsMap = make(map[string]net.Conn)
}

func (client *Client) StartControlConnection() {
	sugar.Info("Starting Control connection")
	/*conf := &tls.Config{
		//InsecureSkipVerify: true,
	}*/
	//net.Dial("tcp", "localhost:9999")
	conn, err := net.Dial("tcp", "localhost:9999") /*tls.Dial("tcp", "manoj.migtunnel.net:9999", conf)*/
	if err != nil {
		sugar.Errorw("Failed to establish control connection ", "Error", err.Error())
		panic(err)
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	controlConnections["data"] = conn
	tunnels.SaveControlConnection(conn)
	mutex.Unlock()

	for {
		message, err := proto.ReceiveMessage(conn)
		sugar.Debug("Received Message", message)
		if err != nil {
			if err.Error() == "EOF" {
				panic("Server closed control connection stopping client now")
			}
			sugar.Errorw("Error on control connection ", "Error", err.Error())
		}
		if message.MessageType == "init-request" {
			tunnel := createNewTunnel(message)
			sugar.Infow("Created a new Tunnel", message)
			localConn := createLocalConnection(tunnels.GetPortForHostName(message.HostName))
			sugar.Infow("Created Local Connection", localConn.RemoteAddr())
			go func() {
				_, err := io.Copy(localConn, tunnel)
				if err != nil {
					closeConnections(localConn, tunnel)
				}
			}()
			sugar.Infow("Writing data to local Connection")
			_, err := io.Copy(tunnel, localConn)
			if err != nil {
				closeConnections(localConn, tunnel)
			}

			sugar.Infow("Finished Writing data to tunnel")
			closeConnections(localConn, tunnel)
		}
		if message.MessageType == "ack-tunnel-create" {
			sugar.Infow("Received Ack for creating tunnel from the upstream server")
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
		sugar.Infow("Detected closed Local connection")
		conn.Close()
		return true
	}
	return false
}

func createNewTunnel(message *proto.Message) net.Conn {
	/*conf := &tls.Config{
		//InsecureSkipVerify: true,
	}*/
	conn, _ := net.Dial("tcp", "localhost:2121") /*tls.Dial("tcp", "manoj.migtunnel.net:2121", conf)*/
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
