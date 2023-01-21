package main

import (
	"context"
	"crypto/tls"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang/server"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

type Main struct {
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

	/*sigc := make(chan os.Signal, 1)

	signal.Notify(sigc)

	go func() {
		s := <-sigc
		log.Println(s.String())
	}()*/

	go func() {
		log.Println("Metrics endpoing at /metrics")
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	/*log.Println("Starting Main with supervisor")
	supervisor := suture.NewSimple("Main")
	service := &Main{}
	ctx, cancel := context.WithCancel(context.Background())
	supervisor.Add(service)
	errors := supervisor.ServeBackground(ctx)
	log.Println(<-errors)
	cancel()*/
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

	tunnelServerConfig := server.TunnelServerConfig{
		ClientTunnelServerPort:  2121,
		ClientControlServerPort: 9999,
		ServerHttpServerPort:    2020,
		ServerAdminServerPort:   9090,
		ServerTlsConfig:         config,
	}
	server.Start(tunnelServerConfig)
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
