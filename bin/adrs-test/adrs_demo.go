package main

import (
	"net"
	"runtime"

	"github.com/chzyer/adrs/dns"
	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

func main() {
	runtime.GOMAXPROCS(2)
	addr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		logex.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logex.Fatal(err)
	}
	conn.SetReadBuffer(512)

	logex.Info("started...")

	for {
		b := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(b)
		if err != nil {
			logex.Fatal(err)
		}
		go func() {
			r := utils.NewRecordReader(b[:n])
			msg, err := dns.NewDNSMessage(r)
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

			r2 := utils.NewRecordReader(data[:n])
			m2, err := dns.NewDNSMessage(r2)
			logex.Pretty(m2)
			if err != nil {
				logex.Fatal(err)
			}

			_, err = conn.WriteToUDP(data[:n], addr)
			if err != nil {
				logex.Fatal(err)
			}
			logex.Info("write to", addr)
		}()
	}
	logex.Info(conn)
}
