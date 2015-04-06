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
