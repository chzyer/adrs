package main

import (
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/guard"
	"github.com/chzyer/adrs/mailman"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"github.com/chzyer/adrs/wiseman"
	"gopkg.in/logex.v1"
)

func main() {
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
	mailPool := mailman.NewMailPool()

	un, err := uninet.NewUniListener(tcpListen.(*uninet.TcpURL), udpListen.(*uninet.UdpURL), pool)
	if err != nil {
		logex.Fatal(err)
	}

	g, err := guard.NewGuard(frontCorridor, backCorridor, un)
	if err != nil {
		logex.Fatal(err)
	}

	mailMan := mailman.NewMailMan(pool)

	wm, err := wiseman.NewWiseMan(frontCorridor, backCorridor, mailMan, mailPool)
	if err != nil {
		logex.Fatal(err)
	}

	g.Start()
	wm.ServeAll()
}
