package main

import (
	"flag"

	"github.com/mohanson/daze"
	"github.com/mohanson/naix"
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
		naix.NewServer(*fListen, *fServer, *fCipher).Run()
		daze.Hang()
	case "client":
		naix.NewClient(*fListen, *fServer, *fCipher).Run()
		daze.Hang()
	case "middle":
		naix.NewMiddle(*fListen, *fServer).Run()
		daze.Hang()
	}
}
