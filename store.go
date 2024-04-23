package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "store"

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

func CASTransformfunc(s string) PathKey {
	hash := sha1.Sum([]byte(s))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

func DefaultPathTransformFunc(s string) PathKey {
	return PathKey{PathName: s, FileName: s}
}

type PathTransformFunc func(string) PathKey

type StoreOpts struct {
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{opts}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	return err == nil
}

func (s *Store) DeleteFile(key string) error {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	defer log.Printf("deleted %s", pathKey.FileName)
	return os.Remove(fullPathWithRoot)
}

func (s *Store) DeleteRoot(key string) error {
	pathKey := s.PathTransformFunc(key)
	rootFolder := fmt.Sprintf("%s/%s", s.Root, strings.Split(pathKey.PathName, "/")[0])
	defer log.Printf("deleted %s", rootFolder)
	return os.RemoveAll(rootFolder)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(s.Root + "/" + pathKey.FullPath())
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}

	fmt.Printf("written (%d bytes) to %s", n, fullPathWithRoot)

	return nil
}
