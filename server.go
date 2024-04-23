package main

import (
	"fmt"

	"github.com/sammyshehter/file-storage/p2p"
)

type FileServerOpts struct {
	storeOpts        StoreOpts
	tcpTransportOpts p2p.TCPTransportOpts
}

type FileServer struct {
	store     *Store
	transport *p2p.TCPTransport
	quitch    chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		store:     NewStore(opts.storeOpts),
		transport: p2p.NewTCPTransport(opts.tcpTransportOpts),
		quitch:    make(chan struct{}),
	}
}

func (fs *FileServer) Stop() {
	fmt.Println("Stop signal received")
	close(fs.quitch)
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

	fs.loop()

	return nil
}
