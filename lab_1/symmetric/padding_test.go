package symmetric

import (
	"bytes"
	"testing"
)

func TestPKCS7AlwaysAddsBlock(t *testing.T) {
	block := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	padded, err := pad(block, 8, PKCS7)
	if err != nil {
		t.Fatal(err)
	}
	if len(padded) != 16 {
		t.Fatalf("len=%d", len(padded))
	}
	un, err := unpad(padded, 8, PKCS7)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(un, block) {
		t.Fatalf("got %v", un)
	}
}

func TestANSIX923AndISO10126(t *testing.T) {
	msg := []byte("abc")
	for _, mode := range []PaddingMode{ANSIX923, ISO10126, PKCS7, Zeros} {
		p, err := pad(msg, 8, mode)
		if err != nil {
			t.Fatal(err)
		}
		u, err := unpad(p, 8, mode)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(u, msg) {
			t.Fatalf("%s: %q", mode, u)
		}
	}
}
