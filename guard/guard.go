package guard

import (
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/uninet"
	"gopkg.in/logex.v1"
)

type Guard struct {
	ln        *uninet.UniListener
	frontDoor customer.Corridor
	backDoor  customer.Corridor
}

func NewGuard(frontDoor, backDoor customer.Corridor, un *uninet.UniListener) (*Guard, error) {
	g := &Guard{
		ln:        un,
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
	block, session, err := g.ln.ReadBlock()
	if err != nil {
		return nil, logex.Trace(err)
	}

	return &customer.Customer{
		Session: session,
		Raw:     block,
	}, nil
}

func (g *Guard) seeingOffCustomer(c *customer.Customer) (err error) {
	err = g.ln.WriteBlock(c.Raw, c.Session)
	if err != nil {
		return logex.Trace(err)
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
