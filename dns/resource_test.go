package dns

import (
	"testing"

	"gopkg.in/logex.v1"
)

type answerMatch struct {
	header   *DNSHeader
	question []*DNSQuestion
	answer   []*DNSResource
	m        string
}

var testAnswer = []answerMatch{
	{
		answer: []*DNSResource{
			{
				Name: []byte{192, 12},
				Type: 1, Class: 1, TTL: 36, RDLength: 4,
				RData: []byte{180, 149, 134, 142},
			},
			{
				Name: []byte{192, 12},
				Type: 1, Class: 1, TTL: 36, RDLength: 4,
				RData: []byte{180, 149, 134, 141},
			},
		},
		m: "211 62 129 128 0 1 0 2 0 0 0 0 5 119 101 105 98 111 3 99 111 109 0 0 1 0 1 192 12 0 1 0 1 0 0 0 36 0 4 180 149 134 142 192 12 0 1 0 1 0 0 0 36 0 4 180 149 134 141",
	},
}

func TestAnswer(t *testing.T) {
	var err error
	for _, q := range testAnswer {
		r := stringToReader(q.m)
		q.header, err = NewDNSHeader(r)
		if err != nil {
			logex.Fatal(err)
		}

		for i := 0; i < int(q.header.QDCount); i++ {
			qq, err := NewDNSQuestion(r)
			if err != nil {
				logex.Fatal(err)
			}
			q.question = append(q.question, qq)
		}

		for i := 0; i < int(q.header.ANCount); i++ {
			answer, err := NewDNSResource(r)
			if err != nil {
				logex.Fatal(err)
			}

			if !answer.Equal(q.answer[i]) {
				logex.Info(answer)
				t.Error("result not except")
			}
		}
	}
}
