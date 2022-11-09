package core

import "testing"

func Test_SegmentByJsonByteLength_SmallSlice_ReturnsAll(t *testing.T) {
	values := []string{"1", "2", "3"}
	input := ToInterfaceSlice(values)

	result, err := SegmentByJsonByteLength(input, 262144, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatal("Slice was not chunked properly")
	}

	if len(result[0]) != 3 {
		t.Fatal("All slice parts were not retrieved")
	}
}

func Test_SegmentByJsonByteLength_MediumSlice_ReturnsAll(t *testing.T) {
	values := []string{"1", "2", "3"}
	input := ToInterfaceSlice(values)

	result, err := SegmentByJsonByteLength(input, 262144, 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatal("Slice was not chunked properly")
	}

	if len(result[0]) != 2 && len(result[1]) != 1 {
		t.Fatal("All slice parts were not retrieved")
	}
}
