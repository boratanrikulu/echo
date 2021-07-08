package echo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
)

// TestNewEchoTCP tests creating EchoTCP.
func TestNewEchoTCP(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		want      *EchoTCP
		wantErr   bool
		testEqual bool
	}{
		{
			name:      "should-not-work-with-empty-address",
			address:   "",
			want:      nil,
			wantErr:   true,
			testEqual: false,
		},
		{
			name:      "should-not-work-with-wrong-address",
			address:   "wrong",
			want:      nil,
			wantErr:   true,
			testEqual: false,
		},
		{
			name:      "should-work-without-defining",
			address:   ":",
			want:      nil,
			wantErr:   false,
			testEqual: false, // Port will be random. Testing equalution is not needed.
		},
		{
			name:      "should-work-with-defining",
			address:   "127.0.0.1:1337",
			want:      nil,
			wantErr:   false,
			testEqual: false, // Port will be random. Testing equalution is not needed.
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEchoTCP(tt.address, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEchoTCP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.testEqual && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEchoTCP() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEchoTCP_handleConnection test handleConnection() by using net.Pipe.
// Sends short and long messages for testing.
func TestEchoTCP_handleConnection(t *testing.T) {
	tests := []struct {
		name     string
		messages [][]byte
	}{
		{
			name: "should-be-equal-all-messages",
			messages: [][]byte{
				[]byte("test"), []byte("test1"), []byte("test2"), []byte("test3"),
			},
		},
		{
			name: "should-work-long-messages",
			messages: [][]byte{
				[]byte("Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat sapiente soluta similique, non provident quam ut eos odio hic beatae dicta amet reiciendis, culpa quibusdam ipsa tenetur, accusamus nulla omnis?"),
				[]byte("Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat sapiente soluta similique, non provident quam ut eos odio hic beatae dicta amet reiciendis, culpa quibusdam ipsa tenetur, accusamus nulla omnis? Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat sapiente soluta similique, non provident quam ut eos odio hic beatae dicta amet reiciendis, culpa quibusdam ipsa tenetur, accusamus nulla omnis?"),
				[]byte("Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat sapiente soluta similique, non provident quam ut eos odio hic beatae dicta amet reiciendis, culpa quibusdam ipsa tenetur, accusamus nulla omnis? Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat sapiente soluta similique, non provident quam ut eos odio hic beatae dicta amet reiciendis, culpa quibusdam ipsa tenetur, accusamus nulla omnis? Lorem ipsum dolor sit amet consectetur adipisicing elit. Repellat sapiente soluta similique, non provident quam ut eos odio hic beatae dicta amet reiciendis, culpa quibusdam ipsa tenetur, accusamus nulla omnis?"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := net.Pipe()
			go handleConnection(server)

			for _, m := range tt.messages {
				_, err := io.Copy(client, bytes.NewReader(m))
				if err != nil {
					t.Errorf("error while writing to server, error: %s", err)
				}

				var got []byte
				buf := make([]byte, 1024)
				for {
					size, err := client.Read(buf[:])
					if err != nil {
						t.Errorf("error while reading from server, error: %s", err)
					}

					got = append(got, buf[:size]...)
					if len(got) >= len(m) {
						break
					}
				}

				if !bytes.Equal(got, m) {
					t.Errorf("got = \"%s\", want = \"%s\"\n", got, m)
				}
			}
		})
	}
}

// TestEchoTCP_handleConnectionWithRealServer test handleConection by using real EchoTCP server.
// Sends 1000, 10000, 1000000 messages for testing.
func TestEchoTCP_handleConnectionWithRealServer(t *testing.T) {
	et, err := NewEchoTCP(":1338", false)
	if err != nil {
		t.Fatalf("error while creating echo server, error: %s", err)
	}

	go et.Run(context.Background())

	m := []byte("test test test")
	conn, err := net.Dial("tcp4", ":1338")
	if err != nil {
		fmt.Printf("error while dialing, error: %s", err)
	}

	tests := []int{1000, 10000, 1000000}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-messages", tt), func(t *testing.T) {
			for n := 0; n < tt; n++ {
				_, err := conn.Write(m)
				if err != nil {
					t.Errorf("error while writing to server: %s", err)
				}

				var got []byte
				buf := make([]byte, 1024)
				for {
					size, err := conn.Read(buf[:])
					if err != nil {
						t.Errorf("error while reading from server: %s", err)
					}

					got = append(got, buf[:size]...)
					if len(got) >= len(m) {
						break
					}
				}
			}
		})
	}
}
