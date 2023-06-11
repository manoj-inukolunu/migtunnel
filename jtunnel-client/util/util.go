package util

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"golang/jtunnel-client/data"
	"net/http"
	"strconv"
)

func GetTunnels(adminUrl string) ([]data.TunnelCreateRequest, error) {
	response, err := http.Get(adminUrl)
	if err != nil {
		return nil, err
	}
	var d []data.TunnelCreateRequest
	err = json.NewDecoder(response.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func RegisterTunnel(adminUrl string, tunnelRegisterRequest data.TunnelCreateRequest) error {
	jsonData, _ := json.Marshal(tunnelRegisterRequest)
	_, err := http.Post(adminUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}

func GetRemoteUrl(isLocal bool, port int16) string {
	if isLocal {
		return "localhost:" + strconv.Itoa(int(port))
	} else {
		return uuid.New().String() + ".migtunnel.net:" + strconv.Itoa(int(port))
	}
}
