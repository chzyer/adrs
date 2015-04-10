package utils

import (
	"errors"
	"io"

	"gopkg.in/logex.v1"
)

var (
	ErrShortRead = errors.New("short read")
)

type RecordReader struct {
	readed     int
	underlying *Block
}

func NewRecordReader(b *Block) *RecordReader {
	return &RecordReader{
		underlying: b,
	}
}

func (r *RecordReader) Block() *Block {
	return r.underlying
}

func (r *RecordReader) ReadN(b []byte, n int) error {
	if len(b) < n {
		return logex.NewError("n large than bytes")
	}
	read, _ := r.Read(b[:n])
	if read != n {
		return logex.Trace(ErrShortRead)
	}
	return nil
}

func (r *RecordReader) Read(b []byte) (int, error) {
	n := copy(b, r.underlying.Bytes()[r.readed:])
	r.readed += n
	return n, nil
}

func (r *RecordReader) ReadBytes(n int) ([]byte, error) {
	return r.read(n)
}

func (r *RecordReader) read(i int) ([]byte, error) {
	if r.readed+i > r.underlying.Len() {
		return nil, logex.Trace(io.EOF)
	}
	r.readed += i
	return r.underlying.Bytes()[r.readed-i : r.readed], nil
}

func (r *RecordReader) ReadUint8() (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, logex.Trace(err)
	}
	return uint8(b), nil
}

func (r *RecordReader) ReadUint32() (uint32, error) {
	d, err := r.read(4)
	if err != nil {
		return 0, logex.Trace(err)
	}
	return ToUint32(d), nil
}

func (r *RecordReader) ReadUint16() (uint16, error) {
	d, err := r.read(2)
	if err != nil {
		return 0, logex.Trace(err)
	}
	return ToUint16(d), nil
}

func (r *RecordReader) ReadByte() (byte, error) {
	d, err := r.read(1)
	if err != nil {
		return 0, logex.Trace(err)
	}
	return d[0], nil
}

func (r *RecordReader) Peek(n int) []byte {
	ret := make([]byte, n)
	copy(ret, r.underlying.Bytes()[r.readed:])
	return ret
}

func (r *RecordReader) RemainBytes() []byte {
	ret := make([]byte, r.underlying.Len()-r.readed)
	copy(ret, r.underlying.Bytes()[r.readed:])
	return ret
}

func (r *RecordReader) Bytes() []byte {
	ret := make([]byte, r.underlying.Len())
	copy(ret, r.underlying.Bytes())
	return ret
}
