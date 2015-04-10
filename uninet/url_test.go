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
	}
	for _, us := range urls {
		u, err := ParseURL(us.Source)
		if us.IsError && err == nil {
			t.Fatal("result not except")
		} else if us.IsError {
			continue
		} else if err != nil {
			t.Fatal(err)
		}

		if u.Host() != us.Host || u.Network() != us.Network {
			logex.Pretty(u.n)
			t.Fatal("result not except")
		}
	}
}
