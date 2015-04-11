package mailman

import (
	"sync"
	"time"

	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type Envelope struct {
	Mail     *Mail
	Customer *customer.Customer
}

type Mail struct {
	From     uninet.Session
	Question *dns.DNSMessage
	Answer   *utils.Block
	reply    chan *utils.Block
}

func (m *Mail) Init() {
	m.From = nil
	m.Question = nil
	m.Answer = nil
}

func (m *Mail) Reply(b *utils.Block) {
	m.reply <- b
}

func (m *Mail) WaitForReply() error {
	select {
	case m.Answer = <-m.reply:
		return nil
	case <-time.After(2 * time.Second):
		return logex.NewTraceError("timeout")
	}
}

type MailPool struct {
	sync.Pool
}

func NewMailPool() *MailPool {
	mp := &MailPool{}
	mp.Pool = sync.Pool{
		New: mp.newMail,
	}
	return mp
}

func (m *MailPool) newMail() interface{} {
	return &Mail{
		reply: make(chan *utils.Block),
	}
}

func (m *MailPool) Get() *Mail {
	return m.Pool.Get().(*Mail)
}

func (mp *MailPool) Put(m *Mail) {
	m.Init()
	mp.Pool.Put(m)
}
