package wiseman

import (
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type WiseMan struct {
	wayIn  customer.Corridor
	wayOut customer.Corridor
}

func NewWiseMan(wayIn, wayOut customer.Corridor) (*WiseMan, error) {
	w := &WiseMan{
		wayIn:  wayIn,
		wayOut: wayOut,
	}
	return w, nil
}

func (w *WiseMan) Serve() {
	var customer *customer.Customer
	for {
		customer = <-w.wayIn
		err := w.Answer(customer)
		if err != nil {
			// oops!, the wise man is passed out!
			logex.Error(err)
			continue
		}
		// say goodbye
		w.wayOut <- customer
	}
}

func (w *WiseMan) Answer(c *customer.Customer) error {
	r := utils.NewRecordReader(c.Question[:c.QuestionN])
	msg, err := dns.NewDNSMessage(r)
	if err != nil {
		return logex.Trace(err)
	}
	c.Msg = msg
	c.Answer = c.Question
	c.AnswerN = c.QuestionN
	logex.Pretty(c)
	return nil
}
