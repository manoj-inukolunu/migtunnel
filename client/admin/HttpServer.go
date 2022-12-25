package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TunnelCreateRequest struct {
	HostName        string
	TunnelName      string
	localServerPort int
}

func StartServer(port int) {

	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/create", createTunnelHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Printf("Could not start admin server on port=%s  error=%s\n", port, err)
		panic(err)
	}

}

func createTunnelHandler(writer http.ResponseWriter, request *http.Request) {
	dec := json.NewDecoder(request.Body)
	message := &TunnelCreateRequest{}
	err := dec.Decode(message)
	if err != nil {
		writer.Write([]byte("Could not read request payload"))
		writer.WriteHeader(400)
		return
	}
	err = CreateTunnel(*message)
	if err != nil {
		writer.Write([]byte("Could not create tunnel" + err.Error()))
		writer.WriteHeader(400)
		return
	}

}

func adminHandler(writer http.ResponseWriter, request *http.Request) {

}
