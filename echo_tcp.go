package echo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
)

var ErrListenerFail = errors.New("starting listener is failed")
var ErrNoAddress = errors.New("server address is not defined")

// EchoTCP keeps net.Listener and verbose status.
type EchoTCP struct {
	listener net.Listener
	verbose  bool
}

// NewEchoTCP returns a new EchoTCP.
func NewEchoTCP(address string, verbose bool) (*EchoTCP, error) {
	if address == "" {
		return nil, ErrNoAddress
	}

	l, err := Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("%w, error: %s", ErrListenerFail, err)
	}

	return &EchoTCP{
		listener: l,
		verbose:  verbose,
	}, nil
}

// Run starts accepting connections.
func (et *EchoTCP) Run(ctx context.Context) {
	et.acceptConnections(ctx)
}

// acceptConnections accepts connects an handle them as goroutines.
// Runs until context is canceled.
func (et *EchoTCP) acceptConnections(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			fmt.Println("context is canceled")
			return
		}

		conn, err := et.listener.Accept()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

// handleConnection reads from the client.
// Sends back to the originating source any data it receives. (rfc 862)
// Runs until client closes the connection.
func handleConnection(conn net.Conn) {
	defer conn.Close()
	cAddr := conn.RemoteAddr().String()
	fmt.Printf("client %s is connected\n", cAddr)

	buf := make([]byte, 1024)
	for {
		size, err := conn.Read(buf[:])

		if err != nil {
			if err != io.EOF {
				fmt.Printf("error: %s\n", err)
			}
			fmt.Printf("client %s is disconnected\n", cAddr)
			return
		}

		_, err = conn.Write(buf[:size])
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
	}
}
