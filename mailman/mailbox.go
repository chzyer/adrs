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
	gMailBox      = make(map[*uninet.DialAddr]*MailBox)
	gMailBoxMutex sync.Mutex
)

func GetMailBox(addr *uninet.DialAddr, pool *utils.BlockPool) (*MailBox, error) {
	gMailBoxMutex.Lock()
	defer gMailBoxMutex.Unlock()

	mb, ok := gMailBox[addr]
	if ok {
		return mb, nil
	}

	mb, err := NewMailBox(addr, pool)
	if err != nil {
		return nil, logex.Trace(err)
	}
	gMailBox[addr] = mb
	return mb, nil
}

type MailBox struct {
	net       *uninet.UniDial
	mailChan  chan *Mail
	blockPool *utils.BlockPool
	box       map[uint16]*Mail
	boxLock   sync.Mutex
}

func NewMailBox(host *uninet.DialAddr, pool *utils.BlockPool) (*MailBox, error) {
	n, err := uninet.NewUniDial(host)
	if err != nil {
		return nil, logex.Trace(err)
	}
	mb := &MailBox{
		net:       n,
		mailChan:  make(chan *Mail, MailBoxSize),
		blockPool: pool,
		box:       make(map[uint16]*Mail),
	}

	go mb.writer()
	go mb.reader()
	return mb, nil
}

func (mb *MailBox) reader() {
	var (
		err   error
		block = mb.blockPool.Get()
		mail  *Mail
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
		mail, ok = mb.box[id]
		if ok {
			delete(mb.box, id)
		}
		mb.boxLock.Unlock()
		if !ok {
			logex.Error("block come but don't know where to go")
			block.Init()
			continue
		}

		mail.Reply(block)
		block = mb.blockPool.Get()
	}
}

func (mb *MailBox) writer() {
	var err error
	for mail := range mb.mailChan {
		err = mb.net.WriteBlockUDP(mail.Question.Block())
		if err != nil {
			logex.Error(err)
			continue
		}
		mb.boxLock.Lock()
		mb.box[mail.Question.Id()] = mail
		mb.boxLock.Unlock()
	}
}

func (mb *MailBox) DeliverAndWaitingForReply(m *Mail) error {

	mb.mailChan <- m
	m.WaitForReply()
	if m.Answer == nil {
		return logex.NewError("not got answer")
	}
	return nil
}
