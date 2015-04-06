package dns

import (
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
	Name []string

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
	qName, err := utils.ReadByFirstByte(r)
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
		RDLength: qRDLength,
		RData:    qRData,
	}, nil
}
