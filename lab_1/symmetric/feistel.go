package symmetric

import (
	"fmt"
)

type FeistelOption func(*FeistelNetwork)

func WithFinalSwap(swap bool) FeistelOption {
	return func(f *FeistelNetwork) { f.finalSwap = swap }
}

func WithBlockSize(n int) FeistelOption {
	return func(f *FeistelNetwork) { f.blockSize = n }
}

type FeistelNetwork struct {
	expander  KeyExpander
	transform RoundTransform
	roundKeys [][]byte
	blockSize int
	finalSwap bool
}

func NewFeistelNetwork(expander KeyExpander, transform RoundTransform, opts ...FeistelOption) *FeistelNetwork {
	f := &FeistelNetwork{
		expander:  expander,
		transform: transform,
		blockSize: 8,
		finalSwap: true,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *FeistelNetwork) BlockSize() int { return f.blockSize }

func (f *FeistelNetwork) SetKey(key []byte) error {
	if f.expander == nil {
		return fmt.Errorf("feistel: key expander is nil")
	}
	keys, err := f.expander.ExpandKey(key)
	if err != nil {
		return err
	}
	f.roundKeys = keys
	return nil
}

func (f *FeistelNetwork) EncryptBlock(block []byte) ([]byte, error) {
	return f.apply(block, f.roundKeys)
}

func (f *FeistelNetwork) DecryptBlock(block []byte) ([]byte, error) {
	rev := reversedKeys(f.roundKeys)
	if f.finalSwap {
		return f.apply(block, rev)
	}
	return f.applyInverse(block, rev)
}

func (f *FeistelNetwork) apply(block []byte, keys [][]byte) ([]byte, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("feistel: round keys are not set")
	}
	if len(block) != f.blockSize {
		return nil, fmt.Errorf("feistel: block size %d, got %d", f.blockSize, len(block))
	}
	if f.blockSize%2 != 0 {
		return nil, fmt.Errorf("feistel: block size must be even")
	}
	half := f.blockSize / 2
	L := append([]byte(nil), block[:half]...)
	R := append([]byte(nil), block[half:]...)
	for i, rk := range keys {
		fOut, err := f.transform.Transform(R, rk)
		if err != nil {
			return nil, fmt.Errorf("feistel round %d: %w", i+1, err)
		}
		if len(fOut) != half {
			return nil, fmt.Errorf("feistel round %d: F produced %d bytes, want %d", i+1, len(fOut), half)
		}
		L, R = R, xorBytes(L, fOut)
	}
	if f.finalSwap {
		L, R = R, L
	}
	return joinHalves(L, R, f.blockSize, half), nil
}

func (f *FeistelNetwork) applyInverse(block []byte, keys [][]byte) ([]byte, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("feistel: round keys are not set")
	}
	if len(block) != f.blockSize {
		return nil, fmt.Errorf("feistel: block size %d, got %d", f.blockSize, len(block))
	}
	if f.blockSize%2 != 0 {
		return nil, fmt.Errorf("feistel: block size must be even")
	}
	half := f.blockSize / 2
	L := append([]byte(nil), block[:half]...)
	R := append([]byte(nil), block[half:]...)
	for i, rk := range keys {
		fOut, err := f.transform.Transform(L, rk)
		if err != nil {
			return nil, fmt.Errorf("feistel inverse round %d: %w", i+1, err)
		}
		if len(fOut) != half {
			return nil, fmt.Errorf("feistel inverse round %d: F produced %d bytes, want %d", i+1, len(fOut), half)
		}
		L, R = xorBytes(R, fOut), L
	}
	return joinHalves(L, R, f.blockSize, half), nil
}

func joinHalves(L, R []byte, blockSize, half int) []byte {
	out := make([]byte, blockSize)
	copy(out, L)
	copy(out[half:], R)
	return out
}

func reversedKeys(keys [][]byte) [][]byte {
	out := make([][]byte, len(keys))
	for i := range keys {
		out[i] = keys[len(keys)-1-i]
	}
	return out
}

func xorBytes(a, b []byte) []byte {
	n := min(len(b), len(a))
	out := make([]byte, len(a))
	copy(out, a)
	for i := 0; i < n; i++ {
		out[i] ^= b[i]
	}
	return out
}
