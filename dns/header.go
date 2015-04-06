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
	option16, err := rr.ReadUint16()
	if err != nil {
		err = logex.Trace(err)
		return
	}
	option := uint64(option16)
	h.QR = HQRTYPE(utils.Read8Bit(option, 16, 1))
	h.OpCode = HOPCODE(utils.Read8Bit(option, 15, 4))

	return
}
