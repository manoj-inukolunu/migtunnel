package http

import (
	"encoding/json"
	"golang/client/admin/data"
	"golang/client/admin/tunnels"
	"log"
	"net/http"
	"strconv"
)

type Client struct {
	ClientConfig data.ClientConfig
}

func (client *Client) NewCustomer(clientConfig data.ClientConfig) {
	client.ClientConfig = clientConfig
}

func (client Client) StartAdminServer() {
	http.HandleFunc("/stats", adminHandler)
	http.HandleFunc("/create", createTunnelHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(int(client.ClientConfig.AdminPort)), nil)
	if err != nil {
		log.Printf("Could not start admin server on port=%s  error=%s\n", client.ClientConfig.AdminPort, err)
		panic(err)
	}
}

func createTunnelHandler(writer http.ResponseWriter, request *http.Request) {
	dec := json.NewDecoder(request.Body)
	message := &data.TunnelCreateRequest{}
	err := dec.Decode(message)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not read request payload"))
		return
	}
	err = tunnels.RegisterTunnel(*message)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("Could not register tunnel" + err.Error()))
		return
	}

}

func adminHandler(writer http.ResponseWriter, request *http.Request) {

}
