/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/tls"
	"github.com/thejerf/suture/v4"
	"golang/client/admin/data"
	"golang/client/admin/http"
	tunnels2 "golang/client/admin/tunnels"
	"golang/proto"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"context"
	markdown "github.com/MichaelMure/go-term-markdown"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		supervisor := suture.NewSimple("Client")
		service := &Main{}
		ctx, cancel := context.WithCancel(context.Background())
		supervisor.Add(service)
		errors := supervisor.ServeBackground(ctx)
		sugar.Error(<-errors)
		cancel()
	},
}

var ControlConnections map[string]net.Conn
var tunnels map[string]net.Conn

const usage = "Welcome to JTunnel .\n\nSource code is at `https://github.com/manoj-inukolunu/jtunnel-go`\n\nTo create a new tunnel\n\nMake a `POST` request to `http://127.0.0.1:1234/create`\nwith the payload\n\n```\n{\n    \"HostName\":\"myhost\",\n    \"TunnelName\":\"Tunnel Name\",\n    \"localServerPort\":\"3131\"\n}\n\n```\n\nThe endpoint you get is `https://myhost.migtunnel.net`\n\nAll the requests to `https://myhost.migtunnel.net` will now\n\nbe routed to your server running on port `3131`\n\n"

type Main struct {
}

func (i *Main) Serve(ctx context.Context) error {
	result := markdown.Render(usage, 80, 6)
	log.Println(string(result))
	ControlConnections = make(map[string]net.Conn)
	tunnels = make(map[string]net.Conn)
	sugar.Infow("Starting Admin Server on ", "port", 1234)
	go http.StartServer(data.ClientConfig{AdminPort: 1234})
	startControlConnection()
	return nil
}

func (i *Main) Stop() {
	sugar.Info("Stopping Client")
}

func createNewTunnel(message *proto.Message) net.Conn {
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	conn, _ := tls.Dial("tcp", "manoj.migtunnel.net:2121", conf)
	mutex := sync.Mutex{}
	mutex.Lock()
	tunnels[message.TunnelId] = conn
	mutex.Unlock()
	proto.SendMessage(message, conn)
	return conn
}

func createLocalConnection(port int16) net.Conn {
	conn, _ := net.Dial("tcp", "localhost:"+strconv.Itoa(int(port)))
	return conn
}

func startControlConnection() {
	sugar.Info("Starting Control connection")
	conf := &tls.Config{
		//InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "manoj.migtunnel.net:9999", conf)
	if err != nil {
		sugar.Errorw("Failed to establish control connection ", "Error", err.Error())
		panic(err)
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	ControlConnections["data"] = conn
	tunnels2.SaveControlConnection(conn)
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
			localConn := createLocalConnection(tunnels2.GetPortForHostName(message.HostName))
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
			tunnels2.UpdateHostNameToPortMap(message.HostName, port)
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

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
