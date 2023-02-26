package util

import (
	"bytes"
	"encoding/json"
	"golang/jtunnel-client/data"
	"net/http"
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
