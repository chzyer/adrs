package utils

import (
	"bytes"
	"io"
	"testing"

	"gopkg.in/logex.v1"
)

func TestRecordReader(t *testing.T) {
	bb := []byte{1, 2, 3, 4, 5, 6}
	block := NewBlockWithByte(bb)
	r := NewRecordReader(block)

	if r.Block() != block {
		t.Fatal("result not except")
	}

	{ //ReadN
		tmp := make([]byte, 1)
		err := r.ReadN(tmp, 5)
		if err == nil {
			t.Fatal("excepting error")
		}

		tmp = make([]byte, 2)
		err = r.ReadN(tmp, 1)
		if err != nil {
			t.Fatal("result not except")
		}
		if !bytes.Equal(tmp, []byte{1, 0}) {
			t.Fatal("result not except")
		}

		tmp = make([]byte, 10)
		err = r.ReadN(tmp, 10)
		if !logex.Is(err, ErrShortRead) {
			t.Fatal("result not except", err)
		}
		r = NewRecordReader(block)
	}

	{ // ReadByte
		n, err := r.ReadByte()
		if err != nil {
			t.Fatal("result not except")
		}
		if n != bb[0] {
			t.Fatal("reuslt not except")
		}

		read, err := r.ReadBytes(5)
		if err != nil {
			t.Fatal("result not except")
		}
		if !bytes.Equal(read, []byte{2, 3, 4, 5, 6}) {
			t.Fatal("result not except")
		}
		_, err = r.ReadBytes(1)
		if !logex.Is(err, io.EOF) {
			t.Fatal("excepting error")
		}
		r = NewRecordReader(block)
	}

	{ //readuint
		u, err := r.ReadUint8()
		if err != nil || u != 1 {
			t.Fatal("result not except")
		}

		u16, err := r.ReadUint16()
		if err != nil || u16 != 515 {
			t.Fatal("result not except")
		}

		_, err = r.ReadUint32()
		if err == nil {
			t.Fatal("excepting error")
		}

		r = NewRecordReader(block)
		u32, err := r.ReadUint32()
		if err != nil || u32 != 16909060 {
			t.Fatal("result not except")
		}
	}

	{ // Peek
		if !bytes.Equal(r.Peek(1), []byte{5}) {
			t.Fatal("result not except")
		}

		if !bytes.Equal(r.RemainBytes(), []byte{5, 6}) {
			t.Fatal("result not except")
		}

		if !bytes.Equal(r.Bytes(), bb) {
			t.Fatal("result not except")
		}
	}
}
