package conf

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/logex.v1"
)

var data = `
[Listen]
tcp: 0.0.0.0
udp: 0.0.0.0
http: 0.0.0.0

[AddrGroup]
google:
	tcp://8.8.8.8
	tcp://8.8.4.4
default:
	114.114.114.114
	8.8.8.8

[DnsGroup]
192.168.1.1
192.168.1.2

[Router]
google:
	*.google.com
	*.facebook.com
	*.twitter.com
default:
	*
`

func TestParser(t *testing.T) {
	conf := new(RawConfig)
	err := Parse(conf, strings.NewReader(data))
	if err != nil {
		logex.Error(err)
		t.Fatal(err)
	}

	if !reflect.DeepEqual(conf.Listen, map[string]string{
		"http": "0.0.0.0",
		"tcp":  "0.0.0.0",
		"udp":  "0.0.0.0",
	}) {
		t.Fatal(conf.Listen)
	}

	if !reflect.DeepEqual(conf.AddrGroup, map[string][]string{
		"default": {"114.114.114.114", "8.8.8.8"},
		"google":  {"tcp://8.8.8.8", "tcp://8.8.4.4"},
	}) {
		t.Fatal(conf.AddrGroup)
	}

	if !reflect.DeepEqual(conf.DnsGroup, []string{
		"192.168.1.1", "192.168.1.2",
	}) {
		t.Fatal(conf.DnsGroup)
	}

	if !reflect.DeepEqual(conf.Router, map[string][]string{
		"default": {"*"},
		"google":  {"*.google.com", "*.facebook.com", "*.twitter.com"},
	}) {
		t.Fatal(conf.Router)
	}
}
