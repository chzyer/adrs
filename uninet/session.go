package uninet

import (
	"net"
)

type NetType int

func (n *NetType) String() string {
	switch *n {
	case NET_UDP:
		return "udp"
	case NET_TCP:
		return "tcp"
	case NET_HTTP:
		return "http"
	}
	return "unknown"
}

const (
	NET_UDP NetType = iota
	NET_TCP
	NET_HTTP
)

type Session interface {
	NetType() NetType
}
type BaseSession struct {
	Type NetType
}

func NewBaseSession(netType NetType) *BaseSession {
	return &BaseSession{
		Type: netType,
	}
}

func (b *BaseSession) NetType() NetType {
	return b.Type
}

func (s *BaseSession) Network() string {
	return s.Type.String()
}

func NewUdpSession(addr *net.UDPAddr) *UdpSession {
	return &UdpSession{
		BaseSession: NewBaseSession(NET_UDP),
		Addr:        addr,
	}
}

type TcpSession struct {
	*BaseSession
	Addr *net.TCPAddr
	conn *tcpConn
}

func NewTcpSession(addr *net.TCPAddr, conn *tcpConn) *TcpSession {
	return &TcpSession{
		BaseSession: NewBaseSession(NET_TCP),
		Addr:        addr,
		conn:        conn,
	}
}

type UdpSession struct {
	*BaseSession
	Addr *net.UDPAddr
}

type HttpSession struct {
	*BaseSession
}
