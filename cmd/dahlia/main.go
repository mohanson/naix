package main

import (
	"flag"

	"github.com/mohanson/dahlia"
	"github.com/mohanson/daze"
)

var (
	fListen = flag.String("l", "", "listen address")
	fServer = flag.String("s", "", "server address")
	fCipher = flag.String("k", "", "password")
)

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "server":
		dahlia.NewServer(*fListen, *fServer, *fCipher).Run()
	case "client":
		dahlia.NewClient(*fListen, *fServer, *fCipher).Run()
	}
	daze.Hang()
}
