package mailman

import "github.com/chzyer/adrs/uninet"

type MailFmt struct {
	AddrGroup map[string]*AddrGroup
	Router    map[string]string
	Listen    []*uninet.URLer
}

type AddrGroup struct {
	Name string
	Data []*uninet.URLer
}
