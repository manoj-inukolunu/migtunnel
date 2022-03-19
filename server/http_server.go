package server

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"jtunnel-go/proto"
	"log"
	"net"
)

func FastHttp() {
	log.Println("Starting Listener on 8080")
	s := &fasthttp.Server{
		Handler: requestHandler,

		// Every response will contain 'Server: My super server' header.
		Name: "Http Server",

		// Other Server settings may be set here.

		MaxRequestBodySize: 4 * 1024 * 1024 * 1024,
	}
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Failed to start listener on 9999", err.Error())
		return
	}
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	//fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())

	hostName := ctx.Host()
	fmt.Printf("Host Name %s\n", hostName)
	//&ctx.Request
	rawRequest := &ctx.Request
	var buf bytes.Buffer
	fmt.Fprintf(&buf, rawRequest.String())
	id := uuid.New().String()
	message := proto.NewMessage(string(hostName), id, "data", buf.Bytes())
	tunnelConn, err := GetTunnelConnectionFromHostName(string(hostName))
	AddHttpConnection(ctx, id)
	log.Println("Writing Message")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	proto.SendMessage(message, tunnelConn)
	for {
		if GetHttpConnection(id) == nil {
			log.Println("Finished Serving request")
			break
		}
	}
}
