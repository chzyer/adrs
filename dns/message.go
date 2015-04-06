package dns

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type Message struct {
	Header   *DNSHeader
	Question []*DNSQuestion
}

func NewMessage(b []byte) (*Message, error) {
	var err error
	m := new(Message)
	rr := utils.NewRecordReader(b)
	m.Header, err = NewDNSHeader(rr)
	if err != nil {
		return nil, logex.Trace(err)
	}

	m.Question, err = m.getQuestions(rr, int(m.Header.QDCount))
	if err != nil {
		return nil, logex.Trace(err)
	}

	return m, nil
}

func (m *Message) getQuestions(r *utils.RecordReader, count int) ([]*DNSQuestion, error) {
	var (
		err error
		ret = make([]*DNSQuestion, count)
	)

	for i := 0; i < int(count); i++ {
		ret[i], err = NewDNSQuestion(r)
		if err != nil {
			err = logex.Trace(err)
			return nil, err
		}
	}
	return ret, nil
}
