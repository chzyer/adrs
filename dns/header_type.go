package dns

import "github.com/chzyer/adrs/utils"

// The second segment of query header section
type DNSHeaderOption struct {
	QR     H_QR
	OpCode H_OPCODE
	AA     H_AA
	TC     H_TC
	RD     H_RD
	RA     H_RA
	Z      H_Z
	Rcode  H_RCODE
}

func NewDNSHeaderOption(option uint64) *DNSHeaderOption {
	return &DNSHeaderOption{
		QR:     H_QR(utils.Read8Bit(option, 15, 1)),
		OpCode: H_OPCODE(utils.Read8Bit(option, 11, 4)),
		AA:     H_AA(utils.ReadBitBool(option, 10)),
		TC:     H_TC(utils.ReadBitBool(option, 9)),
		RD:     H_RD(utils.ReadBitBool(option, 8)),
		RA:     H_RA(utils.ReadBitBool(option, 7)),
		Z:      H_Z(utils.Read8Bit(option, 4, 3)),
		Rcode:  H_RCODE(utils.Read8Bit(option, 0, 4)),
	}
}

// A one bit field that specifies whether this message is a
// query (0), or a response (1).
type H_QR int

const (
	QR_QUERY H_QR = 0
	QR_RESP       = 1
)

// A four bit field that specifies kind of query in this
// message.  This value is set by the originator of a query
// and copied into the response.  The values are:
//     0               a standard query (QUERY)
//     1               an inverse query (IQUERY)
//     2               a server status request (STATUS)
//     3-15            reserved for future use
type H_OPCODE int

const (
	OPCODE_QUERY  H_OPCODE = 0
	OPCODE_IQUERY          = 1
	OPCODE_STATUS          = 2
)

// Authoritative Answer - this bit is valid in responses,
// and specifies that the responding name server is an
// authority for the domain name in question section.
//
// Note that the contents of the answer section may have
// multiple owner names because of aliases.  The AA bit
// corresponds to the name which matches the query name, or
// the first owner name in the answer section.
type H_AA bool

// TrunCation - specifies that this message was truncated
// due to length greater than that permitted on the
// transmission channel.
type H_TC bool

// Recursion Desired - this bit may be set in a query and
// is copied into the response.  If RD is set, it directs
// the name server to pursue the query recursively.
// Recursive query support is optional.
type H_RD bool

// Recursion Available - this be is set or cleared in a
// response, and denotes whether recursive query support is
// available in the name server.
type H_RA bool

// Reserved for future use.  Must be zero in all queries
// and responses.
type H_Z int

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
type H_RCODE int

const (
	RCODE_NO_ERR   H_RCODE = 0
	RCODE_FMT_ERR          = 1
	RCODE_SVR_FAIL         = 2
	RCODE_NAME_ERR         = 3
	RCODE_NOT_IMP          = 4
	RCODE_REFUSED          = 5
	RCODE_RESERVED         = 6
)
