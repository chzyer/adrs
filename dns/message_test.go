package dns

import (
	"bytes"
	"testing"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

func init() {
	inTest = true
}

type msgMatch struct {
	msg *DNSMessage
	m   string
}

var testMsg = []msgMatch{
	{
		msg: &DNSMessage{
			Header: &DNSHeader{
				ID: 54078,
				Option: &DNSHeaderOption{
					QR:     1,
					OpCode: 0,
					AA:     0,
					TC:     0,
					RD:     1,
					RA:     1,
					Z:      0,
					Rcode:  0,
				},
				QDCount: 1,
				ANCount: 2,
				NSCount: 0,
				ARCount: 0,
			},
			Questions: []*DNSQuestion{
				{
					QName: []string{
						"weibo",
						"com",
					},
					QType:  1,
					QClass: 1,
				},
			},
			Resources: []*DNSResource{
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
		},
		m: "211 62 129 128 0 1 0 2 0 0 0 0 5 119 101 105 98 111 3 99 111 109 0 0 1 0 1 192 12 0 1 0 1 0 0 0 36 0 4 180 149 134 142 192 12 0 1 0 1 0 0 0 36 0 4 180 149 134 141",
	},
}

func TestMessage(t *testing.T) {
	for _, q := range testMsg {
		by := stringToByte(q.m)
		r := stringToReader(q.m)
		msg, err := NewDNSMessage(r)
		if err != nil {
			t.Fatal(err)
		}

		if !msg.Equal(q.msg) {
			logex.Pretty(msg, q.msg)
			t.Error("result not except")
		}

		b := utils.NewBlockWithByte(make([]byte, 512))
		b.Length = 0
		if err := msg.WriteTo(utils.NewRecordWriter(b)); err != nil {
			logex.Error(err)
			t.Fatal(err)
		}

		if !bytes.Equal(b.Bytes(), by) {
			logex.Error(b.Bytes())
			logex.Error(by)
			t.Fatal("result not except")
		}
	}
}
