package utils

import "testing"

type bitBoolRes struct {
	Pos uint
	Res bool
}

func TestRead8Bit(t *testing.T) {
	data := uint64(150) // 10010110

	t8b := [][3]uint{
		{1, 4, 11}, // 1011
		{5, 3, 4},  // 100
	}
	for _, e := range t8b {
		if Read8Bit(data, e[0], e[1]) != uint8(e[2]) {
			t.Fatal("result not except: ", e)
		}
	}

	tbb := []bitBoolRes{
		{1, true},
		{2, true},
		{3, false},
		{0, false},
		{7, true},
	}
	for _, e := range tbb {
		if ReadBitBool(data, e.Pos) != e.Res {
			t.Fatal("result not except", e)
		}
	}
}
