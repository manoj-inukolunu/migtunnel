package util

import (
	"golang/jtunnel-client/data"
	"golang/jtunnel-client/db"
	"io"
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
}

func NewTeeReader(requestId string, tunnelConn net.Conn, localConn net.Conn, db db.LocalDb) TeeReader {
	return TeeReader{
		responseData: []byte{},
		requestData:  []byte{},
		requestId:    requestId,
		tunnelConn:   tunnelConn,
		localConn:    localConn,
		timeStamp:    time.Now().UnixNano(),
		Db:           db,
	}
}

func (t *TeeReader) ReadFromTunnel() error {
	oneKB := 32 * 1024
	buf := make([]byte, oneKB)
	for {
		nr, err := t.tunnelConn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if nr > 0 {
			_, err := t.localConn.Write(buf[0:nr])
			if err != nil {
				return err
			}
			t.requestData = append(t.requestData, buf[0:nr]...)
		}
	}
}

func (t *TeeReader) WriteToTunnel() error {
	oneKB := 32 * 1024
	buf := make([]byte, oneKB)
	for {
		nr, err := t.localConn.Read(buf)
		if err != nil {
			if err == io.EOF {
				t.tunnelConn.Close()
				t.localConn.Close()
				return t.save()
			}
			return err
		}
		if nr > 0 {
			_, err := t.tunnelConn.Write(buf[0:nr])
			if err != nil {
				return err
			}
			t.responseData = append(t.responseData, buf[0:nr]...)

		}
	}
}

func (t *TeeReader) save() error {
	return t.Db.Save(data.TunnelData{
		Id:           t.timeStamp,
		TunnelId:     t.requestId,
		IsReplay:     false,
		RequestData:  t.requestData,
		ResponseData: t.responseData,
	})
}
