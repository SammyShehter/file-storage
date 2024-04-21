package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func CASTransformfunc(s string) string {
	hash := sha1.Sum([]byte(s))
	hashStr := hex.EncodeToString(hash[:])


}

type PathTransformFunc func(string) string

var DefaultPathTransformFunc = func(s string) string {
	return s
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{opts}
}

func (s *Store) writeStream (key string, t io.Reader) error {
	pathname := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
		return err
	}

	pathAndFileName := pathname + "/data"

	f, err := os.Create(pathAndFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, t)
	if err != nil {
		return err
	}

	fmt.Printf("written (%d bytes) to %s", n, pathAndFileName)

	return nil
} 
