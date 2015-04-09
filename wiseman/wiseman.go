package wiseman

import (
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/mailman"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type WiseMan struct {
	frontDoor customer.Corridor
	backDoor  customer.Corridor
	mailMan   *mailman.MailMan
}

func NewWiseMan(frontDoor, backDoor customer.Corridor) (*WiseMan, error) {
	w := &WiseMan{
		frontDoor: frontDoor,
		backDoor:  backDoor,
	}
	return w, nil
}

func (w *WiseMan) ServeAll() {
	var customer *customer.Customer
	for {
		customer = <-w.frontDoor
		err := w.Serve(customer)
		if err != nil {
			// oops!, the wise man is passed out!
			logex.Error(err)
			continue
		}
		// say goodbye
		w.backDoor <- customer
	}
}

func (w *WiseMan) Serve(c *customer.Customer) error {
	r := utils.NewRecordReader(c.Question)
	msg, err := dns.NewDNSMessage(r)
	if err != nil {
		c.Question.Recycle()
		return logex.Trace(err)
	}
	c.Msg = msg
	// looking up the wikis.
	// ask others

	mail := w.writeMail(c)
	if err := w.mailMan.SendMailAndGotReply(mail); err != nil {
		c.Question.Recycle()
		return logex.Trace(err)
	}

	msg, err = w.readMail(mail)
	if err != nil {
		return logex.Trace(err)
	}
	c.Answer = msg.Block()
	return nil
}

func (w *WiseMan) writeMail(c *customer.Customer) *mailman.Mail {
	return nil
}

func (w *WiseMan) readMail(m *mailman.Mail) (*dns.DNSMessage, error) {
	return nil, nil
}
