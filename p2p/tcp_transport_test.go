package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddress := ":3000"
	tr := NewTCPTransport(listenAddress)

	assert.Equal(t, listenAddress, tr.listenAddress)
	assert.Nil(t, tr.ListenAndAccept())
}