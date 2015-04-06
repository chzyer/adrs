package dns

// A one bit field that specifies whether this message is a query (0), or a response (1).
type HQRTYPE int

const (
	_ HQRTYPE = 0

	QR_QUERY = 0
	QR_RESP  = 1
)

// A four bit field that specifies kind of query in this
// message.  This value is set by the originator of a query
// and copied into the response.  The values are:
//     0               a standard query (QUERY)
//     1               an inverse query (IQUERY)
//     2               a server status request (STATUS)
//     3-15            reserved for future use
type HOPCODE int

const (
	_ HOPCODE = 0

	OPCODE_QUERY  = 0
	OPCODE_IQUERY = 1
	OPCODE_STATUS = 2
)

// Authoritative Answer - this bit is valid in responses,
// and specifies that the responding name server is an
// authority for the domain name in question section.
//
// Note that the contents of the answer section may have
// multiple owner names because of aliases.  The AA bit
// corresponds to the name which matches the query name, or
// the first owner name in the answer section.
type HAA bool

const (
	_ HAA = false

	AA_NO  = false
	AA_YES = true
)

// TrunCation - specifies that this message was truncated
// due to length greater than that permitted on the
// transmission channel.
type HTC int

const (
	_ HTC = 0
	TC
)
