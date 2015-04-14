package utils

import "gopkg.in/logex.v1"

var (
	ErrBlockrOverflow = logex.NewError("block overflow")
	ErrShortWritten   = logex.NewError("short written")
)

type RecordWriter struct {
	underlying *Block
}

func NewRecordWriter(b *Block) *RecordWriter {
	return &RecordWriter{
		underlying: b,
	}
}

func (w *RecordWriter) Write(b []byte) (int, error) {
	return w.underlying.Write(b)
}

func (w *RecordWriter) WriteSafe(b []byte) error {
	n, err := w.underlying.Write(b)
	if err != nil {
		return logex.Trace(err)
	} else if n != len(b) {
		return logex.Trace(ErrShortWritten)
	}
	return nil
}

func (w *RecordWriter) Block() *Block {
	return w.underlying
}

func (w *RecordWriter) WriteUint8(d uint8) error {
	if w.underlying.RemainLength() < 1 {
		return ErrBlockrOverflow
	}
	Uint8WriteTo(d, w.underlying.RemainBytes())
	w.underlying.Length += 1
	return nil
}

func (w *RecordWriter) WriteUint16(d uint16) error {
	if w.underlying.RemainLength() < 2 {
		return ErrBlockrOverflow
	}
	Uint16WriteTo(d, w.underlying.RemainBytes())
	w.underlying.Length += 2
	return nil
}

func (w *RecordWriter) WriteUint32(d uint32) error {
	if w.underlying.RemainLength() < 4 {
		return ErrBlockrOverflow
	}
	Uint32WriteTo(d, w.underlying.RemainBytes())
	w.underlying.Length += 4
	return nil
}

func (w *RecordWriter) WriteStringWithHeader(ss []string) (err error) {
	var n int
	for idx := range ss {
		if err = w.WriteUint8(uint8(len(ss[idx]))); err != nil {
			return logex.Trace(err)
		}
		if n, err = w.Write([]byte(ss[idx])); err != nil {
			return logex.Trace(err)
		} else if n != len(ss[idx]) {
			return ErrShortWritten
		}
	}
	if err = w.WriteUint8(0); err != nil {
		return logex.Trace(err)
	}
	return
}
