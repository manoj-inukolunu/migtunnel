package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/thejerf/suture/v4"
	"golang/server"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

type Main struct {
}

func (i *Main) Stop() {
	fmt.Println("Stopping the service")
}

func main() {
	supervisor := suture.NewSimple("Main")
	service := &Main{}
	ctx, cancel := context.WithCancel(context.Background())
	supervisor.Add(service)
	errors := supervisor.ServeBackground(ctx)
	fmt.Println(<-errors)
	cancel()

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
		fmt.Println("Using TLS")
		cer, err := tls.LoadX509KeyPair(data["certFile"], data["keyFile"])
		if err != nil {
			log.Println(err)
			return err
		}
		config = &tls.Config{Certificates: []tls.Certificate{cer}}
		fmt.Println("Created TLS config")
	} else {
		fmt.Println("Not using TLS")
	}

	tunnelServerConfig := server.TunnelServerConfig{
		ClientTunnelServerPort:  9999,
		ClientControlServerPort: 2121,
		ServerHttpServerPort:    2020,
		ServerAdminServerPort:   9090,
		ServerTlsConfig:         config,
	}
	server.Start(tunnelServerConfig)

	return nil
}
