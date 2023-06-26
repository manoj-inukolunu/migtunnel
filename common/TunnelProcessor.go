package common

import (
	"io"
	"net"
)

type TeeTunnel struct {
	Src net.Conn
	Dst net.Conn
}

func (t *TeeTunnel) Init(src net.Conn, dst net.Conn) {
	t.Src = src
	t.Dst = dst
}

func (t *TeeTunnel) CopySrcToDest() error {
	/*for {
		size := 32 * 1024
		buf := make([]byte, size)
		numRead, err := t.Src.Read(buf)
		if err != nil {
			if err == io.EOF {
				if err := t.Src.Close(); err != nil {
					log.Println("Failed to close source connection", err.Error())
					return err
				}
				return nil
			}
			return err
		}
		numWrite, err := t.Dst.Write(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Reached end Nothing to write")
				return nil
			}
			return err
		}

		if numRead != numWrite {
			return errors.New("num Written is not equal to Num Read")
		}
	}*/
	_, err := io.Copy(t.Dst, t.Src)
	return err
}
