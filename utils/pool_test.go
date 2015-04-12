package utils

import (
	"bytes"
	"testing"

	"gopkg.in/logex.v1"
)

var pool = NewBlockPool()

func TestBlock(t *testing.T) {
	block := pool.Get()

	bb := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	copy(block.Block, bb)
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

	n, _ := b2.Write(bb)
	if n != len(bb) || b2.Length != len(bb) {
		t.Fatal("result not except")
	}

	// plus header bytes
	if !bytes.Equal(append(Uint16To(uint16(len(bb))), bb...), b2.PlusHeaderBytes()) {
		logex.Info(b2.PlusHeaderBytes())
		t.Fatal("result not except")
	}
}
