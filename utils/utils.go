package utils

func toUint16(data []byte) (ret uint16) {
	for i := range data {
		ret <<= 8
		ret += uint16(data[i]) // TODO: bit protection
	}
	return
}

func Read8Bit(data uint64, start, length uint) uint8 {
	// generate length*1
	mask := uint64(1<<length - 1)
	data >>= start - length
	return uint8(data & mask)
}
