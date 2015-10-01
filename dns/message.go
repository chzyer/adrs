package dns

import (
	"strings"
	"time"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type DNSMessage struct {
	Header    *DNSHeader
	Questions []*DNSQuestion
	Resources []*DNSResource
	block     *utils.Block
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

	m.block = r.Block()
	return m, nil
}

func (m *DNSMessage) Copy(b *utils.Block) *DNSMessage {
	return &DNSMessage{
		Header:    m.Header,
		Questions: m.Questions,
		Resources: m.Resources,
		block:     b,
	}
}

func (m *DNSMessage) GetDeadline() time.Time {
	if len(m.Resources) == 0 {
		return time.Time{}
	}
	return m.Resources[0].Deadline
}

func (m *DNSMessage) WriteTo(w *utils.RecordWriter) (err error) {
	if err = m.Header.WriteTo(w); err != nil {
		return logex.Trace(err)
	}

	for _, q := range m.Questions {
		if err = q.WriteTo(w); err != nil {
			return logex.Trace(err)
		}
	}

	for _, r := range m.Resources {
		if err = r.WriteTo(w); err != nil {
			return logex.Trace(err)
		}
	}
	return
}

func (m *DNSMessage) Key() string {
	qs := m.Questions
	if len(qs) == 0 {
		return ""
	}
	return qs[0].Key()
}

func (m *DNSMessage) GetQueryAddrString() string {
	return strings.Join(m.GetQueryAddr(), ".")
}

func (m *DNSMessage) GetQueryAddr() []string {
	qs := m.Questions
	if len(qs) == 0 {
		return nil
	}
	return qs[0].QName
}

func (m *DNSMessage) Id() uint16 {
	return m.Header.ID
}

// will write to block if block is empty
func (m *DNSMessage) Block() *utils.Block {
	if m.block.Length == 0 {
		err := m.WriteTo(utils.NewRecordWriter(m.block))
		if err != nil {
			logex.Error(err)
		}
	}
	return m.block
}

func (m *DNSMessage) Equal(m2 *DNSMessage) bool {
	if m != nil && m2 == nil || m == nil && m2 != nil {
		return false
	}

	if !m.Header.Equal(m2.Header) {
		return false
	}

	if len(m.Questions) != len(m2.Questions) ||
		len(m.Resources) != len(m2.Resources) {
		return false
	}

	for idx := range m.Questions {
		if !m.Questions[idx].Equal(m2.Questions[idx]) {
			return false
		}
	}

	for idx := range m.Resources {
		if !m.Resources[idx].Equal(m2.Resources[idx]) {
			return false
		}
	}

	return true
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
