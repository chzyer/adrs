package mailman

import (
	"sync"

	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var (
	MailBoxSize   = 10
	gMailBox      = make(map[string]*MailBox)
	gMailBoxMutex sync.Mutex
)

func GetMailBoxs(addrs []uninet.URLer, pool *utils.BlockPool) ([]*MailBox, error) {
	gMailBoxMutex.Lock()
	defer gMailBoxMutex.Unlock()
	var err error

	var mbs []*MailBox
	for _, addr := range addrs {
		mb, ok := gMailBox[addr.String()]
		if !ok {
			mb, err = NewMailBox(addr, pool)
			if err != nil {
				return nil, logex.Trace(err)
			}
			gMailBox[addr.String()] = mb
		}

		mbs = append(mbs, mb)
	}

	return mbs, nil
}

type EnvelopePackage struct {
	envelope *Envelope
	reply    chan *Envelope
}

type MailBox struct {
	net       *uninet.UniDial
	pkgChan   chan *EnvelopePackage
	blockPool *utils.BlockPool
	box       map[uint16]*EnvelopePackage
	boxLock   sync.Mutex
}

func NewMailBox(host uninet.URLer, pool *utils.BlockPool) (*MailBox, error) {
	n, err := uninet.NewUniDial(host)
	if err != nil {
		return nil, logex.Trace(err)
	}
	mb := &MailBox{
		net:       n,
		blockPool: pool,
		pkgChan:   make(chan *EnvelopePackage, MailBoxSize),
		box:       make(map[uint16]*EnvelopePackage),
	}

	go mb.writer()
	go mb.reader()
	return mb, nil
}

func (mb *MailBox) reader() {
	var (
		err   error
		block = mb.blockPool.Get()
		pkg   *EnvelopePackage
		id    uint16
		ok    bool
	)

	for {
		if err = mb.net.ReadBlock(block); err != nil {
			// remote side close the connection
			logex.Error(err)
			break
		}

		id = dns.PeekHeaderID(block)
		mb.boxLock.Lock()
		pkg, ok = mb.box[id]
		if ok {
			delete(mb.box, id)
		}
		mb.boxLock.Unlock()
		if !ok {
			logex.Error("block come but don't know where to go")
			block.Init()
			continue
		}

		pkg.envelope.Mail.Reply = block
		pkg.reply <- pkg.envelope
		block = mb.blockPool.Get()
	}
}

func (mb *MailBox) writer() {
	var (
		err  error
		mail *Mail
	)
	for pkg := range mb.pkgChan {
		mail = pkg.envelope.Mail
		err = mb.net.WriteBlock(mail.Content.Block())
		if err != nil {
			logex.Error(err)
			continue
		}
		mb.boxLock.Lock()
		mb.box[mail.Content.Id()] = pkg
		mb.boxLock.Unlock()
	}
}

func (mb *MailBox) Deliver(envelope *Envelope, reply chan *Envelope) {
	mb.pkgChan <- &EnvelopePackage{
		envelope: envelope,
		reply:    reply,
	}
}
