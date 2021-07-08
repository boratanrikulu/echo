package echo

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

// Listen returns a listener that supports SO_REUSEADDR and SO_REUSEPORT.
func Listen(network, address string) (net.Listener, error) {
	l := net.ListenConfig{Control: control}
	return l.Listen(context.Background(), network, address)
}

// control adds SO_REUSEADDR and SO_REUSEPORT support to the listener.
func control(network, address string, c syscall.RawConn) error {
	var err error
	c.Control(func(fd uintptr) {
		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
		if err != nil {
			return
		}

		err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
		if err != nil {
			return
		}
	})
	return err
}
