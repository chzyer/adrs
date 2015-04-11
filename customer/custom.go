package customer

import (
	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
)

type Corridor chan *Customer
type Customer struct {
	Session uninet.Session
	Raw     *utils.Block
	Msg     *dns.DNSMessage
}

func (c *Customer) SetRaw(b *utils.Block) {
	c.Raw.Recycle()
	c.Raw = b
}

func (c *Customer) LetItGo() {
	if c.Raw != nil {
		c.Raw.Recycle()
	}
	if c.Msg != nil {
		c.Msg.Block().Recycle()
	}
}
