package utils

import (
	"io"

	"gopkg.in/logex.v1"
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

func (r *RecordReader) read(i int) ([]byte, error) {
	if r.readed+i > len(r.underlying) {
		return nil, logex.Trace(io.EOF)
	}
	r.readed += i
	return r.underlying[r.readed-i : r.readed], nil
}

func (r *RecordReader) ReadUint16() (uint16, error) {
	d, err := r.read(2)
	if err != nil {
		return 0, logex.Trace(err)
	}
	return toUint16(d), nil
}

func (r *RecordReader) ReadByte() (byte, error) {
	d, err := r.read(1)
	if err != nil {
		return 0, logex.Trace(err)
	}
	return d[0], nil
}
