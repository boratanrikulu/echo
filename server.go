package echo

import (
	"context"
	"fmt"
	"log"

	"github.com/boratanrikulu/echo/config"
)

type Server interface {
	Address(string) Server
	Banner(bool) Server
	Verbose(bool) Server
	Run(context.Context)
}

type server struct {
	address string
	banner  bool
	verbose bool
}

// NewServer returns a new server.
func NewServer() Server {
	s := server{}
	return &s
}

// Address sets listening address for the server.
func (s *server) Address(a string) Server {
	s.address = a
	return s
}

// Verbose sets verbose feature.
func (s *server) Verbose(v bool) Server {
	s.verbose = v
	return s
}

// Banner sets banner on/off.
func (s *server) Banner(b bool) Server {
	s.banner = b
	return s
}

// Run starts server and listening the clients.
// Hangs until ctx is canceled.
func (s *server) Run(ctx context.Context) {
	if s.banner {
		fmt.Printf("%s\n\n", config.Banner)
	}

	et, err := NewEchoTCP(s.address, s.verbose)
	if err != nil {
		log.Fatal(err) // exit if creating EchoTCP is failed.
	}
	defer et.listener.Close()

	fmt.Printf("server is started at %s\n", s.address)
	et.Run(ctx)
}
