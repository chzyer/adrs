package conf

import (
	"io"

	"github.com/chzyer/adrs/uninet"
	"gopkg.in/logex.v1"
)

type RawConfig struct {
	TTLRatio     int                 `json:"TTLRatio"`
	DisableCache string              `json:"DisableCache"`
	Listen       map[string]string   `json:"Listen,omitempty"`
	AddrGroup    map[string][]string `json:"AddrGroup,omitempty"`
	DnsGroup     []string            `json:"DnsGroup,omitempty"`
	Router       map[string][]string `json:"Router,omitempty"`
}

type Config struct {
	TTLRatio     int
	Listen       *uninet.ListenConfig
	AddrGroup    map[string][]uninet.URLer
	Router       map[string][]uninet.URLer // map[domain] []server
	DisableCache bool
}

func NewConfig(r io.Reader) (*Config, error) {
	raw := new(RawConfig)
	if err := Parse(raw, r); err != nil {
		return nil, logex.Trace(err)
	}

	listenConfig, err := ConvListenConfig(raw.Listen)
	if err != nil {
		return nil, logex.Trace(err)
	}
	addrGroup, err := ConvAddrGroup(raw.AddrGroup)
	if err != nil {
		return nil, logex.Trace(err)
	}
	router, err := ConvRouter(raw.Router, addrGroup)
	if err != nil {
		return nil, logex.Trace(err)
	}

	c := &Config{
		Listen:       listenConfig,
		AddrGroup:    addrGroup,
		Router:       router,
		DisableCache: raw.DisableCache == "true",
		TTLRatio:     raw.TTLRatio,
	}
	return c, nil

}

func ConvAddrGroup(raw map[string][]string) (map[string][]uninet.URLer, error) {
	addrGroup := make(map[string][]uninet.URLer)
	for k, vs := range raw {
		for _, v := range vs {
			u, err := uninet.ParseURL(v)
			if err != nil {
				return nil, logex.Trace(err, v)
			}
			addrGroup[k] = append(addrGroup[k], u)
		}
	}
	return addrGroup, nil
}

func ConvRouter(raw map[string][]string, group map[string][]uninet.URLer) (map[string][]uninet.URLer, error) {
	router := make(map[string][]uninet.URLer)
	var ok bool
	for groupName, domains := range raw {
		for _, domain := range domains {
			router[domain], ok = group[groupName]
			if !ok {
				return nil, logex.NewTraceError("group name '", groupName, "' not found")
			}
		}
	}
	return router, nil
}

func ConvListenConfig(raw map[string]string) (*uninet.ListenConfig, error) {
	listenConfig := new(uninet.ListenConfig)
	if tcp := raw["tcp"]; tcp != "" {
		if tcpUrl, err := uninet.ParseURLEx(tcp, uninet.NET_TCP); err != nil {
			return nil, logex.Trace(err, "parse listen.tcp")
		} else {
			listenConfig.Tcp = tcpUrl.(*uninet.TcpURL)
		}
	}

	if udp := raw["udp"]; udp != "" {
		if udpUrl, err := uninet.ParseURLEx(udp, uninet.NET_UDP); err != nil {
			return nil, logex.Trace(err, "parse listen.udp")
		} else {
			listenConfig.Udp = udpUrl.(*uninet.UdpURL)
		}
	}

	if http := raw["http"]; http != "" {
		if httpUrl, err := uninet.ParseURLEx(http, uninet.NET_HTTP); err != nil {
			return nil, logex.Trace(err, "parse listen.udp")
		} else {
			listenConfig.Http = httpUrl.(*uninet.HttpURL)
		}
	}

	return listenConfig, nil
}
