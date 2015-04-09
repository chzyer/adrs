package mailman

import (
	"net"

	"github.com/chzyer/adrs/dns"
)

type Mail struct {
	From     *net.UDPAddr
	To       *net.UDPAddr
	Question *dns.DNSMessage
	Answer   *dns.DNSMessage
}

// determine where the mail to go
type MailMan struct {
}

func (m *MailMan) SendMailAndGotReply(mail *Mail) error {
	return nil
}
