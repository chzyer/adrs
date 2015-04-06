package dns

import (
	"testing"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type questionMatch struct {
	header   *DNSHeader
	question *DNSQuestion
	m        string
}

var testQuestion = []questionMatch{
	{
		question: &DNSQuestion{
			QName:  []string{"0xdf", "com"},
			QType:  1,
			QClass: 1,
		},
		m: "51 242 1 0 0 1 0 0 0 0 0 0 4 48 120 100 102 3 99 111 109 0 0 1 0 1",
	},
}

func TestQuestion(t *testing.T) {
	var err error
	var question *DNSQuestion
	for _, q := range testQuestion {
		r := stringToReader(q.m)
		q.header, err = NewDNSHeader(r)
		if err != nil {
			logex.Fatal(err)
		}
		question, err = NewDNSQuestion(r)
		if err != nil {
			logex.Fatal(err)
		}

		if !utils.CmpString(question.QName, q.question.QName) ||
			question.QType != q.question.QType ||
			question.QClass != q.question.QClass ||
			false {

			logex.Pretty(question, q.question)
			t.Fatal("result not except!")
		}
	}
}
