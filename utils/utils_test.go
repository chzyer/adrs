package utils

import "testing"

func TestRead8Bit(t *testing.T) {
	data := uint64(150) // 10010110

	ex := [][3]uint{
		{5, 4, 11}, // 1011
		{8, 3, 4},  // 100
	}
	for _, e := range ex {
		if Read8Bit(data, e[0], e[1]) != uint8(e[2]) {
			t.Fatal("result not except: ", e)
		}
	}
}
