package dns

import (
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var (
	DNSQuestionSize = 2 * 3
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
	QType uint16

	// a two octet code that specifies the class of the query.
	// For example, the QCLASS field is IN for the Internet.
	QClass uint16
}

func NewDNSQuestion(r *utils.RecordReader) (*DNSQuestion, error) {
	var (
		err     error
		length  uint8
		qName   []string
		segment []byte = make([]byte, 256)
	)

	// read QName
	for {
		length, err = r.ReadByte()
		if err != nil {
			return nil, logex.Trace(err)
		}

		if length == 0 {
			break
		}

		err := r.ReadN(segment, int(length))
		if err != nil {
			return nil, logex.Trace(err)
		}

		qName = append(qName, string(segment[:length]))
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
		QType:  qType,
		QClass: qClass,
	}, nil
}
