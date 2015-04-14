package dns

import "github.com/chzyer/adrs/utils"

// The second segment of query header section
type DNSHeaderOption struct {
	QR     QR
	OpCode OPCODE
	AA     AA
	TC     TC
	RD     RD
	RA     RA
	Z      Z
	Rcode  RCODE
}

func NewDNSHeaderOption(option uint64) *DNSHeaderOption {
	return &DNSHeaderOption{
		QR:     QR(utils.Read8Bit(option, 15, 1)),
		OpCode: OPCODE(utils.Read8Bit(option, 11, 4)),
		AA:     AA(utils.ReadBit(option, 10)),
		TC:     TC(utils.ReadBit(option, 9)),
		RD:     RD(utils.ReadBit(option, 8)),
		RA:     RA(utils.ReadBit(option, 7)),
		Z:      Z(utils.Read8Bit(option, 4, 3)),
		Rcode:  RCODE(utils.Read8Bit(option, 0, 4)),
	}
}

func (d *DNSHeaderOption) WriteTo(w *utils.RecordWriter) error {
	var data uint16
	data |= uint16(d.QR) << 15
	data |= uint16(d.OpCode) << 11
	data |= uint16(d.AA) << 10
	data |= uint16(d.TC) << 9
	data |= uint16(d.RD) << 8
	data |= uint16(d.RA) << 7
	data |= uint16(d.Z) << 4
	data |= uint16(d.Rcode)

	if err := w.WriteUint16(data); err != nil {
		return err
	}
	return nil
}

// A one bit field that specifies whether this message is a
// query (0), or a response (1).
type QR int

const (
	QR_QUERY QR = 0
	QR_RESP     = 1
)

// A four bit field that specifies kind of query in this
// message.  This value is set by the originator of a query
// and copied into the response.  The values are:
//     0               a standard query (QUERY)
//     1               an inverse query (IQUERY)
//     2               a server status request (STATUS)
//     3-15            reserved for future use
type OPCODE int

const (
	OPCODE_QUERY  OPCODE = 0
	OPCODE_IQUERY        = 1
	OPCODE_STATUS        = 2
)

// Authoritative Answer - this bit is valid in responses,
// and specifies that the responding name server is an
// authority for the domain name in question section.
//
// Note that the contents of the answer section may have
// multiple owner names because of aliases.  The AA bit
// corresponds to the name which matches the query name, or
// the first owner name in the answer section.
type AA uint8

// TrunCation - specifies that this message was truncated
// due to length greater than that permitted on the
// transmission channel.
type TC uint8

// Recursion Desired - this bit may be set in a query and
// is copied into the response.  If RD is set, it directs
// the name server to pursue the query recursively.
// Recursive query support is optional.
type RD uint8

// Recursion Available - this be is set or cleared in a
// response, and denotes whether recursive query support is
// available in the name server.
type RA uint8

// Reserved for future use.  Must be zero in all queries
// and responses.
type Z int

// Response code - this 4 bit field is set as part of
// responses.  The values have the following
// interpretation:
//  0               No error condition
//  1               Format error - The name server was
//                  unable to interpret the query.
//  2               Server failure - The name server was
//                  unable to process this query due to a
//                  problem with the name server.
//  3               Name Error - Meaningful only for
//                  responses from an authoritative name
//                  server, this code signifies that the
//                  domain name referenced in the query does
//                  not exist.
//  4               Not Implemented - The name server does
//                  not support the requested kind of query.
//  5               Refused - The name server refuses to
//                  perform the specified operation for
//                  policy reasons.  For example, a name
//                  server may not wish to provide the
//                  information to the particular requester,
//                  or a name server may not wish to perform
//                  a particular operation (e.g., zone
//                  transfer) for particular data.
//  6-15            Reserved for future use.
type RCODE int

const (
	RCODE_NO_ERR   RCODE = 0
	RCODE_FMT_ERR        = 1
	RCODE_SVR_FAIL       = 2
	RCODE_NAME_ERR       = 3
	RCODE_NOT_IMP        = 4
	RCODE_REFUSED        = 5
	RCODE_RESERVED       = 6
)
