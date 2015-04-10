package guard

import (
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type Guard struct {
	pool      *utils.BlockPool
	net       *uninet.UniNet
	frontDoor customer.Corridor
	backDoor  customer.Corridor
}

func NewGuard(frontDoor, backDoor customer.Corridor, un *uninet.UniNet, pool *utils.BlockPool) (*Guard, error) {
	g := &Guard{
		pool:      pool,
		net:       un,
		frontDoor: frontDoor,
		backDoor:  backDoor,
	}
	return g, nil
}

func (g *Guard) Start() {
	go g.LeadingIn()
	go g.LeadingOut()
}

func (g *Guard) waitingCustomer() (*customer.Customer, error) {
	b := g.pool.Get()
	addr, err := g.net.ReadBlockUDP(b)
	if err != nil {
		b.Recycle()
		return nil, logex.Trace(err)
	}

	return &customer.Customer{
		Type:     uninet.NET_UDP,
		From:     addr,
		Question: b,
	}, nil
}

func (g *Guard) seeingOffCustomer(c *customer.Customer) (err error) {
	switch c.Type {
	case uninet.NET_UDP:
		err = g.net.WriteBlockUDP(c.Answer.Block(), c.From)
		if err != nil {
			logex.Trace(err)
		}
	default:
		err = logex.NewError("type not supported yet!", c)
	}

	c.LetItGo()
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

		g.frontDoor <- customer
	}
}

func (g *Guard) LeadingOut() {
	var (
		customer *customer.Customer
	)
	for {
		customer = <-g.backDoor
		g.seeingOffCustomer(customer)
	}
}