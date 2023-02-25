package util

import (
	"github.com/dgraph-io/badger/v3"
	"io"
	"log"
	"net"
)

type TeeReader struct {
	db           *badger.DB
	requestData  []byte
	responseData []byte
	requestId    string
	tunnelConn   net.Conn
	localConn    net.Conn
}

func NewTeeReader(db *badger.DB, requestId string, tunnelConn net.Conn, localConn net.Conn) TeeReader {
	return TeeReader{
		db:           db,
		responseData: []byte{},
		requestData:  []byte{},
		requestId:    requestId,
		tunnelConn:   tunnelConn,
		localConn:    localConn,
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
	return t.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(t.requestId+":"+"request"), t.requestData)
		if err != nil {
			return err
		}
		log.Println(t.responseData)
		err = txn.Set([]byte(t.requestId+":"+"response"), t.responseData)
		if err != nil {
			return err
		}
		return nil
	})
}
