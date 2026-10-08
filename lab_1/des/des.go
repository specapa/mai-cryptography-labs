package des

import (
	"fmt"

	"lab1/permutation"
	"lab1/symmetric"
)

const (
	blockSize = 8
	keySize   = 8
	halfSize  = 4
)

type KeyExpander struct{}

func (KeyExpander) ExpandKey(key []byte) ([][]byte, error) {
	if len(key) != keySize {
		return nil, fmt.Errorf("DES key must be %d bytes, got %d", keySize, len(key))
	}
	cd, err := permutation.Permute(key, PC1, permutation.MSBToLSB, permutation.StartFromOne)
	if err != nil {
		return nil, err
	}
	c := high28(cd)
	d := low28(cd)
	keys := make([][]byte, 16)
	for i := 0; i < 16; i++ {
		c = rotl28(c, Shifts[i])
		d = rotl28(d, Shifts[i])
		packed := pack56(c, d)
		rk, err := permutation.Permute(packed, PC2, permutation.MSBToLSB, permutation.StartFromOne)
		if err != nil {
			return nil, err
		}
		keys[i] = rk
	}
	return keys, nil
}

type RoundFunction struct{}

func (RoundFunction) Transform(block, roundKey []byte) ([]byte, error) {
	if len(block) != halfSize {
		return nil, fmt.Errorf("DES F expects 4-byte half-block, got %d", len(block))
	}
	if len(roundKey) != 6 {
		return nil, fmt.Errorf("DES round key must be 6 bytes, got %d", len(roundKey))
	}
	expanded, err := permutation.Permute(block, E, permutation.MSBToLSB, permutation.StartFromOne)
	if err != nil {
		return nil, err
	}
	xored := xor6(expanded, roundKey)
	sOut := make([]byte, 4)
	for k := 0; k < 8; k++ {
		chunk, err := sixBits(xored, k)
		if err != nil {
			return nil, err
		}
		row := ((chunk & 0b100000) >> 4) | (chunk & 0b000001)
		col := (chunk >> 1) & 0b1111
		val := SBoxes[k][row*16+col]
		if err := writeNibble(sOut, k, val); err != nil {
			return nil, err
		}
	}
	return permutation.Permute(sOut, P, permutation.MSBToLSB, permutation.StartFromOne)
}

type Cipher struct {
	feistel *symmetric.FeistelNetwork
}

func New() *Cipher {
	return &Cipher{
		feistel: symmetric.NewFeistelNetwork(
			KeyExpander{},
			RoundFunction{},
			symmetric.WithBlockSize(blockSize),
			symmetric.WithFinalSwap(true),
		),
	}
}

func (c *Cipher) BlockSize() int { return blockSize }

func (c *Cipher) SetKey(key []byte) error {
	return c.feistel.SetKey(key)
}

func (c *Cipher) EncryptBlock(block []byte) ([]byte, error) {
	return c.cryptBlock(block, c.feistel.EncryptBlock)
}

func (c *Cipher) DecryptBlock(block []byte) ([]byte, error) {
	return c.cryptBlock(block, c.feistel.DecryptBlock)
}

func (c *Cipher) cryptBlock(block []byte, inner func([]byte) ([]byte, error)) ([]byte, error) {
	if len(block) != blockSize {
		return nil, fmt.Errorf("DES block must be %d bytes, got %d", blockSize, len(block))
	}
	ip, err := permutation.Permute(block, IP, permutation.MSBToLSB, permutation.StartFromOne)
	if err != nil {
		return nil, err
	}
	mid, err := inner(ip)
	if err != nil {
		return nil, err
	}
	return permutation.Permute(mid, FP, permutation.MSBToLSB, permutation.StartFromOne)
}

func high28(b []byte) uint32 {
	return uint32(b[0])<<20 | uint32(b[1])<<12 | uint32(b[2])<<4 | uint32(b[3]>>4)
}

func low28(b []byte) uint32 {
	return uint32(b[3]&0x0F)<<24 | uint32(b[4])<<16 | uint32(b[5])<<8 | uint32(b[6])
}

func pack56(c, d uint32) []byte {
	c &= 0x0FFFFFFF
	d &= 0x0FFFFFFF
	b := make([]byte, 7)
	b[0] = byte(c >> 20)
	b[1] = byte(c >> 12)
	b[2] = byte(c >> 4)
	b[3] = byte((c&0xF)<<4) | byte(d>>24)
	b[4] = byte(d >> 16)
	b[5] = byte(d >> 8)
	b[6] = byte(d)
	return b
}

func rotl28(x uint32, n uint) uint32 {
	x &= 0x0FFFFFFF
	return ((x << n) | (x >> (28 - n))) & 0x0FFFFFFF
}

func xor6(a, b []byte) []byte {
	out := make([]byte, 6)
	for i := 0; i < 6; i++ {
		out[i] = a[i] ^ b[i]
	}
	return out
}

func sixBits(data []byte, group int) (byte, error) {
	var v byte
	base := group * 6
	for i := 0; i < 6; i++ {
		bit, err := permutation.Bit(data, 1+base+i, permutation.MSBToLSB, permutation.StartFromOne)
		if err != nil {
			return 0, err
		}
		v = (v << 1) | bit
	}
	return v, nil
}

func writeNibble(dst []byte, group int, val byte) error {
	for i := 0; i < 4; i++ {
		bit := (val >> (3 - i)) & 1
		if err := permutation.SetBit(dst, 1+group*4+i, bit, permutation.MSBToLSB, permutation.StartFromOne); err != nil {
			return err
		}
	}
	return nil
}
