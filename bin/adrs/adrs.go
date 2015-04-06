package main

import "flag"

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
	_ = args
}
