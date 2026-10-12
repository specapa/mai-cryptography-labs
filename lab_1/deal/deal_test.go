package deal

import (
	"bytes"
	"encoding/hex"
	"os"
	"testing"

	"lab1/des"
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

func TestAdapterMatchesDESEncrypt(t *testing.T) {
	key := mustHex(t, "133457799BBCDFF1")
	block := mustHex(t, "0123456789ABCDEF")
	want, err := encryptDES(key, block)
	if err != nil {
		t.Fatal(err)
	}
	got, err := NewDESAdapter(nil).Transform(block, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("adapter: got %X, want %X", got, want)
	}
}

func encryptDES(key, block []byte) ([]byte, error) {
	c := des.New()
	if err := c.SetKey(key); err != nil {
		return nil, err
	}
	return c.EncryptBlock(block)
}

func TestRoundTripAllKeySizes(t *testing.T) {
	plain := mustHex(t, "00112233445566778899AABBCCDDEEFF")
	keys := [][]byte{
		bytes.Repeat([]byte{0x13, 0x34, 0x57, 0x79, 0x9B, 0xBC, 0xDF, 0xF1}, 2),
		bytes.Repeat([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}, 3),
		bytes.Repeat([]byte{0x0E, 0x32, 0x92, 0x32, 0xEA, 0x6D, 0x0D, 0x73}, 4),
	}
	for _, key := range keys {
		c := New()
		if err := c.SetKey(key); err != nil {
			t.Fatalf("key %d: %v", len(key)*8, err)
		}
		ct, err := c.EncryptBlock(plain)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(ct, plain) {
			t.Fatalf("DEAL-%d left the block unchanged", len(key)*8)
		}
		pt, err := c.DecryptBlock(ct)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(pt, plain) {
			t.Fatalf("DEAL-%d round-trip: got %X want %X", len(key)*8, pt, plain)
		}
	}
}

func TestKeyScheduleLengths(t *testing.T) {
	exp := NewKeyExpander(nil)
	cases := []struct {
		keyBytes int
		rounds   int
	}{{16, 6}, {24, 6}, {32, 8}}
	for _, tc := range cases {
		keys, err := exp.ExpandKey(bytes.Repeat([]byte{0xA5}, tc.keyBytes))
		if err != nil {
			t.Fatal(err)
		}
		if len(keys) != tc.rounds {
			t.Fatalf("key %d: got %d round keys, want %d", tc.keyBytes, len(keys), tc.rounds)
		}
		for i, rk := range keys {
			if len(rk) != 8 {
				t.Fatalf("R%d length %d", i+1, len(rk))
			}
		}
	}
}

func TestRejectsBadKeyLength(t *testing.T) {
	c := New()
	if err := c.SetKey(make([]byte, 8)); err == nil {
		t.Fatal("expected error for 8-byte key")
	}
}

func TestModeRoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{0x13, 0x34, 0x57, 0x79, 0x9B, 0xBC, 0xDF, 0xF1}, 2)
	iv := mustHex(t, "0123456789ABCDEFFEDCBA9876543210")
	delta := mustHex(t, "000000000000000000000000000000A5")
	plain := []byte("MAI cryptography lab 1: DEAL round-trip payload over Feistel+DES")
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
			ctx, err := symmetric.NewSymmetricContext(New(), key, mode, pad, iv, delta)
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
	payload := bytes.Repeat([]byte("deal file payload\n"), 40)
	if err := os.WriteFile(in, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, err := symmetric.NewSymmetricContext(
		New(),
		bytes.Repeat([]byte{0x0E, 0x32, 0x92, 0x32, 0xEA, 0x6D, 0x0D, 0x73}, 2),
		symmetric.CBC,
		symmetric.PKCS7,
		mustHex(t, "A1B2C3D4E5F6071890AABBCCDDEEFF00"),
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
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(payload, got) {
		t.Fatal("file round-trip mismatch")
	}
}
