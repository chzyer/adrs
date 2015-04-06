package dns

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type DNSHeader struct {
	ID         uint16
	QR         HQRTYPE
	OpCode     HOPCODE
	AA         HAA
	underlying []byte
}

func NewHeader(b []byte) (h *DNSHeader, err error) {
	rr := utils.NewRecordReader(b)
	h = new(DNSHeader)
	h.ID, err = rr.ReadUint16()
	if err != nil {
		err = logex.Trace(err)
		return
	}
	option, err := rr.ReadUint16()
	if err != nil {
		err = logex.Trace(err)
		return
	}
	h.QR = HQRTYPE(option >> 15)
	h.OpCode = HOPCODE(uint8(option>>11) & uint8(15))

	return
}
