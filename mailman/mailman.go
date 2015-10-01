package mailman

import (
	"strings"

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
	router      map[string][]uninet.URLer
	incomingBox chan *Envelope
	outgoingBox chan *Envelope
}

func NewMailMan(router map[string][]uninet.URLer, pool *utils.BlockPool, incomingBox, outgoingBox chan *Envelope) *MailMan {
	mm := &MailMan{
		pool:        pool,
		router:      router,
		incomingBox: incomingBox,
		outgoingBox: outgoingBox,
	}

	return mm
}

func (m *MailMan) FigureOutWhereToGo(mail *Mail) (ret []uninet.URLer) {
	host := mail.Content.GetQueryAddr()
	if len(host) == 0 {
		return nil
	}
	segmentSize := len(host)
	for ; segmentSize > 1; segmentSize-- {
		domain := strings.Join(host[len(host)-segmentSize:], ".")
		if segmentSize != len(host) {
			domain = "*." + domain
		}
		if u, ok := m.router[domain]; ok {
			ret = append(ret, u...)
		} else if segmentSize == len(host) {
			if u, ok := m.router["*."+domain]; ok {
				ret = append(ret, u...)
			}
		}
	}

	if len(ret) == 0 {
		u, ok := m.router["*"]
		if ok {
			ret = append(ret, u...)
		}
	}
	return ret
}

func (m *MailMan) CheckingOutgoingBox() {
	for {
		envelope := <-m.outgoingBox
		to := m.FigureOutWhereToGo(envelope.Mail)
		if len(to) == 0 {
			logex.Error("can't not figure where to go", envelope.Mail)
			m.incomingBox <- envelope
			continue
		}
		logex.Info("query domain:", to, envelope.Mail.Content.GetQueryAddrString())

		mbs, err := GetMailBoxs(to, m.pool)
		if err != nil {
			logex.Error(err)
			continue
		}

		logex.Debug("we found a mailbox")

		mbs[0].Deliver(envelope, m.incomingBox)
	}
}
