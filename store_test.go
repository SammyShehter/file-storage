package main

import (
	"bytes"
	"testing"
)

func TestStore (t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	store := NewStore(opts)

	if err := store.writeStream("test", bytes.NewReader([]byte("test data"))); err != nil {
		t.Errorf("writeStream failed: %v", err)
	}
}