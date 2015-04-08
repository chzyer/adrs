package main

import (
	"flag"

	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/guard"
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
	flag.StringVar(&f.Listen, "-l", ":53", "listen udp")
	flag.Parse()
	return f
}

func main() {
	args := NewFlag()
	wayIn := make(customer.Corridor)
	wayOut := make(customer.Corridor)
	host := &uninet.Host{
		UDP: args.Listen,
	}

	pool := utils.NewBlockPool()
	un, err := uninet.NewUniNet(host, pool)
	if err != nil {
		logex.Fatal(err)
	}

	g, err := guard.NewGuard(wayIn, wayOut, un, pool)
	if err != nil {
		logex.Fatal(err)
	}
	g.Start()

	wm, err := wiseman.NewWiseMan(wayIn, wayOut)
	if err != nil {
		logex.Fatal(err)
	}
	wm.Serve()
}
