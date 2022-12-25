package admin

import (
	"fmt"
	control_manager "golang/control-manager"
	tunnel_manager "golang/tunnel-manager"
	"log"
	"net/http"
)

func StartAdminServer(port int) {

	http.HandleFunc("/tunnels", listTunnels)
	http.HandleFunc("/control", listControlConnections)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Printf("Could not start admin server on port=%s  error=%s\n", port, err)
		panic(err)
	}
	fmt.Println("Started Admin Server on ", 8090)
}

func listControlConnections(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte(control_manager.ListAllConnectionsAsString()))
	if err != nil {
		fmt.Println("Failed to list tunnels", err)
	}
}

func listTunnels(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte(tunnel_manager.ListAllConnectionsAsString()))
	if err != nil {
		fmt.Println("Failed to list tunnels", err)
	}
}
