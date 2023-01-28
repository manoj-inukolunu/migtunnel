package client

import (
	"encoding/json"
	"golang/jtunnel-client/admin/data"
	"golang/jtunnel-client/admin/tunnels"
	"log"
	"net/http"
	"strconv"
)

type Client struct {
	ClientConfig data.ClientConfig
}

func (client *Client) NewClient(clientConfig data.ClientConfig) {
	client.ClientConfig = clientConfig
}

func (client *Client) GetClientConfig() data.ClientConfig {
	return client.ClientConfig
}

func (client *Client) StartAdminServer() {
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/register", registerTunnelHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(int(client.ClientConfig.AdminPort)), nil)
	if err != nil {
		log.Printf("Could not start admin server on port=%s  error=%s\n", client.ClientConfig.AdminPort, err)
		panic(err)
	}
}

func registerTunnelHandler(writer http.ResponseWriter, request *http.Request) {
	dec := json.NewDecoder(request.Body)
	message := &data.TunnelData{}
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

func listHandler(writer http.ResponseWriter, request *http.Request) {
	list := tunnels.GetRegisteredTunnels()

	var listOfTunnels []data.TunnelData
	list.Range(func(key, value any) bool {
		listOfTunnels = append(listOfTunnels, data.TunnelData{HostName: key.(string),
			LocalServerPort: value.(int16)})
		return true
	})
	str, _ := json.Marshal(listOfTunnels)
	writer.Write(str)
}
