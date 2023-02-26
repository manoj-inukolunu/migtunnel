package client

import (
	"encoding/json"
	data2 "golang/jtunnel-client/data"
	"golang/jtunnel-client/db"
	"golang/jtunnel-client/tunnels"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	ClientConfig data2.ClientConfig
	Db           db.LocalDb
}

func NewClient(clientConfig data2.ClientConfig, db db.LocalDb) Client {
	client := Client{
		ClientConfig: clientConfig,
		Db:           db,
	}
	return client
}

func NewMemClient(clientConfig data2.ClientConfig) Client {
	client := Client{
		ClientConfig: clientConfig,
	}
	return client
}

func (client *Client) GetClientConfig() data2.ClientConfig {
	return client.ClientConfig
}

func (client *Client) StartAdminServer() {
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/register", registerTunnelHandler)
	http.HandleFunc("/all", func(writer http.ResponseWriter, request *http.Request) {
		start, _ := strconv.ParseInt(request.URL.Query().Get("start"), 10, 64)
		limit, _ := strconv.ParseInt(request.URL.Query().Get("limit"), 10, 32)
		data, err := client.Db.ListWithoutData(start, int(limit))
		if err != nil {
			writer.Write([]byte(err.Error()))
		}
		bytes, _ := json.Marshal(data)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(bytes)
	})
	http.HandleFunc("/request/", func(writer http.ResponseWriter, request *http.Request) {
		splits := strings.Split(request.URL.Path, "/")
		requestId, _ := strconv.ParseInt(splits[2], 10, 64)
		data, err := client.Db.Get(requestId)
		if err != nil {
			writer.Write([]byte(err.Error()))
		}
		if splits[3] == "requestData" {
			writer.Write(data.RequestData)
		} else if splits[3] == "responseData" {
			writer.Write(data.ResponseData)
		}

	})
	err := http.ListenAndServe(":"+strconv.Itoa(int(client.ClientConfig.AdminPort)), nil)
	if err != nil {
		log.Printf("Could not start admin server on port=%s  error=%s\n", client.ClientConfig.AdminPort, err)
		panic(err)
	}
}

func registerTunnelHandler(writer http.ResponseWriter, request *http.Request) {
	dec := json.NewDecoder(request.Body)
	message := &data2.TunnelCreateRequest{}
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

	var listOfTunnels []data2.TunnelCreateRequest
	list.Range(func(key, value any) bool {
		listOfTunnels = append(listOfTunnels, data2.TunnelCreateRequest{HostName: key.(string),
			LocalServerPort: value.(int16)})
		return true
	})
	str, _ := json.Marshal(listOfTunnels)
	writer.Write(str)
}
