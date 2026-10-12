package deal

import (
	"fmt"

	"lab1/des"
	"lab1/symmetric"
)

const (
	blockSize = 16
	desBlock  = 8
)

var scheduleKey = []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}

var (
	C1 = []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	C2 = []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	C3 = []byte{0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	C4 = []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type KeyExpander struct {
	NewCipher func() symmetric.SymmetricAlgorithm
}

func NewKeyExpander(newCipher func() symmetric.SymmetricAlgorithm) *KeyExpander {
	if newCipher == nil {
		newCipher = func() symmetric.SymmetricAlgorithm { return des.New() }
	}
	return &KeyExpander{NewCipher: newCipher}
}

func (e *KeyExpander) ExpandKey(key []byte) ([][]byte, error) {
	var ks [][]byte
	switch len(key) {
	case 16, 24, 32:
		for i := 0; i < len(key); i += desBlock {
			ks = append(ks, append([]byte(nil), key[i:i+desBlock]...))
		}
	default:
		return nil, fmt.Errorf("DEAL key must be 16, 24 or 32 bytes, got %d", len(key))
	}
	newCipher := e.NewCipher
	if newCipher == nil {
		newCipher = func() symmetric.SymmetricAlgorithm { return des.New() }
	}
	sched := newCipher()
	if err := sched.SetKey(scheduleKey); err != nil {
		return nil, err
	}
	enc := func(x []byte) ([]byte, error) {
		return sched.EncryptBlock(x)
	}

	var (
		roundKeys [][]byte
		err       error
	)
	switch len(ks) {
	case 2:
		roundKeys, err = expand128(ks[0], ks[1], enc)
	case 3:
		roundKeys, err = expand192(ks[0], ks[1], ks[2], enc)
	default:
		roundKeys, err = expand256(ks[0], ks[1], ks[2], ks[3], enc)
	}
	return roundKeys, err
}

func expand128(k1, k2 []byte, enc func([]byte) ([]byte, error)) ([][]byte, error) {
	r := make([][]byte, 6)
	var err error
	r[0], err = enc(k1)
	if err != nil {
		return nil, err
	}
	r[1], err = enc(xor8(k2, r[0]))
	if err != nil {
		return nil, err
	}
	r[2], err = enc(xor8(xor8(k1, r[1]), C1))
	if err != nil {
		return nil, err
	}
	r[3], err = enc(xor8(xor8(k2, r[2]), C2))
	if err != nil {
		return nil, err
	}
	r[4], err = enc(xor8(xor8(k1, r[3]), C3))
	if err != nil {
		return nil, err
	}
	r[5], err = enc(xor8(xor8(k2, r[4]), C4))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func expand192(k1, k2, k3 []byte, enc func([]byte) ([]byte, error)) ([][]byte, error) {
	r := make([][]byte, 6)
	var err error
	r[0], err = enc(k1)
	if err != nil {
		return nil, err
	}
	r[1], err = enc(xor8(k2, r[0]))
	if err != nil {
		return nil, err
	}
	r[2], err = enc(xor8(k3, r[1]))
	if err != nil {
		return nil, err
	}
	r[3], err = enc(xor8(xor8(k1, r[2]), C1))
	if err != nil {
		return nil, err
	}
	r[4], err = enc(xor8(xor8(k2, r[3]), C2))
	if err != nil {
		return nil, err
	}
	r[5], err = enc(xor8(xor8(k3, r[4]), C3))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func expand256(k1, k2, k3, k4 []byte, enc func([]byte) ([]byte, error)) ([][]byte, error) {
	r := make([][]byte, 8)
	var err error
	r[0], err = enc(k1)
	if err != nil {
		return nil, err
	}
	r[1], err = enc(xor8(k2, r[0]))
	if err != nil {
		return nil, err
	}
	r[2], err = enc(xor8(k3, r[1]))
	if err != nil {
		return nil, err
	}
	r[3], err = enc(xor8(k4, r[2]))
	if err != nil {
		return nil, err
	}
	r[4], err = enc(xor8(xor8(k1, r[3]), C1))
	if err != nil {
		return nil, err
	}
	r[5], err = enc(xor8(xor8(k2, r[4]), C2))
	if err != nil {
		return nil, err
	}
	r[6], err = enc(xor8(xor8(k3, r[5]), C3))
	if err != nil {
		return nil, err
	}
	r[7], err = enc(xor8(xor8(k4, r[6]), C4))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func xor8(a, b []byte) []byte {
	out := make([]byte, desBlock)
	for i := 0; i < desBlock; i++ {
		out[i] = a[i] ^ b[i]
	}
	return out
}

type Cipher struct {
	feistel *symmetric.FeistelNetwork
}

func New() *Cipher {
	return NewWithDES(func() symmetric.SymmetricAlgorithm { return des.New() })
}

func NewWithDES(newDES func() symmetric.SymmetricAlgorithm) *Cipher {
	return &Cipher{
		feistel: symmetric.NewFeistelNetwork(
			NewKeyExpander(newDES),
			NewDESAdapter(newDES),
			symmetric.WithBlockSize(blockSize),
			symmetric.WithFinalSwap(false),
		),
	}
}

func (c *Cipher) BlockSize() int { return blockSize }

func (c *Cipher) SetKey(key []byte) error {
	return c.feistel.SetKey(key)
}

func (c *Cipher) EncryptBlock(block []byte) ([]byte, error) {
	return c.feistel.EncryptBlock(block)
}

func (c *Cipher) DecryptBlock(block []byte) ([]byte, error) {
	return c.feistel.DecryptBlock(block)
}
