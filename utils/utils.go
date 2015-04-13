package utils

import (
	"time"

	"gopkg.in/logex.v1"
)

var cacheTime time.Time

func init() {
	cacheTime = time.Now()
	go func() {
		for cacheTime = range time.Tick(time.Second) {
		}
	}()
}

func Now() time.Time {
	return cacheTime
}

func Uint16To(b uint16) []byte {
	mask := uint16(1<<8 - 1)
	return []byte{uint8(b >> 8), uint8(b & mask)}
}

func ToUint16(data []byte) (ret uint16) {
	for i := range data {
		ret <<= 8
		ret += uint16(data[i])
	}
	return
}
func Uint8WriteTo(d uint8, b []byte) {
	b[0] = d
}

func Uint16WriteTo(d uint16, b []byte) {
	mask := uint16(1<<8 - 1)
	for i := 1; i >= 0; i-- {
		b[i] = uint8(d & mask)
		d = d >> 8
	}
}

func Uint32WriteTo(d uint32, b []byte) {
	mask := uint32(1<<8 - 1)
	for i := 3; i >= 0; i-- {
		b[i] = uint8(d & mask)
		d = d >> 8
	}
}

func ToUint32(data []byte) (ret uint32) {
	for i := range data {
		ret <<= 8
		ret += uint32(data[i])
	}
	return
}

func Read8Bit(data uint64, start, length uint) uint8 {
	// generate length*1
	mask := uint64(1<<length - 1)
	data >>= start
	return uint8(data & mask)
}

func ReadBit(data uint64, pos uint) uint8 {
	data >>= pos
	if data&1 == 1 {
		return 1
	}
	return 0
}

func ReadBitBool(data uint64, pos uint) bool {
	data >>= pos
	if data&1 == 1 {
		return true
	}
	return false
}

func CmpStrings(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func ReadByFirstByte(r *RecordReader) ([]string, error) {
	var (
		ret     []string
		length  uint8
		err     error
		segment []byte = make([]byte, 1<<8)
	)
	for {
		length, err = r.ReadByte()
		if err != nil {
			return nil, logex.Trace(err)
		}

		if length == 0 {
			break
		}

		err := r.ReadN(segment, int(length))
		if err != nil {
			return nil, logex.Trace(err)
		}

		ret = append(ret, string(segment[:length]))
	}
	return ret, nil
}
