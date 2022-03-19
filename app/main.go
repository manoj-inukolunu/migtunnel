package main

import "jtunnel-go/server"

func main() {
	go server.FastHttp()

	server.TunnelServer(9999)
}
