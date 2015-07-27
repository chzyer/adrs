package uninet

import (
	"bytes"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/chzyer/adrs/utils"
	"gopkg.in/logex.v1"
)

var (
	testHost = "localhost:4324"
)

func TestUDPDial(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", testHost)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	go func() {
		b := make([]byte, 512)
		for {
			n, addr, err := conn.ReadFromUDP(b)
			if err != nil {
				if IsClosedError(err) {
					break
				}
				logex.Error(err)
				t.Fatal(err)

			}
			n, err = conn.WriteToUDP(b[:n], addr)
			if err != nil {
				logex.Error(err)
				t.Fatal(err)
			}
		}
	}()

	url, err := ParseURL("udp://" + testHost)
	if err != nil {
		t.Fatal(err)
	}

	sources := [][]byte{
		[]byte("hello\n"), []byte("fuck"), []byte("caocaocao"),
	}
	a := sources
	for i := 0; i < 50; i++ {
		sources = append(sources, a...)
	}
	var wg sync.WaitGroup
	for _, source := range sources {
		wg.Add(1)
		source := source
		go func() {
			defer wg.Done()
			err := testDial(url.(*UdpURL), source)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}
	wg.Wait()
}

func testDial(u *UdpURL, source []byte) error {
	dialer, err := DialUDP(u)
	if err != nil {
		return err
	}
	defer dialer.Close()

	pool := utils.NewBlockPool()
	block := pool.Get()
	block.Write(source)
	if err := dialer.WriteBlock(block); err != nil {
		return err
	}
	block = pool.Get()
	if err := dialer.ReadBlock(block); err != nil {
		return err
	}
	if !bytes.Equal(source, block.Bytes()) {
		return logex.NewError("result not except")
	}
	return nil
}

func TestUDPListen(t *testing.T) {
	runtime.GOMAXPROCS(4)
	url, err := ParseURL("udp://" + testHost)
	if err != nil {
		t.Fatal(err)
	}
	ln, err := ListenUDP(url.(*UdpURL))
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
				if IsClosedError(err) {
					break
				}
				t.Fatal(err)
			}
			if err := ln.WriteBlock(block, session); err != nil {
				t.Fatal(err)
			}
			block.Recycle()
		}
	}()

	sources := [][]byte{
		[]byte("hello\n"), []byte("fuck"), []byte("caocaocao"),
	}
	ex := sources
	for i := 0; i < 10; i++ {
		sources = append(sources, ex...)
	}
	var wg sync.WaitGroup
	wg.Add(len(sources))

	for _, source := range sources {
		source := source
		go func() {
			defer wg.Done()
			err := mockUdpConn(testHost, source)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

	wg.Wait()
}

func mockUdpConn(host string, source []byte) error {
	addr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return logex.NewError(err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second))

	n, err := conn.Write(source)
	if err != nil {
		return logex.NewError(err)
	}

	b := make([]byte, 512)
	n, err = conn.Read(b)
	if err != nil {
		return err
	}

	if !bytes.Equal(b[:n], source) {
		return logex.NewError("result not except")
	}
	return nil
}
