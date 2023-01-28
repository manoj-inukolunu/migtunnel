package checks

import (
	"go.uber.org/zap"
	"golang/client/admin/tunnels"
	"net"
	"time"
)

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

func CheckLocalServerPorts() {
	registeredTunnels := tunnels.GetRegisteredTunnels()

	registeredTunnels.Range(func(key, value any) bool {
		if rawConnect("localhost", value.(string)) {
			sugar.Infof("SUCCESS Connection to server on port %s successful", value)
		} else {
			sugar.Errorf("FAIL could not connect to server on port %s", value)
		}

		return true
	})

}

func rawConnect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		sugar.Errorf("Could not connect to port %s %s", port, err.Error())
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return true
}
