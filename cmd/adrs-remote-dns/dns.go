package main

import (
	"flag"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var singleRequestStr = "51 242 1 0 0 1 0 0 0 0 0 0 4 48 120 100 102 3 99 111 109 0 0 1 0 1"

func stringToByte(s string) []byte {
	sp := strings.Split(s, " ")
	ret := make([]byte, len(sp))
	var b int
	for i := range sp {
		b, _ = strconv.Atoi(sp[i])
		ret[i] = byte(b)
	}
	return ret
}

func main() {
	var network string
	var host string
	flag.StringVar(&network, "n", "udp", "network")
	flag.StringVar(&host, "h", "8.8.8.8:53", "host name")
	flag.Parse()

	conn, err := net.Dial(network, host)
	if err != nil {
		logex.Fatal(err)
	}
	defer conn.Close()

	data := stringToByte(singleRequestStr)
	if network == "tcp" {
		data = append(utils.Uint16To(uint16(len(data))), data...)
	}

	go func() {
		for _ = range time.Tick(time.Second) {
			write(conn, data)
		}
	}()
	for {
		read(conn)
	}
}

func write(conn net.Conn, data []byte) {
	n, err := conn.Write(data)
	if err != nil {
		logex.Fatal(err)
	}
	if n != len(data) {
		logex.Fatal("short written")
	}

}

func read(conn net.Conn) {
	ret := make([]byte, 512)
	n, err := conn.Read(ret)
	if err != nil {
		logex.Fatal(err)
	}
	logex.Info(ret[:n])
}
