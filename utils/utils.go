package utils

import "gopkg.in/logex.v1"

func toUint16(data []byte) (ret uint16) {
	for i := range data {
		ret <<= 8
		ret += uint16(data[i])
	}
	return
}

func toUint32(data []byte) (ret uint32) {
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

func ReadBitBool(data uint64, pos uint) bool {
	data >>= pos
	if data&1 == 1 {
		return true
	}
	return false
}

func CmpString(s1, s2 []string) bool {
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
