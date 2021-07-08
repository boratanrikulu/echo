package echo

import (
	"net"
	"testing"
)

func BenchmarkEchoTCP_handleConnection(b *testing.B) {
	m := []byte("test test test")
	server, client := net.Pipe()
	go handleConnection(server)

	for n := 0; n < b.N; n++ {
		_, err := client.Write(m)
		if err != nil {
			b.Errorf("error while writing to server: %s", err)
		}

		var got []byte
		buf := make([]byte, 1024)
		for {
			size, err := client.Read(buf[:])
			if err != nil {
				b.Errorf("error while reading from server: %s", err)
			}

			got = append(got, buf[:size]...)
			if len(got) >= len(m) {
				break
			}
		}
	}
}
