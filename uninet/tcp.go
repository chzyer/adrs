package uninet

import (
	"net"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type tcpConn struct {
	*net.TCPConn
}

func newTcpConn(conn *net.TCPConn) *tcpConn {
	return &tcpConn{
		TCPConn: conn,
	}
}

func (c *tcpConn) ReadBlock(b *utils.Block) error {
	n, err := c.Read(b.All)
	if err != nil {
		return logex.Trace(err)
	}

	b.Length = n
	return nil
}

func (c *tcpConn) WriteBlock(b *utils.Block) error {
	n, err := c.Write(utils.Uint16To(uint16(b.Length)))
	if err != nil {
		return logex.Trace(err)
	}
	if n != 2 {
		return logex.Trace(ErrShortWritten)
	}
	n, err = c.Write(b.Bytes())
	if err != nil {
		return logex.Trace(err)
	}
	if n != b.Length {
		return logex.Trace(ErrShortWritten)
	}
	return nil
}

type TCPListener struct {
	listener     *net.TCPListener
	incomingConn chan *connHeader
}

type connHeader struct {
	conn   *tcpConn
	header int
}

func ListenTcp(url *TcpURL) (*TCPListener, error) {
	l, err := net.Listen(url.Network(), url.Host())
	if err != nil {
		return nil, logex.Trace(err)
	}
	ln := l.(*net.TCPListener)

	tl := &TCPListener{
		listener:     ln,
		incomingConn: make(chan *connHeader),
	}
	go tl.listen()
	return tl, nil
}

func (t *TCPListener) listen() {
	var (
		c   *net.TCPConn
		err error
	)
	for {
		c, err = t.listener.AcceptTCP()
		if err != nil {
			if IsClosedError(err) {
				break
			}
			logex.Warn(err)
			break
		}
		go t.accept(newTcpConn(c))
	}
}

// check header!
func (t *TCPListener) accept(c *tcpConn) bool {
	header := make([]byte, 2)
	n, err := c.Read(header)
	if err != nil || n != 2 {
		// EOF
		c.Close()
		return false
	}

	length := int(utils.ToUint16(header))
	if length > utils.BLOCK_CAP {
		c.Close()
		logex.Info("close")
		return false
	}
	t.incomingConn <- &connHeader{c, length}
	return true
}

func (t *TCPListener) Close() error {
	return logex.TraceIfError(t.listener.Close())
}

func (t *TCPListener) ReadBlock(b *utils.Block) (*TcpSession, error) {
	connHeader := <-t.incomingConn
	b.Length = connHeader.header
	n, err := connHeader.conn.Read(b.Bytes())
	if err != nil {
		return nil, logex.Trace(err)
	}
	if n != connHeader.header {
		return nil, logex.Trace(ErrShortRead)
	}

	return NewTcpSession(connHeader.conn.RemoteAddr().(*net.TCPAddr), connHeader.conn), nil
}

func (t *TCPListener) WriteBlock(b *utils.Block, session *TcpSession) error {
	err := session.conn.WriteBlock(b)
	if err != nil {
		return logex.Trace(err)
	}

	// if other message?
	/*
		if !t.accept(session.conn) {
			return nil
		}
	*/
	return nil
}
