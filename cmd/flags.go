package main

import (
	"flag"
	"strconv"
)

type Flags struct {
	fs *flag.FlagSet

	portNumber int
}

func NewFlagSet() *Flags {
	fgs := &Flags{
		fs: flag.NewFlagSet("gopher-translator", flag.ContinueOnError),
	}

	fgs.fs.IntVar(&fgs.portNumber, "port", 8080, "Server Port Number")

	return fgs
}

func (g *Flags) PortNumber() string {
	return strconv.Itoa(g.portNumber)
}

func (g *Flags) Init(args []string) error {
	return g.fs.Parse(args)
}
