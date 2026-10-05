package permutation

import (
	"bytes"
	"testing"
)

func TestIdentityMSBFromOne(t *testing.T) {
	in := []byte{0xA5, 0x3C}
	pblock := make([]int, 16)
	for i := range pblock {
		pblock[i] = i + 1
	}
	out, err := Permute(in, pblock, MSBToLSB, StartFromOne)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(in, out) {
		t.Fatalf("identity permutation changed bits: %x -> %x", in, out)
	}
}

func TestReverseMSBFromZero(t *testing.T) {
	in := []byte{0b10000000}
	pblock := []int{7, 6, 5, 4, 3, 2, 1, 0}
	out, err := Permute(in, pblock, MSBToLSB, StartFromZero)
	if err != nil {
		t.Fatal(err)
	}
	if out[0] != 0b00000001 {
		t.Fatalf("got %08b, want 00000001", out[0])
	}
}

func TestLSBToMSBStartZero(t *testing.T) {
	in := []byte{0x00, 0x01}
	bit, err := Bit(in, 0, LSBToMSB, StartFromZero)
	if err != nil {
		t.Fatal(err)
	}
	if bit != 1 {
		t.Fatalf("bit 0 of {00 01} in LSB-first numbering must be 1, got %d", bit)
	}
	bit, err = Bit(in, 8, LSBToMSB, StartFromZero)
	if err != nil {
		t.Fatal(err)
	}
	if bit != 0 {
		t.Fatalf("bit 8 must be 0")
	}
}

func TestOutOfRange(t *testing.T) {
	_, err := Permute([]byte{0x00}, []int{9}, MSBToLSB, StartFromOne)
	if err == nil {
		t.Fatal("expected out-of-range error")
	}
}

func TestExpandPermutation(t *testing.T) {
	in := []byte{0b10100000}
	pblock := []int{1, 1, 2, 3}
	out, err := Permute(in, pblock, MSBToLSB, StartFromOne)
	if err != nil {
		t.Fatal(err)
	}
	if out[0] != 0b11010000 {
		t.Fatalf("got %08b, want 11010000", out[0])
	}
}
