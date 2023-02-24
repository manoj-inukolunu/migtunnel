package admin

import (
	control_manager "golang/server/control"
	tunnel_manager "golang/tunnel-manager"
	"log"
	"net/http"
)

func StartAdminServer(port int, manager control_manager.Server) {

	http.HandleFunc("/tunnels", listTunnels)
	http.HandleFunc("/control", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(manager.ListAllConnectionsAsString()))
		if err != nil {
			log.Println("Failed to list tunnels", err)
		}
	})

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Printf("Could not start admin server on port=%s  error=%s\n", port, err)
		panic(err)
	}
	log.Println("Started Admin Server on ", 8090)
}

func listTunnels(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte(tunnel_manager.ListAllConnectionsAsString()))
	if err != nil {
		log.Println("Failed to list tunnels", err)
	}
}
