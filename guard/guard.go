package guard

import (
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type Guard struct {
	pool   *utils.BlockPool
	net    *uninet.UniNet
	wayIn  customer.Corridor
	wayOut customer.Corridor
}

func NewGuard(in, out customer.Corridor, un *uninet.UniNet, pool *utils.BlockPool) (*Guard, error) {
	g := &Guard{
		pool:   pool,
		net:    un,
		wayIn:  in,
		wayOut: out,
	}
	return g, nil
}

func (g *Guard) Start() {
	go g.LeadingIn()
	go g.LeadingOut()
}

func (g *Guard) waitingCustomer() (*customer.Customer, error) {
	n, b, addr, err := g.net.ReadBlockUDP()
	if err != nil {
		return nil, logex.Trace(err)
	}
	return &customer.Customer{
		Type:      uninet.NET_UDP,
		From:      addr,
		Question:  b,
		QuestionN: n,
	}, nil
}

func (g *Guard) seeingOffCustomer(c *customer.Customer) (err error) {
	switch c.Type {
	case uninet.NET_UDP:
		err = g.net.WriteBlockUDP(c.Answer[:c.AnswerN], c.From)
		if err != nil {
			logex.Trace(err)
		}
	default:
		err = logex.NewError("type not supported yet!", c)
	}

	c.Recycle(g.pool)
	return
}

func (g *Guard) LeadingIn() {
	var (
		customer *customer.Customer
		err      error
	)
	for {
		customer, err = g.waitingCustomer()
		if err != nil {
			logex.Error(err)
			continue
		}

		g.wayIn <- customer
	}
}

func (g *Guard) LeadingOut() {
	var (
		customer *customer.Customer
	)
	for {
		customer = <-g.wayOut
		g.seeingOffCustomer(customer)
	}
}
