package uninet

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var (
	ErrNotSuchProtocol = logex.NewError("not supported such protocol yet: %s")
)

type UniDial struct {
	netType   NetType
	udp       *UDPDialer
	tcp       *TCPDialer
	RetryTime int
}

func NewUniDial(addr URLer) (*UniDial, error) {
	ud := &UniDial{
		netType: addr.NetType(),
	}
	var (
		err error
	)

	switch addr.NetType() {
	case NET_UDP:
		ud.udp, err = DialUDP(addr.(*UdpURL))
		if err != nil {
			return nil, logex.Trace(err)
		}
	case NET_TCP:
		ud.tcp, err = DialTCP(addr.(*TcpURL))
		if err != nil {
			return nil, logex.Trace(err)
		}
	default:
		return nil, logex.TraceError(ErrNotSuchProtocol).Format(addr.Network())
	}

	return ud, nil
}

func (u *UniDial) ReadBlock(b *utils.Block) (err error) {
	switch u.netType {
	case NET_UDP:
		err = u.udp.ReadBlock(b)
		if err != nil {
			return logex.Trace(err)
		}
	case NET_TCP:
		err = u.tcp.ReadBlock(b)
		if err != nil {
			return logex.Trace(err)
		}
	}
	return
}

func (u *UniDial) WriteBlock(b *utils.Block) (err error) {
	switch u.netType {
	case NET_UDP:
		err = u.udp.WriteBlock(b)
		if err != nil {
			return logex.Trace(err)
		}
	case NET_TCP:
		err = u.tcp.WriteBlock(b)
		if err != nil {
			return logex.Trace(err)
		}
	}
	return
}
