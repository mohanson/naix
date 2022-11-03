package naix

import (
	"errors"
	"io"
	"log"
	"math"
	"net"

	"github.com/mohanson/daze"
	"github.com/mohanson/daze/protocol/ashe"
)

// Naix is an encrypted port forwarding protocol. Unlike common port forwarding tools, it needs to configure a server
// and a client, and the communication between the server and the client is encrypted to bypass firewall detection.

// Server implemented the naix protocol.
type Server struct {
	Cipher []byte
	Closer io.Closer
	Listen string
	Server string
}

// Close listener. Established connections will not be closed.
func (s *Server) Close() error {
	if s.Closer != nil {
		return s.Closer.Close()
	}
	return nil
}

// Serve incoming connections. Parameter cli will be closed automatically when the function exits.
func (s *Server) Serve(ctx *daze.Context, cli net.Conn) error {
	dec := ashe.Server{Cipher: s.Cipher}
	spy, err := dec.ServeCipher(ctx, cli)
	if err != nil {
		return err
	}
	srv, err := daze.Dial("tcp", s.Server)
	if err != nil {
		return err
	}
	daze.Link(spy, srv)
	return nil
}

// Run it.
func (s *Server) Run() error {
	l, err := net.Listen("tcp", s.Listen)
	if err != nil {
		return err
	}
	s.Closer = l
	log.Println("main: listen and serve on", s.Listen)

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
				if err := s.Serve(ctx, cli); err != nil {
					log.Printf("conn: %08x  error %s", ctx.Cid, err)
				}
				log.Printf("conn: %08x closed", ctx.Cid)
			}(cli)
		}
	}()
	return nil
}

// NewServer returns a new Server.
func NewServer(listen string, server string, cipher string) *Server {
	return &Server{
		Cipher: daze.Salt(cipher),
		Listen: listen,
		Server: server,
	}
}

// Client implemented the naix protocol.
type Client struct {
	Cipher []byte
	Closer io.Closer
	Listen string
	Server string
}

// Close listener. Established connections will not be closed.
func (c *Client) Close() error {
	if c.Closer != nil {
		return c.Closer.Close()
	}
	return nil
}

// Serve incoming connections. Parameter cli will be closed automatically when the function exits.
func (c *Client) Serve(ctx *daze.Context, cli net.Conn) error {
	srv, err := daze.Dial("tcp", c.Server)
	if err != nil {
		return err
	}
	enc := ashe.Client{Cipher: c.Cipher}
	spy, err := enc.WithCipher(ctx, srv)
	if err != nil {
		srv.Close()
		return err
	}
	daze.Link(cli, spy)
	return nil
}

// Run it.
func (c *Client) Run() error {
	l, err := net.Listen("tcp", c.Listen)
	if err != nil {
		return err
	}
	c.Closer = l
	log.Println("main: listen and serve on", c.Listen)

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
				if err := c.Serve(ctx, cli); err != nil {
					log.Printf("conn: %08x  error %s", ctx.Cid, err)
				}
				log.Printf("conn: %08x closed", ctx.Cid)
			}(cli)
		}
	}()
	return nil
}

// NewClient returns a new Client.
func NewClient(listen string, server string, cipher string) *Client {
	return &Client{
		Cipher: daze.Salt(cipher),
		Listen: listen,
		Server: server,
	}
}

// Middle implemented the naix protocol.
type Middle struct {
	Closer io.Closer
	Listen string
	Server string
}

// Close listener. Established connections will not be closed.
func (m *Middle) Close() error {
	if m.Closer != nil {
		return m.Closer.Close()
	}
	return nil
}

// Serve incoming connections. Parameter cli will be closed automatically when the function exits.
func (m *Middle) Serve(ctx *daze.Context, cli net.Conn) error {
	srv, err := daze.Dial("tcp", m.Server)
	if err != nil {
		return err
	}
	daze.Link(cli, srv)
	return nil
}

// Run it.
func (m *Middle) Run() error {
	l, err := net.Listen("tcp", m.Listen)
	if err != nil {
		return err
	}
	m.Closer = l
	log.Println("main: listen and serve on", m.Listen)

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
				if err := m.Serve(ctx, cli); err != nil {
					log.Printf("conn: %08x  error %s", ctx.Cid, err)
				}
				log.Printf("conn: %08x closed", ctx.Cid)
			}(cli)
		}
	}()
	return nil
}

// NewMiddle returns a new Middle.
func NewMiddle(listen string, server string) *Middle {
	return &Middle{
		Listen: listen,
		Server: server,
	}
}
