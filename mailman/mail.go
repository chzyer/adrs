package mailman

import (
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
)

type Envelope struct {
	Mail     *Mail
	Customer *customer.Customer
}

type Mail struct {
	From    uninet.Session
	Content *dns.DNSMessage
	Reply   *utils.Block
}
