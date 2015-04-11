package uninet

import (
	"net/url"
	"strings"

	"gopkg.in/logex.v1"
)

type URLer interface {
	Network() string
	Host() string
	NetType() NetType
}

type BaseURL struct {
	netType NetType
	host    string
}
type UdpURL struct{ *BaseURL }
type TcpURL struct{ *BaseURL }
type HttpURL struct{ *BaseURL }

func NewBaseURL(t NetType, host string) *BaseURL {
	return &BaseURL{t, host}
}

func (b *BaseURL) NetType() NetType {
	return b.netType
}
func (b *BaseURL) Network() string {
	return b.netType.String()
}
func (b *BaseURL) Host() string {
	return b.host
}

func ParseURL(u string) (URLer, error) {
	var (
		netType NetType
		host    string
	)

	u_, err := url.Parse(u)
	if err != nil {
		return nil, logex.Trace(err)
	}

	// init host
	host = u_.Host
	if host == "" && u_.Path != "" && !strings.HasPrefix(u_.Path, "/") {
		host = u_.Path
	}
	if host == "" {
		return nil, logex.NewError("missing host")
	}

	// netType
	switch u_.Scheme {
	case "", "udp":
		netType = NET_UDP
	case "tcp":
		netType = NET_TCP
	case "http":
		netType = NET_HTTP
	default:
		return nil, logex.NewError("unsupported protocol")
	}

	if !strings.Contains(host, ":") {
		switch netType {
		case NET_UDP, NET_TCP:
			host = host + ":53"
		case NET_HTTP:
			host = host + ":80"
		}
	}

	base := NewBaseURL(netType, host)

	switch netType {
	case NET_UDP:
		return &UdpURL{base}, nil
	case NET_TCP:
		return &TcpURL{base}, nil
	case NET_HTTP:
		return &HttpURL{base}, nil
	}

	return nil, logex.NewTraceError("unknown error")
}
