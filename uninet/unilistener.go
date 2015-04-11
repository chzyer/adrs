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

func NewUniListener(tcpUrl *TcpURL, udpUrl *UdpURL, pool *utils.BlockPool) (*UniListener, error) {
	var (
		udp *UDPListener
		tcp *TCPListener
		err error
	)

	u := &UniListener{
		pool:         pool,
		blockSession: make(chan *blockSession),
	}

	if udpUrl != nil {
		udp, err = ListenUDP(udpUrl)
		if err != nil {
			return nil, logex.Trace(err)
		}
		u.udp = udp
		go u.readingUDP()
	}

	if tcpUrl != nil {
		tcp, err = ListenTcp(tcpUrl)
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
			return logex.NewTraceError("tcp not init but use to send data")
		}
		return logex.TraceIfError(u.tcp.WriteBlock(b, session.(*TcpSession)))
	case NET_UDP:
		if u.udp == nil {
			return logex.NewTraceError("tcp not init but use to send data")
		}
		return logex.TraceIfError(u.udp.WriteBlock(b, session.(*UdpSession)))
	}
	return logex.NewTraceError("unknown session", session)
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