package mailman

import (
	"github.com/chzyer/adrs/uninet"
	"gopkg.in/logex.v1"
)

var (
	MailBoxSize = 10
	gMailBox    = make(map[string]*MailBox)
)

func GetMailBox(host string) (*MailBox, error) {
	return nil, nil
}

type MailBox struct {
	net      *uninet.UniDial
	mailChan chan *Mail
}

func NewMailBox(host *uninet.DialAddr) (*MailBox, error) {
	n, err := uninet.NewUniDial(host)
	if err != nil {
		return nil, logex.Trace(err)
	}
	m := &MailBox{
		net:      n,
		mailChan: make(chan *Mail, MailBoxSize),
	}
	return m, nil
}

func (mb *MailBox) reader() {
}

func (mb *MailBox) writer() {
	var err error
	for mail := range mb.mailChan {
		err = mb.net.WriteBlockUDP(mail.Question.Block())
		if err != nil {
			logex.Error(err)
		}
	}
}

func (mb *MailBox) DeliverAndWaitingForReply(m *Mail) error {
	if m.To == nil {
		return logex.NewTraceError("mail missing to")
	}
	mb.mailChan <- m
	m.WaitForReply()
	return nil
}
