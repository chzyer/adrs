package main

import (
	"net"

	"gopkg.in/logex.v1"
)

type Data struct {
	Addr *net.UDPAddr
	Data []byte
}

func main() {
	addr, _ := net.ResolveUDPAddr("udp", "localhost:5003")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logex.Fatal(err)
	}
	defer conn.Close()
	conn.SetReadBuffer(512)
	conn.SetWriteBuffer(512)
	ch := make(chan Data, 1)

	go func() {
		for {
			ret := <-ch
			n, err := conn.WriteToUDP(ret.Data, ret.Addr)
			if err != nil {
				logex.Fatal(err)
			}
			logex.Info("write", ret.Data, ret.Addr, n)
		}
	}()

	for {
		data := make([]byte, 512)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			logex.Fatal(err)
		}
		logex.Info("recv", data[:n], addr, n)
		ch <- Data{addr, data[:n]}
	}
}
