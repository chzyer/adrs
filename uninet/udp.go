package uninet

import (
	"errors"
	"net"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var (
	ErrShortWritten = errors.New("short written")
)

type udpConn struct {
	*net.UDPConn
}

func newUDPConn(conn *net.UDPConn) *udpConn {
	/*
		conn.SetWriteBuffer(512)
		conn.SetReadBuffer(512)
	*/

	return &udpConn{
		UDPConn: conn,
	}
}

func (c *udpConn) ReadBlock(b *utils.Block) error {
	n, err := c.Read(b.Block)
	if err != nil {
		return logex.Trace(err)
	}

	b.Length = n
	return nil
}

func (c *udpConn) ReadBlockFrom(b *utils.Block) (*UdpSession, error) {
	n, addr, err := c.ReadFromUDP(b.Block)
	if err != nil {
		return nil, logex.Trace(err)
	}

	b.Length = n
	return NewUdpSession(addr), nil
}

func (c *udpConn) WriteBlock(b *utils.Block) error {
	n, err := c.Write(b.Bytes())
	if err != nil {
		return logex.Trace(err)
	}
	if n != b.Length {
		return logex.Trace(ErrShortWritten)
	}
	return nil
}

func (c *udpConn) WriteBlockTo(b *utils.Block, session *UdpSession) error {
	n, err := c.WriteToUDP(b.Bytes(), session.Addr)
	if err != nil {
		return logex.Trace(err)
	}
	if n != b.Length {
		return logex.Trace(ErrShortWritten)
	}
	return nil
}

type UDPListener struct {
	conn *udpConn
}

func (u *UDPListener) Close() error {
	return u.conn.Close()
}

func ListenUDP(url *UdpURL) (*UDPListener, error) {
	addr, err := net.ResolveUDPAddr(url.Network(), url.Host())
	if err != nil {
		return nil, logex.Trace(err)
	}
	conn, err := net.ListenUDP(url.Network(), addr)
	if err != nil {
		return nil, logex.Trace(err)
	}

	u := &UDPListener{
		conn: newUDPConn(conn),
	}
	return u, nil
}

func (u *UDPListener) WriteBlock(b *utils.Block, s *UdpSession) error {
	return logex.Trace(u.conn.WriteBlockTo(b, s))
}

func (u *UDPListener) ReadBlock(b *utils.Block) (*UdpSession, error) {
	session, err := u.conn.ReadBlockFrom(b)
	if err != nil {
		return nil, logex.Trace(err)
	}
	return session, nil
}

type UDPDialer struct {
	conn *udpConn
}

func DialUDP(addr *UdpURL) (*UDPDialer, error) {
	c, err := net.Dial(addr.Network(), addr.Host())
	if err != nil {
		return nil, logex.Trace(err)
	}
	conn := c.(*net.UDPConn)
	u := &UDPDialer{
		conn: newUDPConn(conn),
	}
	return u, nil
}

func (u *UDPDialer) ReadBlock(b *utils.Block) error {
	return logex.Trace(u.conn.ReadBlock(b))
}

func (u *UDPDialer) WriteBlock(b *utils.Block) error {
	return logex.Trace(u.conn.WriteBlock(b))
}

func (u *UDPDialer) Close() error {
	return u.conn.Close()
}
