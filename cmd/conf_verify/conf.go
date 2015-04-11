package main

import (
	"bytes"
	"flag"
	"io/ioutil"

	"github.com/chzyer/adrs/conf"
	"gopkg.in/logex.v1"
)

func main() {
	confPath := ""
	flag.StringVar(&confPath, "c", "adrs.conf", "conf path")
	flag.Parse()

	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		logex.Fatal(err)
	}
	c := new(conf.RawConfig)
	err = conf.Parse(c, bytes.NewReader(data))
	if err != nil {
		logex.Fatal(err)
	}
	logex.Pretty(c)
}
