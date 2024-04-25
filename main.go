package main

import (
	"fmt"
	"log"

	"github.com/sammyshehter/file-storage/p2p"
)

func OnPeer(p p2p.Peer) error {
	fmt.Println("Some calculation on connect hook")
	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
		OnPeer:        OnPeer,
	}

	storeOpts := StoreOpts{
		Root:              listenAddr + "_network",
		PathTransformFunc: CASTransformfunc,
	}

	fileServerOpts := FileServerOpts{
		storeOpts:        storeOpts,
		tcpTransportOpts: tcpTransportOpts,
		bootstrapNodes:   nodes,
	}

	return NewFileServer(fileServerOpts)
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	s2.Start()
}
