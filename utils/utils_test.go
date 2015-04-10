package utils

import "testing"

type bitBoolRes struct {
	Pos uint
	Res bool
}

func TestToUint16(t *testing.T) {
	data := []struct {
		source []byte
		dest   uint16
	}{
		{[]byte{1, 2}, 258},
		{[]byte{0, 14}, 14},
		{[]byte{0, 3, 14}, 782},
	}

	for _, d := range data {
		if ToUint16(d.source) != d.dest {
			t.Fatal("result not except", ToUint16(d.source), d.dest)
		}
	}
}

func TestToUint32(t *testing.T) {
	data := []struct {
		source []byte
		dest   uint32
	}{
		{[]byte{1, 2}, 258},
		{[]byte{0, 14}, 14},
		{[]byte{0, 3, 14}, 782},
		{[]byte{6, 0, 3, 14}, 100664078},
		{[]byte{3, 6, 0, 3, 14}, 100664078},
	}

	for _, d := range data {
		if ToUint32(d.source) != d.dest {
			t.Fatal("result not except", ToUint32(d.source), d.dest)
		}
	}
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

func TestCmpString(t *testing.T) {
	data := []struct {
		a, b []string
		cmp  bool
	}{
		{[]string{"a", "b"}, []string{"a", "c"}, false},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"a", "b"}, false},
	}

	for _, d := range data {
		if CmpStrings(d.a, d.b) != d.cmp {
			t.Fatal("result not except")
		}
	}

}

func TestReadByFirstByte(t *testing.T) {
	data := []struct {
		source  []byte
		dest    []string
		isError bool
	}{
		{[]byte{4, 'a', 'b', 'c', 'd', 0}, []string{"abcd"}, false},
		{[]byte{4, 'a', 'b', 'c', 'd', 1, 'j', 0}, []string{"abcd", "j"}, false},
		{[]byte{4, 'c'}, nil, true},
		{[]byte{1, 'c'}, nil, true},
	}
	for _, d := range data {
		dest, err := ReadByFirstByte(NewRecordReader(NewBlockWithByte(d.source)))
		if d.isError && err == nil {
			t.Fatal("excepting error")
		} else if d.isError {
			continue
		} else if err != nil {
			t.Fatal("result not except error", err)
		}

		if !CmpStrings(d.dest, dest) {
			t.Fatal("result not except")
		}
	}
}
