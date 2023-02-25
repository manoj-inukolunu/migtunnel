package client

import (
	"crypto/tls"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"golang/jtunnel-client/admin/tunnels"
	"golang/jtunnel-client/util"
	"golang/proto"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var db *badger.DB

func init() {
	db, _ = badger.Open(badger.DefaultOptions("/Users/minukolunu/badgerT"))
}

func ListAll(writer http.ResponseWriter) error {
	return db.View(func(txn *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				_, err := writer.Write([]byte(fmt.Sprintf("key=%s, value=%s\n", k, v)))
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil

	})
}

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
			localConn, localConnErr := createLocalConnection(tunnels.GetPortForHostName(message.HostName))
			if localConnErr != nil {
				log.Printf("Could not connect to local server on port %d "+
					"Please check if server is running.\n", tunnels.GetPortForHostName(message.HostName))
				continue
			}
			log.Println("Created Local Connection", localConn.RemoteAddr())
			tunnelProcessor := util.NewTeeReader(db, message.TunnelId, tunnel, localConn)
			log.Println("Writing data to local Connection")
			sig := make(chan bool)
			go func() {
				err := tunnelProcessor.ReadFromTunnel()
				if err != nil && !strings.Contains(err.Error(), "use of closed") {
					log.Println("Error reading from tunnel ", err.Error())
				}
				sig <- true
			}()
			err := tunnelProcessor.WriteToTunnel()
			if err != nil && !strings.Contains(err.Error(), "use of closed") {
				log.Println("Error writing to tunnel ", err.Error())
			}
			log.Println("Finished Writing data to tunnel")
			<-sig
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
		err := localConn.Close()
		if err != nil && !strings.Contains(err.Error(), "use of closed") {
			log.Println("Error while closing local connection ", err.Error())
			return
		}
		if !checkClosed(tunnel) {
			err := tunnel.Close()
			if err != nil && !strings.Contains(err.Error(), "use of closed") {
				log.Println("Error while closing tunnel connection ", err.Error())
				return
			}
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

func createLocalConnection(port int16) (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(int(port)))
	return conn, err
}
