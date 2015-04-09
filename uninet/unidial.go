package uninet

import (
	"net"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type DialAddr struct {
	UDP string
}

type UniDial struct {
	Addr      *DialAddr
	UDP       *UDP
	RetryTime int
}

func NewUniDial(addr *DialAddr) (*UniDial, error) {
	udp, err := NewDialUDP(addr.UDP)
	if err != nil {
		return nil, logex.Trace(err)
	}

	ud := &UniDial{
		UDP: udp,
	}
	return ud, nil
}

func (u *UniDial) ReadBlockUDP(b *utils.Block) (addr *net.UDPAddr, err error) {
	addr, err = u.UDP.ReadBlock(b)
	if err != nil {
		return nil, logex.Trace(err)
	}

	return
}

func (u *UniDial) WriteBlockUDP(b *utils.Block) (err error) {
	i := 0
	for ; i < u.RetryTime || u.RetryTime == 0; i++ {
		err = u.UDP.WriteBlock(b)
		if err == nil {
			break
		}
	}

	if err != nil {
		return logex.Trace(err)
	}

	if i > 0 {
		logex.Warn("write udp to", u.Addr.UDP, "fail", i, "times")
	}
	return
}
