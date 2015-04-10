package uninet

import (
	"net/url"
	"strings"

	"gopkg.in/logex.v1"
)

type URL struct {
	Type NetType
	n    *url.URL
}

func ParseURL(u string) (*URL, error) {
	u_, err := url.Parse(u)
	if err != nil {
		return nil, logex.Trace(err)
	}
	url := &URL{
		n: u_,
	}

	if u_.Host == "" && u_.Path != "" {
		u_.Host = u_.Path
		u_.Path = ""
	}

	switch url.Network() {
	case "":
		url.n.Scheme = "udp"
		fallthrough
	case "udp":
		url.Type = NET_UDP
	case "tcp":
		url.Type = NET_TCP
	case "http":
		url.Type = NET_HTTP
	default:
		return nil, logex.NewError("unsupported protocol")
	}

	if !strings.Contains(u_.Host, ":") {
		switch url.Type {
		case NET_UDP, NET_TCP:
			u_.Host = u_.Host + ":53"
		case NET_HTTP:
			u_.Host = u_.Host + ":80"
		}
	}

	return url, nil
}

func (u *URL) Network() string {
	return u.n.Scheme
}

func (u *URL) Host() string {
	return u.n.Host
}
