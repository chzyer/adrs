package uninet

import (
	"net"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
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

type Host struct {
	UDP  string
	TCP  string
	HTTP string
}

type UniNet struct {
	RetryTime int
	host      *Host
	pool      *utils.BlockPool
	UDP       *UDP
}

func NewUniNet(host *Host, pool *utils.BlockPool) (*UniNet, error) {
	udp, err := NewUDP(host.UDP)
	if err != nil {
		return nil, logex.Trace(err)
	}
	u := &UniNet{
		host: host,
		pool: pool,
		UDP:  udp,
	}
	return u, nil
}

func (u *UniNet) ReadBlockUDP() (n int, b []byte, addr *net.UDPAddr, err error) {
	b = u.pool.Get()
	n, addr, err = u.UDP.ReadBlock(b)
	if err != nil {
		u.pool.Put(b)
		return 0, nil, nil, logex.Trace(err)
	}

	return
}

func (u *UniNet) WriteBlockUDP(b []byte, addr *net.UDPAddr) (err error) {
	i := 0
	for ; i < u.RetryTime || u.RetryTime == 0; i++ {
		err = u.UDP.WriteBlock(b, addr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return logex.Trace(err)
	}

	if i > 0 {
		logex.Warn("write udp to", addr, "fail", i, "times")
	}
	return
}
