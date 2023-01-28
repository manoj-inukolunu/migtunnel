package util

import (
	"bytes"
	"encoding/json"
	"golang/jtunnel-client/admin/data"
	"net/http"
)

func GetTunnels(adminUrl string) ([]data.TunnelData, error) {
	response, err := http.Get(adminUrl)
	if err != nil {
		return nil, err
	}
	var d []data.TunnelData
	err = json.NewDecoder(response.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func RegisterTunnel(adminUrl string, tunnelRegisterRequest data.TunnelData) error {
	jsonData, _ := json.Marshal(tunnelRegisterRequest)
	_, err := http.Post(adminUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	return nil
}
