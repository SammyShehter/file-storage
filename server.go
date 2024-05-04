package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/sammyshehter/file-storage/p2p"
)

type Message struct {
	From    string
	Payload any
}

type DataMessage struct {
	Key  string
	Data []byte
}

type FileServerOpts struct {
	storeOpts      StoreOpts
	transport      p2p.Transport
	bootstrapNodes []string
}

type FileServer struct {
	store          *Store
	transport      p2p.Transport
	quitch         chan struct{}
	bootstrapNodes []string
	peers          map[string]p2p.Peer
	peerLock       sync.Mutex
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		store:          NewStore(opts.storeOpts),
		transport:      opts.transport,
		quitch:         make(chan struct{}),
		bootstrapNodes: opts.bootstrapNodes,
		peers:          make(map[string]p2p.Peer),
	}
}

func (fs *FileServer) Stop() {
	fmt.Println("Stop signal received")
	close(fs.quitch)
}

func (fs *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

func (fs *FileServer) StoreFile(key string, r io.Reader) error {
	buf := new(bytes.Buffer)

	tee := io.TeeReader(r, buf)

	if err := fs.store.Write(key, tee); err != nil {
		return err
	}

	p := &DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}

	return fs.broadcast(&Message{
		From:    "fs.transport.ListedAddr",
		Payload: p,
	})
}

func (fs *FileServer) OnPeer(p p2p.Peer) error {
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()

	fs.peers[p.RemoteAddr().String()] = p
	log.Printf("New connection peer: %s", p.RemoteAddr())
	return nil
}

func (fs *FileServer) loop() {
	defer func() {
		fs.transport.Close()
	}()

	for {
		select {
		case msg := <-fs.transport.Consume():
			var m Message
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&m); err != nil {
				log.Println(err)
			}
			if err := fs.handleMessage(&m); err != nil {
				log.Println(err)
			}
		case <-fs.quitch:
			fmt.Println("Stop the loop")
			return
		}
	}
}

func (s *FileServer) handleMessage(msg *Message) error {
	switch v := msg.Payload.(type) {
	case *DataMessage:
		fmt.Println("received data %+v\n", v)	
	}

	return nil
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
