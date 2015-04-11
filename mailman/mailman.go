package mailman

import (
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

// determine where the mail to go
type MailMan struct {
	pool *utils.BlockPool
}

func NewMailMan(pool *utils.BlockPool) *MailMan {
	return &MailMan{
		pool: pool,
	}
}

func (m *MailMan) SendMailAndGotReply(mail *Mail) error {
	toURL, _ := uninet.ParseURL("udp://8.8.8.8")
	addr := &uninet.DialAddr{
		UDP: toURL.(*uninet.UdpURL),
	}
	mb, err := GetMailBox(addr, m.pool)
	if err != nil {
		return logex.Trace(err)
	}

	err = mb.DeliverAndWaitingForReply(mail)
	if err != nil {
		err = logex.Trace(err)
		return err
	}
	return nil
}
