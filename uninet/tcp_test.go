package uninet

import (
	"bytes"
	"net"
	"runtime"
	"sync"
	"testing"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

func init() {
	inTest = true
}

func TestTCPListen(t *testing.T) {
	runtime.GOMAXPROCS(4)
	url, err := ParseURL("tcp://:12345")
	if err != nil {
		t.Fatal(err)
	}
	ln, err := ListenTcp(url.(*TcpURL))
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		pool := utils.NewBlockPool()
		for {
			block := pool.Get()
			session, err := ln.ReadBlock(block)
			if err != nil {
				t.Fatal(err)
			}
			if err := ln.WriteBlock(block, session); err != nil {
				logex.Error(err)
				t.Fatal(err)
			}
		}
	}()

	addr := "localhost:12345"

	sources := [][]byte{}
	ex := [][]byte{[]byte("hello\n"), []byte("fuck"), []byte("caocaocao")}
	for i := 0; i < 30; i++ {
		sources = append(sources, ex...)
	}
	var wg sync.WaitGroup
	wg.Add(len(sources))

	for _, source := range sources {
		source := source
		go func() {
			defer wg.Done()

			err := mockTcpConn(addr, source)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	wg.Wait()
}

func mockTcpConn(addr string, source []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return logex.NewTraceError(err)
	}
	defer conn.Close()

	source = append(utils.Uint16To(uint16(len(source))), source...)
	if _, err = conn.Write(source); err != nil {
		return logex.NewTraceError(err)
	}

	b := make([]byte, 512)
	n, err := conn.Read(b)
	if err != nil {
		return err
	}
	if n == 2 {
		n_, err := conn.Read(b[n:])
		if err != nil {
			return err
		}
		n += n_
	}

	if !bytes.Equal(b[:n], source) {
		logex.Info(b[:n], source)
		return logex.NewTraceError("result not except1")
	}

	return nil
}
