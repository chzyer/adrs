package uninet

import (
	"testing"

	"gopkg.in/logex.v1"
)

type Exception struct {
	Source  string
	Host    string
	Network string
	IsError bool
}

func TestURL(t *testing.T) {
	urls := []Exception{
		{"tcp://localhost:101", "localhost:101", "tcp", false},
		{"udp://localhost", "localhost:53", "udp", false},
		{"localhost", "localhost:53", "udp", false},
		{"http://localhost", "localhost:80", "http", false},
		{"ftp://localhost", "", "", true},
		{"udp:///hello", "", "", true},
		{"ud-p:///hello", "", "", true},
		{"udp://:53", ":53", "udp", false},
		{":53", "", "", true},
	}
	for _, us := range urls {
		u, err := ParseURL(us.Source)
		if us.IsError && err == nil {
			logex.Pretty(u)
			t.Fatal("excepting error")
		} else if us.IsError {
			continue
		} else if err != nil {
			logex.Error(err)
			t.Fatal(err)
		}

		if u.Host() != us.Host || u.Network() != us.Network {
			logex.Pretty(u, us.Host)
			t.Fatal("result not except")
		}
	}
}
