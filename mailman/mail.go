package mailman

import (
	"sync"

	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
)

type Envelope struct {
	Mail     *Mail
	Customer *customer.Customer
}

type Mail struct {
	From     uninet.Session
	Question *dns.DNSMessage
	Answer   *utils.Block
}

func (m *Mail) Init() {
	m.From = nil
	m.Question = nil
	m.Answer = nil
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
	return &Mail{}
}

func (m *MailPool) Get() *Mail {
	return m.Pool.Get().(*Mail)
}

func (mp *MailPool) Put(m *Mail) {
	m.Init()
	mp.Pool.Put(m)
}
