package uninet

import (
	"errors"
	"net"

	"gopkg.in/logex.v1"
)

var (
	ErrShortWritten = errors.New("short written")
)

type UDP struct {
	conn *net.UDPConn
}

func NewUDP(host string) (*UDP, error) {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return nil, logex.Trace(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, logex.Trace(err)
	}
	conn.SetWriteBuffer(512)
	conn.SetReadBuffer(512)
	u := &UDP{
		conn: conn,
	}
	return u, nil
}

func (c *UDP) ReadBlock(b []byte) (int, *net.UDPAddr, error) {
	n, addr, err := c.conn.ReadFromUDP(b)
	if err != nil {
		return 0, nil, logex.Trace(err)
	}

	return n, addr, nil
}

func (c *UDP) WriteBlock(b []byte, addr *net.UDPAddr) error {
	n, err := c.conn.WriteToUDP(b, addr)
	if err != nil {
		return logex.Trace(err)
	}
	if n != len(b) {
		return logex.Trace(ErrShortWritten)
	}
	return nil
}

func (c *UDP) Close() error {
	return c.conn.Close()
}
