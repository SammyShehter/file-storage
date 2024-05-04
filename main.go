package main

import (
	"bytes"
	"log"
	"time"

	"github.com/sammyshehter/file-storage/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.NOPDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	storeOpts := StoreOpts{
		Root:              listenAddr + "_network",
		PathTransformFunc: CASTransformfunc,
	}

	fileServerOpts := FileServerOpts{
		storeOpts:      storeOpts,
		transport:      tcpTransport,
		bootstrapNodes: nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(2 * time.Second)
	go s2.Start()
	time.Sleep(2 * time.Second)

	data := bytes.NewReader([]byte("my big data file here!"))

	s2.StoreFile("key", data)

	select {}
}
