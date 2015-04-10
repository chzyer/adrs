package dns

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type headerMatch struct {
	DNSHeader
	data []byte
	m    string
}

var testHeader = []headerMatch{
	{
		DNSHeader: DNSHeader{
			ID: 13298,
			Option: &DNSHeaderOption{
				QR:     QR_QUERY,
				OpCode: 0,
				RD:     true,
			},
			QDCount: 1,
		},
		m: "51 242 1 0 0 1 0 0 0 0 0 0 4 48 120 100 102 3 99 111 109 0 0 1 0 1",
	},
	{
		// 36100
		DNSHeader: DNSHeader{
			ID: 13298,
			Option: &DNSHeaderOption{
				QR:     QR_RESP,
				OpCode: 1,
				AA:     true,
				RD:     true,
				RA:     true,
				Rcode:  RCODE_NOT_IMP,
			},
			QDCount: 1,
		},
		m: "51 242 141 132 0 1 0 0 0 0 0 0 4 48 120 100 102 3 99 111 109 0 0 1 0 1",
	},
	/*
		{
			m: "215 91 1 0 0 1 0 0 0 0 0 0 4 49 49 49 49 4 50 50 50 50 0 0 1 0 1",
		},
	*/
}

func stringToReader(s string) *utils.RecordReader {
	return utils.NewRecordReader(utils.NewBlockWithByte(stringToByte(s)))
}

func stringToByte(s string) []byte {
	sp := strings.Split(s, " ")
	ret := make([]byte, len(sp))
	var b int
	for i := range sp {
		b, _ = strconv.Atoi(sp[i])
		ret[i] = byte(b)
	}
	return ret
}

func testGetHeader(i int) headerMatch {
	h := testHeader[i]
	h.data = stringToByte(h.m)
	return h
}

func TestPeek(t *testing.T) {
	block := utils.NewBlockWithByte([]byte{1})
	if PeekHeaderID(block) != 0 {
		t.Fatal("excepting error")
	}
}

func TestHeader(t *testing.T) {
	for i := 0; i < len(testHeader); i++ {
		header := testGetHeader(i)
		block := utils.NewBlockWithByte(header.data)

		h, err := NewDNSHeader(utils.NewRecordReader(block))
		if err != nil {
			t.Fatal(err)
		}

		if h.ID != PeekHeaderID(block) {
			t.Fatal("peek not except")
		}

		if h.Equal(nil) {
			t.Fatal("result not except")
		}
		if !h.Equal(&header.DNSHeader) {
			logex.Pretty(h, header)
			t.Fatal("parse fail")
		}
	}

}
