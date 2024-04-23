package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASTransformfunc,
	}

	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Errorf("Failed to clear root directory. Got %s", err)
	}
}

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
	store := newStore()
	defer teardown(t, store)

	for i := 0; i < 50; i++ {
		testKey := fmt.Sprintf("test_key_%d", i)
		testData := []byte(fmt.Sprintf("test_data_%d", i))

		if err := store.Write(testKey, bytes.NewReader([]byte(testData))); err != nil {
			t.Errorf("Store.Write failed: %v", err)
		}

		if ok := store.Has(testKey); !ok {
			t.Errorf("store.Has failed. Expected %v, got %v", true, ok)
		}

		dataRead, err := store.Read(testKey)
		if err != nil {
			t.Errorf("Read failed: %v", err)
		}

		b, _ := io.ReadAll(dataRead)
		if string(testData) != string(b) {
			t.Errorf("Read failed. Expected %s, got %s", "test data", string(b))
		}

		if err = store.DeleteRoot(testKey); err != nil {
			t.Errorf("DeleteFile failed: %v", err)
		}

		if ok := store.Has(testKey); ok {
			t.Errorf("store.Has failed. Expected %v, got %v", false, ok)
		}
	}
}
