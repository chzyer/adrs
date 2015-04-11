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
	mailPool  *mailman.MailPool
}

func NewWiseMan(frontDoor, backDoor customer.Corridor, mailMan *mailman.MailMan, mailPool *mailman.MailPool) (*WiseMan, error) {
	w := &WiseMan{
		frontDoor: frontDoor,
		backDoor:  backDoor,
		mailPool:  mailPool,
		mailMan:   mailMan,
	}
	return w, nil
}

func (w *WiseMan) ServeAll() {
	var customer *customer.Customer
	for {
		customer = <-w.frontDoor
		err := w.serve(customer)
		if err != nil {
			// oops!, the wise man is passed out!
			logex.Error(err)
			customer.LetItGo()
			continue
		}
		// say goodbye
		w.backDoor <- customer
	}
}

func (w *WiseMan) serve(c *customer.Customer) error {
	r := utils.NewRecordReader(c.Raw)
	msg, err := dns.NewDNSMessage(r)
	if err != nil {
		return logex.Trace(err)
	}
	// looking up the wikis.
	// ask others

	// we don't know where to send yet
	mail := w.writeMail(c, msg)

	// but don't worry, my mail man know
	if err := w.mailMan.SendMailAndGotReply(mail); err != nil {
		return logex.Trace(err)
	}

	if mail.Answer == nil {
		return logex.NewError("oops!")
	}

	// write to wiki in case someone ask the same question
	c.Raw = mail.Answer
	msg, err = w.readMailAndDestory(mail)
	if err != nil {
		return logex.Trace(err)
	}
	c.Msg = msg

	return nil
}

func (w *WiseMan) writeMail(c *customer.Customer, msg *dns.DNSMessage) *mailman.Mail {
	mail := w.mailPool.Get()
	mail.From = c.Session
	mail.Question = msg
	return mail
}

func (w *WiseMan) readMailAndDestory(m *mailman.Mail) (*dns.DNSMessage, error) {
	msg, err := dns.NewDNSMessage(utils.NewRecordReader(m.Answer))
	w.mailPool.Put(m)
	if err != nil {
		return nil, logex.Trace(err)
	}

	return msg, nil
}
