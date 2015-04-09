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
	MailChan chan *Mail
}

func NewMailBox(host *uninet.DialAddr) (*MailBox, error) {
	n, err := uninet.NewUniDial(host)
	if err != nil {
		return nil, logex.Trace(err)
	}
	m := &MailBox{
		net:      n,
		MailChan: make(chan *Mail, MailBoxSize),
	}
	return m, nil
}

func (mb *MailBox) Working() {
	var err error
	for mail := range mb.MailChan {
		err = mb.net.WriteBlockUDP(mail.Question.Block())
		if err != nil {

		}
	}
}

func (mb *MailBox) DeliverAndWaitingForReply(m *Mail) error {

	return nil
}
