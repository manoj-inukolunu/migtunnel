package main

import (
	"crypto/tls"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang/server/admin"
	"golang/server/control"
	myhttp "golang/server/http"
	"golang/server/tunnel"
	tunnelmanager "golang/server/tunnel-manager"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

type TunnelServerConfig struct {
	ClientControlServerPort int
	ServerHttpServerPort    int
	ClientTunnelServerPort  int
	ServerAdminServerPort   int
	ServerTlsConfig         *tls.Config
}

func start(tunnelServerConfig TunnelServerConfig) {
	useTLS := tunnelServerConfig.ServerTlsConfig != nil
	controlManager := control.ControlManager{
		ControlConnections: make(map[string]net.Conn),
		ControlServerPort:  tunnelServerConfig.ClientControlServerPort,
		UseTLS:             useTLS,
		ServerTlsConfig:    tunnelServerConfig.ServerTlsConfig,
	}
	tunnelManager := tunnelmanager.TunnelManager{
		TunnelConnections:  make(map[string]net.Conn),
		HttpServerChannels: make(map[string]chan bool),
	}
	controlManager.InitCronitorHeartbeat()
	controlManager.CheckConnections()
	//Start all the servers
	httpServer := myhttp.Server{
		Port:           tunnelServerConfig.ServerHttpServerPort,
		ControlManager: controlManager,
		TunnelManager:  tunnelManager,
	}
	tunnelServer := tunnel.Server{
		Port:          tunnelServerConfig.ClientTunnelServerPort,
		TlsConfig:     tunnelServerConfig.ServerTlsConfig,
		TunnelManager: tunnelManager,
		UseTls:        useTLS,
	}
	adminServer := admin.Server{
		TunnelManger:   tunnelManager,
		ControlManager: controlManager,
		Port:           tunnelServerConfig.ServerAdminServerPort,
	}
	go httpServer.Start()
	go adminServer.Start()
	tunnelServer.Start()
	controlManager.Start()

}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover called from main error is ", r)
		}
	}()
	go func() {
		log.Println("Metrics endpoing at /metrics")
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	run()

}

func run() {
	yfile, err := os.ReadFile("server.yaml")
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]string)
	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		return
	}

	var config *tls.Config
	if ok, _ := strconv.ParseBool(data["useTLS"]); ok {
		log.Println("Using TLS")
		cer, err := tls.LoadX509KeyPair(data["certFile"], data["keyFile"])
		if err != nil {
			log.Println(err)
			return
		}
		config = &tls.Config{Certificates: []tls.Certificate{cer}}
		log.Println("Created TLS config")
	} else {
		log.Println("Not using TLS")
	}

	tunnelServerConfig := TunnelServerConfig{
		ClientTunnelServerPort:  2121,
		ClientControlServerPort: 9999,
		ServerHttpServerPort:    2020,
		ServerAdminServerPort:   9090,
		ServerTlsConfig:         config,
	}
	start(tunnelServerConfig)
}
