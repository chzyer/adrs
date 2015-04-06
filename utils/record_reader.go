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
	underlying []byte
}

func NewRecordReader(data []byte) *RecordReader {
	return &RecordReader{
		underlying: data,
	}
}

func (r *RecordReader) ReadN(b []byte, n int) error {
	read, err := r.Read(b[:n])
	if err != nil {
		return logex.Trace(err)
	}
	if read != n {
		return logex.Trace(ErrShortRead)
	}
	return nil
}

func (r *RecordReader) Read(b []byte) (int, error) {
	n := copy(b, r.underlying[r.readed:])
	r.readed += n
	return n, nil
}

func (r *RecordReader) ReadBytes(n int) ([]byte, error) {
	return r.read(n)
}

func (r *RecordReader) read(i int) ([]byte, error) {
	if r.readed+i > len(r.underlying) {
		return nil, logex.Trace(io.EOF)
	}
	r.readed += i
	return r.underlying[r.readed-i : r.readed], nil
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
	copy(ret, r.underlying[r.readed:])
	return ret
}

func (r *RecordReader) RemainBytes() []byte {
	ret := make([]byte, len(r.underlying)-r.readed)
	copy(ret, r.underlying[r.readed:])
	return ret
}

func (r *RecordReader) Bytes() []byte {
	ret := make([]byte, len(r.underlying))
	copy(ret, r.underlying)
	return ret
}
