package des

import (
	"bytes"
	"encoding/hex"
	"testing"

	"lab1/permutation"
	"lab1/symmetric"
)

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestIPThenFPIsIdentity(t *testing.T) {
	in := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	mid, err := permutation.Permute(in, IP, permutation.MSBToLSB, permutation.StartFromOne)
	if err != nil {
		t.Fatal(err)
	}
	out, err := permutation.Permute(mid, FP, permutation.MSBToLSB, permutation.StartFromOne)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(in, out) {
		t.Fatalf("IP∘FP is not identity: %x -> %x", in, out)
	}
}

func TestDESKnownVectorTextbook(t *testing.T) {
	c := New()
	if err := c.SetKey(mustHex(t, "133457799BBCDFF1")); err != nil {
		t.Fatal(err)
	}
	pt := mustHex(t, "0123456789ABCDEF")
	want := mustHex(t, "85E813540F0AB405")
	got, err := c.EncryptBlock(pt)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("encrypt: got %X, want %X", got, want)
	}
	back, err := c.DecryptBlock(got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(back, pt) {
		t.Fatalf("decrypt: got %X, want %X", back, pt)
	}
}

func TestDESKnownVectorNBS(t *testing.T) {
	c := New()
	if err := c.SetKey(mustHex(t, "0101010101010101")); err != nil {
		t.Fatal(err)
	}
	pt := mustHex(t, "0000000000000000")
	want := mustHex(t, "8CA64DE9C1B123A7")
	got, err := c.EncryptBlock(pt)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("encrypt: got %X, want %X", got, want)
	}
}

func TestDESFIPS81(t *testing.T) {
	c := New()
	if err := c.SetKey(mustHex(t, "0123456789ABCDEF")); err != nil {
		t.Fatal(err)
	}
	pt := []byte("Now is t")
	want := mustHex(t, "3FA40E8A984D4815")
	got, err := c.EncryptBlock(pt)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("encrypt: got %X, want %X", got, want)
	}
}

func TestSBoxLectureExample(t *testing.T) {
	chunk := byte(0b101100)
	row := ((chunk & 0b100000) >> 4) | (chunk & 0b000001)
	col := (chunk >> 1) & 0b1111
	if row != 2 || col != 6 {
		t.Fatalf("row=%d col=%d, want 2,6", row, col)
	}
	if SBoxes[0][row*16+col] != 2 {
		t.Fatalf("S1[2][6]=%d, FIPS value is 2", SBoxes[0][row*16+col])
	}
}

func TestModeRoundTrip(t *testing.T) {
	key := mustHex(t, "133457799BBCDFF1")
	iv := mustHex(t, "0123456789ABCDEF")
	delta := mustHex(t, "00000000000000A5")
	plain := []byte("MAI cryptography lab 1: Feistel/DES round-trip payload!")
	modes := []symmetric.CipherMode{
		symmetric.ECB, symmetric.CBC, symmetric.PCBC,
		symmetric.CFB, symmetric.OFB, symmetric.CTR, symmetric.RandomDelta,
	}
	pads := []symmetric.PaddingMode{
		symmetric.PKCS7, symmetric.ANSIX923, symmetric.ISO10126, symmetric.Zeros,
	}
	for _, mode := range modes {
		for _, pad := range pads {
			msg := plain
			if pad == symmetric.Zeros {
				msg = []byte("no-trailing-zero-payload")
			}
			ciph := New()
			ctx, err := symmetric.NewSymmetricContext(ciph, key, mode, pad, iv, delta)
			if err != nil {
				t.Fatalf("%s/%s: %v", mode, pad, err)
			}
			var enc, dec []byte
			if err := ctx.Encrypt(msg, &enc); err != nil {
				t.Fatalf("%s/%s encrypt: %v", mode, pad, err)
			}
			if err := ctx.Decrypt(enc, &dec); err != nil {
				t.Fatalf("%s/%s decrypt: %v", mode, pad, err)
			}
			if !bytes.Equal(msg, dec) {
				t.Fatalf("%s/%s round-trip failed\nplain=%q\ndec  =%q", mode, pad, msg, dec)
			}
		}
	}
}

func TestFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	in := dir + "/in.bin"
	enc := dir + "/enc.bin"
	out := dir + "/out.bin"
	payload := bytes.Repeat([]byte("rb-tree allocator test source\n"), 50)
	if err := writeFile(in, payload); err != nil {
		t.Fatal(err)
	}
	ciph := New()
	ctx, err := symmetric.NewSymmetricContext(
		ciph,
		mustHex(t, "0E329232EA6D0D73"),
		symmetric.CBC,
		symmetric.PKCS7,
		mustHex(t, "A1B2C3D4E5F60718"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := ctx.EncryptFile(in, enc); err != nil {
		t.Fatal(err)
	}
	if err := ctx.DecryptFile(enc, out); err != nil {
		t.Fatal(err)
	}
	got, err := readFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(payload, got) {
		t.Fatal("file round-trip mismatch")
	}
}
