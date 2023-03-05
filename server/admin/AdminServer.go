package admin

import (
	control_manager "golang/server/control"
	"golang/server/tunnel-manager"
	"log"
	"net/http"
)

type Server struct {
	TunnelManger   tunnel_manager.TunnelManager
	ControlManager control_manager.ControlManager
	Port           int
}

func (s *Server) Start() {

	http.HandleFunc("/tunnels", s.listTunnels)
	http.HandleFunc("/control", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(s.ControlManager.ListAllConnectionsAsString()))
		if err != nil {
			log.Println("Failed to list tunnels", err)
		}
	})

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Printf("Could not start ui server on port=%s  error=%s\n", 8090, err)
		panic(err)
	}
	log.Println("Started Admin ControlManager on ", 8090)
}

func (s *Server) listTunnels(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte(s.TunnelManger.ListAllConnectionsAsString()))
	if err != nil {
		log.Println("Failed to list tunnels", err)
	}
}
