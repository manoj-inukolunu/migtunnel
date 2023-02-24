package main

import (
	"context"
	"crypto/tls"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang/server/admin"
	"golang/server/control"
	myhttp "golang/server/http"
	"golang/server/tunnel"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

type Main struct {
}

type TunnelServerConfig struct {
	ClientControlServerPort int
	ServerHttpServerPort    int
	ClientTunnelServerPort  int
	ServerAdminServerPort   int
	ServerTlsConfig         *tls.Config
}

func start(tunnelServerConfig TunnelServerConfig) {
	useTLS := tunnelServerConfig.ServerTlsConfig != nil
	controlManager := control.Server{
		ControlConnections: make(map[string]net.Conn),
		ControlServerPort:  tunnelServerConfig.ClientControlServerPort,
		UseTLS:             useTLS,
	}
	controlManager.InitCronitorHeartbeat()
	controlManager.CheckConnections()
	httpChan := make(chan string)
	//Start all the servers
	httpServer := myhttp.Server{
		TunnelChannel:  httpChan,
		Port:           tunnelServerConfig.ServerHttpServerPort,
		ControlManager: controlManager,
	}
	tunnelServer := tunnel.Server{
		Port:           tunnelServerConfig.ClientTunnelServerPort,
		HttpServerChan: httpChan,
		TlsConfig:      tunnelServerConfig.ServerTlsConfig,
	}
	go httpServer.Start()
	go admin.StartAdminServer(tunnelServerConfig.ServerAdminServerPort, controlManager)
	tunnelServer.Start()
	controlManager.Start()

}

func (i *Main) Stop() {
	log.Println("Stopping the service")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover called from main error is , program exiting", r)
		}
	}()
	go func() {
		log.Println("Metrics endpoing at /metrics")
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
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

func (i *Main) Serve(ctx context.Context) error {
	yfile, err := os.ReadFile("server.yaml")
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]string)
	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		return err2
	}

	var config *tls.Config
	if ok, _ := strconv.ParseBool(data["useTLS"]); ok {
		log.Println("Using TLS")
		cer, err := tls.LoadX509KeyPair(data["certFile"], data["keyFile"])
		if err != nil {
			log.Println(err)
			return err
		}
		config = &tls.Config{Certificates: []tls.Certificate{cer}}
		log.Println("Created TLS config")
	} else {
		log.Println("Not using TLS")
	}

	tunnelServerConfig := TunnelServerConfig{
		ClientTunnelServerPort:  9999,
		ClientControlServerPort: 2121,
		ServerHttpServerPort:    2020,
		ServerAdminServerPort:   9090,
		ServerTlsConfig:         config,
	}
	start(tunnelServerConfig)

	return nil
}
