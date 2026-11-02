package number

import (
	"errors"
	"math/big"
)

var ErrInvalidInput = errors.New("mod must ve greater than 0")

func eulerValidateData(n *big.Int) bool {
	if n == nil || n.Sign() < 1 {
		return false
	}
	return true
}

func (s *Service) EulerPhiDefinition(n *big.Int) (*big.Int, error) {
	if !eulerValidateData(n) {
		return nil, ErrInvalidInput
	}

	one := new(big.Int)
	one.SetInt64(1)

	if n.Cmp(one) == 0 {
		return one, nil
	}

	sm := new(big.Int)
	sm.SetInt64(0)
	k := new(big.Int)
	k.SetInt64(0)

	for k.Cmp(n) < 0 {
		if one.Cmp(s.GCD(k, n)) == 0 {
			s.fastIncrement(sm)
		}
		s.fastIncrement(k)
	}

	return sm, nil
}

func (s *Service) EulerPhiFactorization(n *big.Int) (*big.Int, error) {
	if !eulerValidateData(n) {
		return nil, ErrInvalidInput
	}
	panic("Undefined")
}

func (s *Service) EulerPhiDFT(n *big.Int) (*big.Int, error) {
	if !eulerValidateData(n) {
		return nil, ErrInvalidInput
	}
	panic("Undefined")
}
