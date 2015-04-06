package dns

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type DNSHeader struct {
	ID         uint16
	QR         H_QRTYPE
	OpCode     H_OPCODE
	AA         H_AA
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
	h.QR = H_QRTYPE(option >> 15)
	h.OpCode = H_OPCODE(uint8(option>>11) & uint8(15))
	h.AA = AA(uint8(option>>10) & 1)
	return
}
