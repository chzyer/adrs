package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/guard"
	"github.com/chzyer/adrs/mailman"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"github.com/chzyer/adrs/wiseman"
	"gopkg.in/logex.v1"
)

func main() {
	var (
		mailManNum = 1
		wiseManNum = 1
	)

	frontCorridor := make(customer.Corridor)
	backCorridor := make(customer.Corridor)
	udpListen, err := uninet.ParseURL("udp://:53")
	if err != nil {
		logex.Fatal(err)
	}
	tcpListen, err := uninet.ParseURL("tcp://:53")
	if err != nil {
		logex.Fatal(err)
	}

	pool := utils.NewBlockPool()
	incomingBox, outgoingBox := mailman.MakeBoxes()

	un, err := uninet.NewUniListener(tcpListen.(*uninet.TcpURL), udpListen.(*uninet.UdpURL), pool)
	if err != nil {
		logex.Fatal(err)
	}

	// fool guard
	g, err := guard.NewGuard(frontCorridor, backCorridor, un)
	if err != nil {
		logex.Fatal(err)
	}
	g.Start()

	// wiseman
	for i := 0; i < wiseManNum; i++ {
		wm, err := wiseman.NewWiseMan(frontCorridor, backCorridor, incomingBox, outgoingBox)
		if err != nil {
			logex.Fatal(err)
		}
		go wm.ServeAll()
	}

	// mailman
	for i := 0; i < mailManNum; i++ {
		mailMan := mailman.NewMailMan(pool, incomingBox, outgoingBox)
		go mailMan.CheckingOutgoingBox()
	}

	// waiting for signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	<-c
}
