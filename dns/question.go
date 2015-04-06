package dns

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

// CLASS fields appear in resource records.
// QCLASS fields appear in the question section of a query.  QCLASS values
// are a superset of CLASS values; every CLASS is a valid QCLASS.
type QCLASS uint16

const (
	_ QCLASS = 0

	// the Internet
	QCLASS_IN = 1

	// the CSNET class (Obsolete - used only for examples in
	// some obsolete RFCs)
	QCLASS_CS = 2

	// the CHAOS class
	QCLASS_CH = 3

	// Hesiod [Dyer 87]
	QCLASS_HS = 4

	// any class
	QCLASS_WILD = 255
)

// QTYPE fields appear in the question part of a query.  QTYPES are a
// superset of TYPEs, hence all TYPEs are valid QTYPEs.
type QTYPE uint16

const (
	_ QTYPE = 0

	// a host address
	QTYPE_A = 1

	// an authoritative name server
	QTYPE_NS = 2

	// a mail destination (Obsolete - use MX)
	QTYPE_MD = 3

	// a mail forwarder (Obsolete - use MX)
	QTYPE_MF = 4

	// the canonical name for an alias
	QTYPE_CNAME = 5

	// marks the start of a zone of authority
	QTYPE_SOA = 6

	// a mailbox domain name (EXPERIMENTAL)
	QTYPE_MB = 7

	// a mail group member (EXPERIMENTAL)
	QTYPE_MG = 8

	// a mail rename domain name (EXPERIMENTAL)
	QTYPE_MR = 9

	// a null RR (EXPERIMENTAL)
	QTYPE_NULL = 10

	// a well known service description
	QTYPE_WKS = 11

	// a domain name pointer
	QTYPE_PTR = 12

	// host information
	QTYPE_HINFO = 13

	// mailbox or mail list information
	QTYPE_MINFO = 14

	// mail exchange
	QTYPE_MX = 15

	// text strings
	QTYPE_TXT = 16

	//A request for a transfer of an entire zone
	QTYPE_AXFR = 252

	//A request for mailbox-related records (MB, MG or MR)
	QTYPE_MAILB = 253

	//A request for mail agent RRs (Obsolete - see MX)
	QTYPE_MAILA = 254

	//A request for all records
	QTYPE_WILD = 255
)

// Question section format
//                                    1  1  1  1  1  1
//      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                     QNAME                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QTYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QCLASS                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//
type DNSQuestion struct {
	// a domain name represented as a sequence of labels, where
	// each label consists of a length octet followed by that
	// number of octets.  The domain name terminates with the
	// zero length octet for the null label of the root.  Note
	// that this field may be an odd number of octets; no
	// padding is used.
	QName []string

	// a two octet code which specifies the type of the query.
	// The values for this field include all codes valid for a
	// TYPE field, together with some more general codes which
	// can match more than one type of RR.
	QType QTYPE

	// a two octet code that specifies the class of the query.
	// For example, the QCLASS field is IN for the Internet.
	QClass QCLASS
}

func NewDNSQuestion(r *utils.RecordReader) (*DNSQuestion, error) {
	var err error

	// read QName
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

	return &DNSQuestion{
		QName:  qName,
		QType:  QTYPE(qType),
		QClass: QCLASS(qClass),
	}, nil
}
