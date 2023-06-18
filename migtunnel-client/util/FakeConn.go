package util

import (
	"bytes"
	"net"
	"time"
)

type FakeConn struct {
	Reader bytes.Reader
}

func (c *FakeConn) Read(b []byte) (n int, err error) {
	read, err := c.Reader.Read(b)
	if err != nil {
		return -1, err
	}
	return read, nil
}

func (f *FakeConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (f *FakeConn) Close() error {
	return nil
}

func (f *FakeConn) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr returns the remote network address, if known.
func (f *FakeConn) RemoteAddr() net.Addr {
	return nil
}

func (f *FakeConn) SetDeadline(t time.Time) error {
	return nil
}

func (f *FakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (f *FakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
