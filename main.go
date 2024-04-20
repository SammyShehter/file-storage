package main

import (
	"fmt"
	"log"

	"github.com/sammyshehter/file-storage/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("SOme calculation on connect hook")
	return nil
}

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
		OnPeer:        OnPeer,
	}

	tr := p2p.NewTCPTransport(opts)
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("received message: %s", string(msg.Payload))
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %v", err)
	}
	select {}
}
