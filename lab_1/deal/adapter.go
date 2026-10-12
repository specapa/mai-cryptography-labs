package deal

import (
	"fmt"

	"lab1/des"
	"lab1/symmetric"
)

type DESAdapter struct {
	NewCipher func() symmetric.SymmetricAlgorithm
}

func NewDESAdapter(newCipher func() symmetric.SymmetricAlgorithm) *DESAdapter {
	if newCipher == nil {
		newCipher = func() symmetric.SymmetricAlgorithm { return des.New() }
	}
	return &DESAdapter{NewCipher: newCipher}
}

func (a *DESAdapter) Transform(half, roundKey []byte) ([]byte, error) {
	if len(half) != desBlock {
		return nil, fmt.Errorf("DEAL F expects %d-byte half-block, got %d", desBlock, len(half))
	}
	if len(roundKey) != desBlock {
		return nil, fmt.Errorf("DEAL round key must be %d bytes, got %d", desBlock, len(roundKey))
	}
	ciph := a.NewCipher()
	if ciph == nil {
		return nil, fmt.Errorf("DEAL adapter: cipher factory returned nil")
	}
	if err := ciph.SetKey(roundKey); err != nil {
		return nil, err
	}
	return ciph.EncryptBlock(half)
}
