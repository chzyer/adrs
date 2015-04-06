package main

import (
	"net"

	"github.com/chzyer/adrs/dns"
	"gopkg.in/logex.v1"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		logex.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logex.Fatal(err)
	}

	logex.Info("started...")

	b := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			logex.Fatal(err)
		}

		msg, err := dns.NewDNSMessage(b[:n])
		if err != nil {
			logex.Fatal(err)
		}
		logex.Pretty(msg)
		if addr != nil {
			logex.Info(n, b[:n], addr)
		}

	}
	logex.Info(conn)
}
