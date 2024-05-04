package p2p

import (
	"testing"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       GOBDecoder{},
	}
	tr := NewTCPTransport(opts)

	if err := tr.ListenAndAccept(); err != nil {
		t.Fatalf("failed to listen and accept: %v", err)
	}
}
