package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/tidwall/evio"
	"jtunnel-go/proto"
	"log"
	"net/http"
)

var response = "HTTP/1.1 200 OK\nDate: Sun, 10 Oct 2010 23:26:07 GMT\nServer: Apache/2.2.8 (Ubuntu) mod_ssl/2.2.8 OpenSSL/0.9.8g\nLast-Modified: Sun, 26 Sep 2010 22:04:35 GMT\nETag: \"45b6-834-49130cc1182c0\"\nAccept-Ranges: bytes\nContent-Length: 12\nConnection: close\nContent-Type: text/html\n\nHello world!"

func FastHttp() {
	log.Println("Starting Listener on 8080")
	var events evio.Events
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		reader := bytes.NewReader(in)
		request, err := http.ReadRequest(bufio.NewReader(reader))
		if err != nil {
			fmt.Println("Failed to parse Request", err.Error())
		}
		id := uuid.New().String()
		message := proto.NewMessage(request.Host, id, "data", in)
		tunnelConn, err := GetTunnelConnectionFromHostName(request.Host)
		AddHttpConnection(id)
		log.Println("Writing Message")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		proto.SendMessage(message, tunnelConn)
		for {
			if GetHttpConnection(id) != nil {
				log.Println("Found Response")
				resp := GetHttpConnection(id)
				out = resp
				RemoveHttpConnection(id)
				log.Println("Finished Serving request")
				break
			}
		}
		return
	}
	if err := evio.Serve(events, "tcp://localhost:8080"); err != nil {
		panic(err.Error())
	}
}
