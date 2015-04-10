package utils

import (
	"bytes"
	"testing"
)

var pool = NewBlockPool()

func TestBlock(t *testing.T) {
	block := pool.Get()

	bb := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	copy(block.All, bb)
	block.Length = len(bb)

	if !bytes.Equal(bb, block.Bytes()) {
		t.Fatal("result not except")
	}

	if block.Len() != len(bb) {
		t.Fatal("result not except")
	}

	block.Recycle()
	block.Recycle()

	b2 := pool.Get()
	if b2 != block {
		t.Fatal("pool not working")
	}

}
