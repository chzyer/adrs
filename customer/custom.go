package customer

import (
	"net"

	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
)

type Corridor chan *Customer
type Customer struct {
	Type      uninet.NetType
	From      *net.UDPAddr
	Question  []byte
	QuestionN int
	Msg       *dns.DNSMessage
	Answer    []byte
	AnswerN   int
}

func (c *Customer) Recycle(p utils.BlockPooler) {
	if c.Question != nil {
		p.Put(c.Question)
	}
	if c.Answer != nil {
		p.Put(c.Answer)
	}
}
