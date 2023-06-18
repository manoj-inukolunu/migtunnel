package util

import (
	"golang/migtunnel-client/data"
	"golang/migtunnel-client/db"
	"golang/migtunnel-client/tunnels"
	"io"
	"log"
	"net"
	"time"
)

type TeeReader struct {
	requestData  []byte
	responseData []byte
	requestId    string
	tunnelConn   net.Conn
	localConn    net.Conn
	timeStamp    int64
	Db           db.LocalDb
	isReplay     bool
	localPort    int16
}

func NewTeeReader(requestId string, tunnelConn net.Conn, localConn net.Conn, db db.LocalDb, isReplay bool, localServer tunnels.LocalServer) TeeReader {

	return TeeReader{
		responseData: []byte{},
		requestData:  []byte{},
		requestId:    requestId,
		tunnelConn:   tunnelConn,
		localConn:    localConn,
		timeStamp:    time.Now().UnixNano(),
		Db:           db,
		isReplay:     isReplay,
		localPort:    localServer.Port,
	}
}

func (t *TeeReader) TunnelToLocal() error {
	for {
		oneByte := 1
		buf := make([]byte, oneByte)
		nr, err := t.tunnelConn.Read(buf)
		if err != nil {
			if err == io.EOF {
				processTunnelToLocalData(t, buf, nr)
				return nil
			}
			return err
		}
		processTunnelToLocalData(t, buf, nr)
	}
}

func processTunnelToLocalData(t *TeeReader, buf []byte, numRead int) {
	if numRead > 0 {
		written, err := t.localConn.Write(buf[0:numRead])
		if err != nil || written != numRead {
			log.Println("Failed to write to local connection", err.Error())
		}
		t.requestData = append(t.requestData, buf[0:numRead]...)
	}
}

func (t *TeeReader) LocalToTunnel() error {
	for {
		oneByte := 1
		buf := make([]byte, oneByte)
		nr, err := t.localConn.Read(buf)
		if err != nil {
			log.Println("Finished reading from local connection")
			if err == io.EOF {
				processLocalToTunnelData(t, buf, nr)
				closeErr := t.tunnelConn.Close()
				if closeErr != nil {
					log.Println("Failed to close tunnel connection", err.Error())
				}
				t.localConn.Close()
				return t.save()
			}
			return err
		}
		processLocalToTunnelData(t, buf, nr)
	}
}

func processLocalToTunnelData(t *TeeReader, buf []byte, numRead int) {
	if numRead > 0 {
		numWrite, err := t.tunnelConn.Write(buf[0:numRead])
		if err != nil || numWrite != numRead {
			log.Println("Failed to write data to Tunnel ", err.Error())
		}
		t.responseData = append(t.responseData, buf[0:numRead]...)
	}
}

func (t *TeeReader) save() error {
	log.Println(string(t.responseData))
	return t.Db.Save(data.TunnelData{
		Id:           t.timeStamp,
		TunnelId:     t.requestId,
		IsReplay:     t.isReplay,
		RequestData:  t.requestData,
		ResponseData: t.responseData,
		LocalPort:    t.localPort,
	})
}
