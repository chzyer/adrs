package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/chzyer/adrs/conf"
	"github.com/chzyer/adrs/customer"
	"github.com/chzyer/adrs/guard"
	"github.com/chzyer/adrs/mailman"
	"github.com/chzyer/adrs/uninet"
	"github.com/chzyer/adrs/utils"
	"github.com/chzyer/adrs/wiseman"
	"gopkg.in/logex.v1"
)

func GetConfig() (*conf.Config, error) {
	confPath := ""
	flag.StringVar(&confPath, "c", "adrs.conf", "conf path")
	flag.Parse()

	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, logex.Trace(err)
	}
	config, err := conf.NewConfig(bytes.NewReader(data))
	if err != nil {
		return nil, logex.Trace(err)
	}
	return config, nil
}

func main() {
	config, err := GetConfig()
	if err != nil {
		logex.Fatal(err)
	}

	var (
		mailManNum = 1
		wiseManNum = 1
	)

	frontCorridor := make(customer.Corridor)
	backCorridor := make(customer.Corridor)

	pool := utils.NewBlockPool()
	incomingBox, outgoingBox := mailman.MakeBoxes()

	un, err := uninet.NewUniListener(config.Listen, pool)
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
		mailMan := mailman.NewMailMan(config.Router, pool, incomingBox, outgoingBox)
		go mailMan.CheckingOutgoingBox()
	}

	// waiting for signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	<-c
}
