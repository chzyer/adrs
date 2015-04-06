package dns

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type DNSMessage struct {
	Header    *DNSHeader
	Questions []*DNSQuestion
	Resources []*DNSResource
}

func NewDNSMessage(r *utils.RecordReader) (*DNSMessage, error) {
	var err error
	m := new(DNSMessage)
	m.Header, err = NewDNSHeader(r)
	if err != nil {
		return nil, logex.Trace(err)
	}

	m.Questions, err = m.getQuestions(r, int(m.Header.QDCount))
	if err != nil {
		return nil, logex.Trace(err)
	}

	m.Resources, err = m.getResources(r, int(m.Header.ANCount))
	if err != nil {
		return nil, logex.Trace(err)
	}

	return m, nil
}

func (m *DNSMessage) getQuestions(r *utils.RecordReader, count int) ([]*DNSQuestion, error) {
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

func (m *DNSMessage) getResources(r *utils.RecordReader, count int) ([]*DNSResource, error) {
	var (
		err error
		ret = make([]*DNSResource, count)
	)

	for i := 0; i < int(count); i++ {
		ret[i], err = NewDNSResource(r)
		if err != nil {
			err = logex.Trace(err)
			return nil, err
		}

	}
	return ret, nil
}
