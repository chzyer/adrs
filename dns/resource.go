package dns

import (
	"bytes"
	"time"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

// The answer, authority, and additional sections all share the same
// format: a variable number of resource records, where the number of
// records is specified in the corresponding count field in the header.
// Each resource record has the following format:
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                                               /
//     /                      NAME                     /
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     CLASS                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TTL                      |
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                   RDLENGTH                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
//     /                     RDATA                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type DNSResource struct {
	// a domain name to which this resource record pertains.
	Name []byte

	// Two octets containing one of the RR type codes.  This
	// field specifies the meaning of the data in the RDATA
	// field.
	Type uint16

	// Two octets which specify the class of the data in the
	// RDATA field.
	Class uint16

	// A 32 bit unsigned integer that specifies the time
	// interval (in seconds) that the resource record may be
	// cached before it should be discarded.	Zero values are
	// interpreted to mean that the RR can only be used for the
	// transaction in progress, and should not be cached.
	TTL uint32

	// Calculate by TTL
	Deadline time.Time

	// an unsigned 16 bit integer that specifies the length in
	// octets of the RDATA field.
	RDLength uint16

	// a variable length string of octets that describes the
	// resource. The format of this information varies
	// according to the TYPE and CLASS of the resource record.
	// For example, the if the TYPE is A and the CLASS is IN,
	// the RDATA field is a 4 octet ARPA Internet address.
	RData []byte
}

func NewDNSResource(r *utils.RecordReader) (*DNSResource, error) {
	qName, err := r.ReadBytes(2)
	if err != nil {
		return nil, logex.Trace(err)
	}

	qType, err := r.ReadUint16()
	if err != nil {
		return nil, logex.Trace(err)
	}

	qClass, err := r.ReadUint16()
	if err != nil {
		return nil, logex.Trace(err)
	}

	qTTL, err := r.ReadUint32()
	if err != nil {
		return nil, logex.Trace(err)
	}
	deadline := utils.Now().Add(time.Duration(qTTL) * time.Second)

	qRDLength, err := r.ReadUint16()
	if err != nil {
		return nil, logex.Trace(err)
	}

	qRData, err := r.ReadBytes(int(qRDLength))
	if err != nil {
		return nil, logex.Trace(err)
	}

	// RDATA NOT IMPLEMENT!
	return &DNSResource{
		Name:     qName,
		Type:     qType,
		Class:    qClass,
		TTL:      qTTL,
		Deadline: deadline,
		RDLength: qRDLength,
		RData:    qRData,
	}, nil
}

func (r *DNSResource) WriteTo(w *utils.RecordWriter) (err error) {
	if err = w.WriteSafe(r.Name); err != nil {
		return logex.Trace(err)
	}

	if err = w.WriteUint16(r.Type); err != nil {
		return logex.Trace(err)
	}

	if err = w.WriteUint16(r.Class); err != nil {
		return logex.Trace(err)
	}

	ttl := uint32(r.Deadline.Sub(utils.Now()).Seconds())
	if err = w.WriteUint32(ttl); err != nil {
		return logex.Trace(err)
	}

	if err = w.WriteUint16(r.RDLength); err != nil {
		return logex.Trace(err)
	}

	if err = w.WriteSafe(r.RData); err != nil {
		return logex.Trace(err)
	}
	return
}

func (r *DNSResource) Equal(r2 *DNSResource) bool {
	if r != nil && r2 == nil || r == nil && r2 != nil {
		return false
	}

	return !(bytes.Equal(r.Name, r2.Name) &&
		r.Class != r2.Class ||
		r.Type != r2.Type ||
		r.TTL != r2.TTL ||
		!bytes.Equal(r.RData, r2.RData))
}
