package mailman

import (
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

func MakeBoxes() (incomingBox, outgoingBox chan *Envelope) {
	incomingBox = make(chan *Envelope, 10)
	outgoingBox = make(chan *Envelope, 10)
	return incomingBox, outgoingBox
}

// determine where the mail to go
type MailMan struct {
	pool        *utils.BlockPool
	incomingBox chan *Envelope
	outgoingBox chan *Envelope
}

func NewMailMan(pool *utils.BlockPool, incomingBox, outgoingBox chan *Envelope) *MailMan {
	mm := &MailMan{
		pool:        pool,
		incomingBox: incomingBox,
		outgoingBox: outgoingBox,
	}

	return mm
}

func (m *MailMan) CheckingOutgoingBox() {
	for {
		envelope := <-m.outgoingBox
		toURL, _ := uninet.ParseURL("udp://8.8.8.8")
		addr := &uninet.DialAddr{
			UDP: toURL.(*uninet.UdpURL),
		}
		mb, err := GetMailBox(addr, m.pool)
		if err != nil {
			logex.Error(err)
			continue
		}

		logex.Info("we found a mailbox")

		mb.Deliver(envelope, m.incomingBox)
	}
}
