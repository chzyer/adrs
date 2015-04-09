package customer

import (
	"net"

	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
)

type Corridor chan *Customer
type Customer struct {
	Type     uninet.NetType
	From     *net.UDPAddr
	Question *utils.Block
	Answer   *dns.DNSMessage
}

func (c *Customer) Recycle() {
	if c.Question != nil {
		c.Question.Recycle()
	}
	if c.Answer != nil {
		c.Answer.Block().Recycle()
	}
}
