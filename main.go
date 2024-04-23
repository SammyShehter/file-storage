package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sammyshehter/file-storage/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("Some calculation on connect hook")
	return nil
}

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
		OnPeer:        OnPeer,
	}

	storeOpts := StoreOpts{
		Root:              "store_prod",
		PathTransformFunc: CASTransformfunc,
	}

	fileServerOpts := FileServerOpts{
		storeOpts:        storeOpts,
		tcpTransportOpts: tcpTransportOpts,
	}
	s := NewFileServer(fileServerOpts)

	go func () {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()
	
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
