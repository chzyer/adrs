package main

import (
	"net"
	"strconv"
	"strings"

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
	conn, err := net.Dial("tcp", "localhost:53")
	if err != nil {
		logex.Fatal(err)
	}
	defer conn.Close()

	data := stringToByte(singleRequestStr)
	data = append(utils.Uint16To(uint16(len(data))), data...)

	for {
		write(conn, data)
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
