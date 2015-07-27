package uninet

import (
	"strings"
	"sync/atomic"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var (
	ErrShortRead = logex.NewError("short read")
)

type ListenConfig struct {
	Tcp  *TcpURL
	Udp  *UdpURL
	Http *HttpURL
}

type blockSession struct {
	Block   *utils.Block
	Session Session
}

type UniListener struct {
	RetryTime int
	udp       *UDPListener
	tcp       *TCPListener
	pool      *utils.BlockPool
	stop      uint32

	blockSession chan *blockSession
}

func NewUniListener(lnConfig *ListenConfig, pool *utils.BlockPool) (*UniListener, error) {
	var (
		udp *UDPListener
		tcp *TCPListener
		err error
	)

	u := &UniListener{
		pool:         pool,
		blockSession: make(chan *blockSession),
	}

	if lnConfig.Udp != nil {
		udp, err = ListenUDP(lnConfig.Udp)
		if err != nil {
			return nil, logex.Trace(err)
		}
		u.udp = udp
		go u.readingUDP()
	}

	if lnConfig.Tcp != nil {
		tcp, err = ListenTcp(lnConfig.Tcp)
		if err != nil {
			return nil, logex.Trace(err)
		}
		u.tcp = tcp
		go u.readingTCP()
	}

	return u, nil
}

func (u *UniListener) readingUDP() {
	block := u.pool.Get()
	for atomic.LoadUint32(&u.stop) == 0 {
		session, err := u.udp.ReadBlock(block)
		if err != nil {
			block.Init()
			logex.Warn(err)
			continue
		}
		u.blockSession <- &blockSession{block, session}
		block = u.pool.Get()
	}
}

func (u *UniListener) readingTCP() {
	block := u.pool.Get()
	for atomic.LoadUint32(&u.stop) == 0 {
		session, err := u.tcp.ReadBlock(block)
		if err != nil {
			block.Init()
			logex.Warn(err)
			continue
		}
		u.blockSession <- &blockSession{block, session}
		block = u.pool.Get()
	}
}

func (u *UniListener) ReadBlock() (*utils.Block, Session, error) {
	// timeout
	bs := <-u.blockSession
	return bs.Block, bs.Session, nil
}

func (u *UniListener) WriteBlock(b *utils.Block, session Session) error {
	switch session.NetType() {
	case NET_TCP:
		if u.tcp == nil {
			return logex.NewError("tcp not init but use to send data")
		}
		return logex.Trace(u.tcp.WriteBlock(b, session.(*TcpSession)))
	case NET_UDP:
		if u.udp == nil {
			return logex.NewError("tcp not init but use to send data")
		}
		return logex.Trace(u.udp.WriteBlock(b, session.(*UdpSession)))
	}
	return logex.NewError("unknown session", session)
}

func IsClosedError(err error) bool {
	if err == nil {
		return false
	}
	if strings.HasSuffix(err.Error(), "use of closed network connection") {
		return true
	}
	return false
}
