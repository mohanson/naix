package main

import (
	"flag"

	"github.com/mohanson/dahlia"
	"github.com/mohanson/daze"
)

var (
	fListen = flag.String("l", "", "listen address")
	fServer = flag.String("s", "", "server address")
	fCrypto = flag.String("c", "--", "crypto")
	fCipher = flag.String("k", "", "password")
)

func main() {
	flag.Parse()
	d := dahlia.NewDahlia(*fListen, *fServer, *fCipher)
	if (*fCrypto)[0] == 'x' {
		d.XI = true
	}
	if (*fCrypto)[1] == 'x' {
		d.XO = true
	}
	d.Run()
	daze.Hang()
}
