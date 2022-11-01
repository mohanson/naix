package dahlia

import (
	"errors"
	"io"
	"log"
	"math"
	"net"

	"github.com/mohanson/daze"
	"github.com/mohanson/daze/protocol/ashe"
)

// Dahlia is an encrypted port forwarding protocol. Unlike common port forwarding tools, it needs to configure a server
// and a client, and the communication between the server and the client is encrypted to bypass firewall detection.

// Dahlia implemented the dahlia protocol.
type Dahlia struct {
	Cipher []byte
	Closer io.Closer
	Listen string
	Server string
	XI     bool
	XO     bool
}

// Close listener. Established connections will not be closed.
func (d *Dahlia) Close() error {
	if d.Closer != nil {
		return d.Closer.Close()
	}
	return nil
}

// Serve incoming connections. Parameter cli will be closed automatically when the function exits.
func (d *Dahlia) Serve(ctx *daze.Context, cli net.Conn) error {
	var (
		a   io.ReadWriteCloser
		b   io.ReadWriteCloser
		err error
		srv net.Conn
	)
	if d.XI {
		s := ashe.Server{Cipher: d.Cipher}
		a, err = s.ServeCipher(ctx, cli)
		if err != nil {
			return err
		}
	} else {
		a = cli
	}
	srv, err = daze.Dial("tcp", d.Server)
	if err != nil {
		return err
	}
	if d.XO {
		c := ashe.Client{Cipher: d.Cipher}
		b, err = c.WithCipher(ctx, srv)
		if err != nil {
			return err
		}
	} else {
		b = srv
	}
	daze.Link(a, b)
	return nil
}

// Run it.
func (d *Dahlia) Run() error {
	l, err := net.Listen("tcp", d.Listen)
	if err != nil {
		return err
	}
	d.Closer = l
	log.Println("main: listen and serve on", d.Listen)

	go func() {
		idx := uint32(math.MaxUint32)
		for {
			cli, err := l.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Println("main:", err)
				}
				break
			}
			idx++
			ctx := &daze.Context{Cid: idx}
			log.Printf("conn: %08x accept remote=%s", ctx.Cid, cli.RemoteAddr())
			go func(cli net.Conn) {
				defer cli.Close()
				if err := d.Serve(ctx, cli); err != nil {
					log.Printf("conn: %08x  error %s", ctx.Cid, err)
				}
				log.Printf("conn: %08x closed", ctx.Cid)
			}(cli)
		}
	}()
	return nil
}

// NewDahlia returns a new dahlia.
func NewDahlia(listen string, server string, cipher string) *Dahlia {
	return &Dahlia{
		Cipher: daze.Salt(cipher),
		Listen: listen,
		Server: server,
	}
}
