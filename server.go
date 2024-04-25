package main

import (
	"fmt"
	"log"

	"github.com/sammyshehter/file-storage/p2p"
)

type FileServerOpts struct {
	storeOpts        StoreOpts
	tcpTransportOpts p2p.TCPTransportOpts
	bootstrapNodes   []string
}

type FileServer struct {
	store          *Store
	transport      *p2p.TCPTransport
	quitch         chan struct{}
	bootstrapNodes []string
	peers          map[string]p2p.Peer
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		store:          NewStore(opts.storeOpts),
		transport:      p2p.NewTCPTransport(opts.tcpTransportOpts),
		quitch:         make(chan struct{}),
		bootstrapNodes: opts.bootstrapNodes,
		peers:          make(map[string]p2p.Peer),
	}
}

func (fs *FileServer) Stop() {
	fmt.Println("Stop signal received")
	close(fs.quitch)
}

func (fs *FileServer) OnPeer(p p2p.Peer) {
	
}

func (fs *FileServer) loop() {
	defer func() {
		fs.transport.Close()
	}()

	for {
		select {
		case msg := <-fs.transport.Consume():
			fmt.Printf("Received message: %s", string(msg.Payload))
		case <-fs.quitch:
			fmt.Println("Stop the loop")
			return
		}
	}
}

func (fs *FileServer) Start() error {
	if err := fs.transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(fs.bootstrapNodes) != 0 {
		fs.bootstrapNetwork()
	}

	fs.loop()

	return nil
}

func (fs *FileServer) bootstrapNetwork() error {
	for _, addr := range fs.bootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := fs.transport.Dial(addr); err != nil {
				log.Println("dial err: ", err)
			}
		}(addr)
	}
	return nil
}
