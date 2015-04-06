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
		n, addr, err := conn.ReadFromUDP(b)
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

		c, err := net.Dial("udp", "8.8.8.8:53")
		if err != nil {
			logex.Fatal(err)
		}
		_, err = c.Write(b[:n])
		if err != nil {
			logex.Fatal(err)
		}

		data := make([]byte, 512)
		n, err = c.Read(data)
		if err != nil {
			logex.Fatal(err)
		}
		c.Close()
		logex.Info(data[:n])
		m2, err := dns.NewDNSMessage(data[:n])
		logex.Pretty(m2)
		logex.Info(m2.Resources[0].Name)
		if err != nil {
			logex.Fatal(err)
		}

		n, err = conn.WriteToUDP(data[:n], addr)
		if err != nil {
			logex.Fatal(err)
		}
		logex.Info("write to dig:", n)
	}
	logex.Info(conn)
}
