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

func GetMailBox(addr *uninet.DialAddr, pool *utils.BlockPool) (*MailBox, error) {
	gMailBoxMutex.Lock()
	defer gMailBoxMutex.Unlock()

	mb, ok := gMailBox[addr.String()]
	if ok {
		return mb, nil
	}

	mb, err := NewMailBox(addr, pool)
	if err != nil {
		return nil, logex.Trace(err)
	}
	gMailBox[addr.String()] = mb
	return mb, nil
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

func NewMailBox(host *uninet.DialAddr, pool *utils.BlockPool) (*MailBox, error) {
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
		if err = mb.net.UDP.ReadBlock(block); err != nil {
			logex.Error(err)
			continue
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

		pkg.envelope.Mail.Answer = block
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
		err = mb.net.WriteBlockUDP(mail.Question.Block())
		if err != nil {
			logex.Error(err)
			continue
		}
		mb.boxLock.Lock()
		mb.box[mail.Question.Id()] = pkg
		mb.boxLock.Unlock()
	}
}

func (mb *MailBox) Deliver(envelope *Envelope, reply chan *Envelope) {
	mb.pkgChan <- &EnvelopePackage{
		envelope: envelope,
		reply:    reply,
	}
}
