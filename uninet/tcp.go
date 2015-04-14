package uninet

import (
	"net"
	"time"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var (
	readTimeout  = 10 * time.Second
	writeTimeout = 5 * time.Second
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
	n, err := c.Read(b.PlusHeaderBlock())
	if err != nil {
		return logex.Trace(err)
	}

	b.SetLengthPlusHeader(n)
	return nil
}

func (c *tcpConn) WriteBlock(b *utils.Block) error {
	n, err := c.Write(b.PlusHeaderBytes())
	if err != nil {
		return logex.Trace(err)
	}
	if n != b.PlusHeaderLength() {
		return logex.Trace(ErrShortWritten, n, b.PlusHeaderLength())
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
	reply  chan bool
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
func (t *TCPListener) accept(c *tcpConn) {
	var (
		header    = make([]byte, 2)
		reply     = make(chan bool)
		err       error
		n, length int
	)

	for {
		c.SetReadDeadline(time.Now().Add(readTimeout))
		n, err = c.Read(header)
		if err != nil || n != 2 {
			// EOF
			if !inTest {
				logex.Info("close cause by:", err)
			}
			c.Close()
			break
		}

		length = int(utils.ToUint16(header))
		if length > utils.BLOCK_CAP {
			c.Close()
			logex.Info("close")
			break
		}
		t.incomingConn <- &connHeader{c, length, reply}
		if !<-reply {
			c.Close()
			break
		}
	}
}

func (t *TCPListener) Close() error {
	return logex.TraceIfError(t.listener.Close())
}

func (t *TCPListener) ReadBlock(b *utils.Block) (*TcpSession, error) {
	connHeader := <-t.incomingConn
	b.Length = connHeader.header
	connHeader.conn.SetReadDeadline(time.Now().Add(readTimeout))
	n, err := connHeader.conn.Read(b.Bytes())

	if err != nil {
		connHeader.reply <- false
		return nil, logex.Trace(err)
	}

	if n != connHeader.header {
		connHeader.reply <- false
		return nil, logex.Trace(ErrShortRead)
	}
	connHeader.reply <- true

	return NewTcpSession(connHeader.conn.RemoteAddr().(*net.TCPAddr), connHeader.conn), nil
}

func (t *TCPListener) WriteBlock(b *utils.Block, session *TcpSession) error {
	session.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	err := session.conn.WriteBlock(b)
	if err != nil {
		return logex.Trace(err)
	}

	return nil
}

type TCPDialer struct {
	conn *tcpConn
}

func DialTCP(addr *TcpURL) (*TCPDialer, error) {
	c, err := net.Dial(addr.Network(), addr.Host())
	if err != nil {
		return nil, logex.Trace(err)
	}
	conn := c.(*net.TCPConn)
	t := &TCPDialer{
		conn: newTcpConn(conn),
	}
	return t, nil
}

func (t *TCPDialer) ReadBlock(b *utils.Block) error {
	if err := t.conn.ReadBlock(b); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func (t *TCPDialer) WriteBlock(b *utils.Block) error {
	if err := t.conn.WriteBlock(b); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func (t *TCPDialer) Close() error {
	return logex.TraceIfError(t.conn.Close())
}
