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

func (s *Service) factorizationHelper(n *big.Int, m *big.Int) {
	// n / m * (m - 1)
	tmp := new(big.Int)
	tmp.Set(m)
	n.Div(n, tmp)
	s.fastDecrement(tmp)
	n.Mul(n, tmp)
}

func (s *Service) EulerPhiFactorization(n *big.Int) (*big.Int, error) {
	if !eulerValidateData(n) {
		return nil, ErrInvalidInput
	}

	one := new(big.Int)
	one.SetInt64(1)

	if n.Cmp(one) == 0 {
		return one, nil
	}

	result := new(big.Int)
	result.Set(n)

	p := new(big.Int)
	p.SetInt64(2)

	tmp := new(big.Int)

	for n.Cmp(tmp.Mul(p, p)) >= 0 {
		tmp.Mod(n, p)
		if tmp.Sign() == 0 {
			for tmp.Sign() == 0 {
				n.Div(n, p)
				tmp.Mod(n, p)
			}
			s.factorizationHelper(result, p)
		}
		s.fastIncrement(p)
	}

	if n.Cmp(one) > 0 {
		s.factorizationHelper(result, n)
	}

	return result, nil
}

func (s *Service) EulerPhiDFT(n *big.Int) (*big.Int, error) {
	if !eulerValidateData(n) {
		return nil, ErrInvalidInput
	}
	panic("Undefined")
}
