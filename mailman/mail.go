package mailman

import (
	"net"
	"sync"

	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/utils"
)

type Mail struct {
	From     *net.UDPAddr
	To       *net.UDPAddr
	Question *dns.DNSMessage
	Answer   *utils.Block
	reply    chan struct{}
}

func (m *Mail) Init() {
	m.From = nil
	m.To = nil
	m.Question = nil
	m.Answer = nil
}

func (m *Mail) Reply() {
	m.reply <- struct{}{}
}

func (m *Mail) WaitForReply() {
	<-m.reply
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
		reply: make(chan struct{}),
	}
}

func (m *MailPool) Get() *Mail {
	return m.Pool.Get().(*Mail)
}

func (mp *MailPool) Put(m *Mail) {
	m.Init()
	mp.Pool.Put(m)
}
