package main

import (
	"flag"

	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/guard"
	"github.com/chzyer/adrs/mailman"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"github.com/chzyer/adrs/wiseman"
	"gopkg.in/logex.v1"
)

type Flag struct {
	Listen string
}

func NewFlag() *Flag {
	f := new(Flag)
	flag.StringVar(&f.Listen, "-l", "udp://:53", "listen udp")
	flag.Parse()
	return f
}

func main() {
	args := NewFlag()
	frontCorridor := make(customer.Corridor)
	backCorridor := make(customer.Corridor)
	listenURL, err := uninet.ParseURL(args.Listen)
	if err != nil {
		logex.Fatal(err)
	}
	host := &uninet.Host{
		UDP: listenURL,
	}
	pool := utils.NewBlockPool()
	mailPool := mailman.NewMailPool()

	un, err := uninet.NewUniNet(host)
	if err != nil {
		logex.Fatal(err)
	}

	g, err := guard.NewGuard(frontCorridor, backCorridor, un, pool)
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
