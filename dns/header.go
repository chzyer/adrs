package dns

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

type DNSHeader struct {
	ID     uint16
	Option *DNSHeaderOption

	// specifying the number of entries in the question section.
	QDCount uint16

	// specifying the number of resource records in the answer section.
	ANCount uint16

	// specifying the number of name server resource records
	// in the authority records section.
	NSCount uint16

	// specifying the number of resource records in the additional records section.
	ARCount uint16

	underlying []byte
}

func PeekHeaderID(b *utils.Block) uint16 {
	if b.Length < 2 {
		return 0
	}
	return utils.ToUint16(b.All[:2])
}

func NewDNSHeader(rr *utils.RecordReader) (h *DNSHeader, err error) {
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
	h.Option = NewDNSHeaderOption(uint64(option))

	refs := []*uint16{&h.QDCount, &h.ANCount, &h.NSCount, &h.ARCount}

	for _, ref := range refs {
		*ref, err = rr.ReadUint16()
		if err != nil {
			err = logex.Trace(err)
			return
		}
	}

	return
}

func (h *DNSHeader) Equal(h2 *DNSHeader) bool {
	if h != nil && h2 == nil || h == nil && h2 != nil {
		return false
	}

	return !(h.ID != h2.ID ||
		h.Option.QR != h2.Option.QR ||
		h.Option.OpCode != h2.Option.OpCode ||
		h.Option.AA != h2.Option.AA ||
		h.Option.TC != h2.Option.TC ||
		h.Option.RD != h2.Option.RD ||
		h.Option.RA != h2.Option.RA ||
		h.Option.Z != h2.Option.Z ||
		h.Option.Rcode != h2.Option.Rcode ||
		h.QDCount != h2.QDCount ||
		h.ANCount != h2.ANCount ||
		h.NSCount != h2.NSCount ||
		h.ARCount != h2.ARCount)
}
