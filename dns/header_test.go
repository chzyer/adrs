package dns

import (
	"strconv"
	"strings"
	"testing"

	"gopkg.in/logex.v1"
)

func testGetHeader(i int) []byte {
	testHeader := []string{
		"51 242 1 0 0 1 0 0 0 0 0 0 4 48 120 100 102 3 99 111 109 0 0 1 0 1",
	}
	h := testHeader[i]
	sp := strings.Split(h, " ")
	ret := make([]byte, len(sp))
	var b int
	for i := range sp {
		b, _ = strconv.Atoi(sp[i])
		ret[i] = byte(b)
	}
	return ret
}

func TestHeader(t *testing.T) {
	header := testGetHeader(0)

	h, err := NewHeader(header)
	if err != nil {
		t.Fatal(err)
	}

	if h.ID != 13298 || h.QR != QR_QUERY {
		t.Fatal("parse fail")
	}
	logex.Pretty(h)
}
