package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "somebestpicture"
	pathKey := CASTransformfunc(key)
	originalPathName := "7432e4aced777c0dedbead6bcd2ab94ec33682a1"
	expectedPathName := "7432e/4aced/777c0/dedbe/ad6bc/d2ab9/4ec33/682a1"
	if pathKey.PathName != expectedPathName || pathKey.FileName != originalPathName {
		t.Errorf("CASTransformfunc failed. Expected %s, got %s", expectedPathName, pathKey.PathName)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASTransformfunc,
	}
	store := NewStore(opts)
	testKey := "testKey"
	testData := []byte("test data")
	if err := store.writeStream(testKey, bytes.NewReader([]byte(testData))); err != nil {
		t.Errorf("writeStream failed: %v", err)
	}

	dataRead, err := store.Read(testKey)
	if err != nil {
		t.Errorf("Read failed: %v", err)
	}

	b, _ := io.ReadAll(dataRead)

	if string(testData) != string(b) {
		t.Errorf("Read failed. Expected %s, got %s", "test data", string(b))
	}

	h := store.Has(testKey)
	if !h {
		t.Errorf("Has failed. Expected %v, got %v", true, h)
	}

	err = store.DeleteRoot(testKey)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}
