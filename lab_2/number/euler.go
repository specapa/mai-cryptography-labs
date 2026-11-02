package number

import (
	"errors"
	"math"
	"math/big"
)

var ErrInvalidInput = errors.New("n must be greater than 0")
var ErrRepresentation = errors.New("n is not exactly representable as float64")

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

	sm := new(big.Int).SetInt64(0)
	k := new(big.Int).SetInt64(1)
	one := new(big.Int).SetInt64(1)

	for k.Cmp(n) <= 0 {
		if one.Cmp(s.BinaryGCD(k, n)) == 0 {
			s.fastIncrement(sm)
		}
		s.fastIncrement(k)
	}

	return sm, nil
}

func (s *Service) factorizationHelper(n *big.Int, m *big.Int) {
	// n / m * (m - 1)
	tmp := new(big.Int).Set(m)
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

	rest := new(big.Int).Set(n)
	result := new(big.Int).Set(n)
	p := new(big.Int).SetInt64(2)
	tmp := new(big.Int)

	for tmp.Mod(rest, p); tmp.Sign() == 0; tmp.Mod(rest, p) {
		rest.Div(rest, p)
	}

	if n.Bit(0) == 0 {
		s.factorizationHelper(result, p)
	}

	p.SetInt64(3)

	for rest.Cmp(tmp.Mul(p, p)) >= 0 {
		tmp.Mod(rest, p)
		if tmp.Sign() == 0 {
			for tmp.Sign() == 0 {
				rest.Div(rest, p)
				tmp.Mod(rest, p)
			}
			s.factorizationHelper(result, p)
		}
		s.fastAddUint(p, 2)
	}

	if rest.Cmp(one) > 0 {
		s.factorizationHelper(result, rest)
	}

	return result, nil
}

func (s *Service) EulerPhiDFT(n *big.Int) (*big.Int, error) {
	if !eulerValidateData(n) {
		return nil, ErrInvalidInput
	}

	one := new(big.Int).SetInt64(1)

	if n.Cmp(one) == 0 {
		return one, nil
	}

	var sm float64
	result := new(big.Int)
	gcd := new(big.Int)

	k := new(big.Int).Set(one)
	n_float64, accrs := n.Float64()

	if accrs != big.Exact {
		return nil, ErrRepresentation
	}

	for k.Cmp(n) <= 0 {
		gcd.Set(s.BinaryGCD(k, n))
		f_float64, _ := k.Float64()
		g_float64, _ := gcd.Float64()
		d := 2 * math.Pi * f_float64 / n_float64
		sm += g_float64 * math.Cos(d)
		s.fastIncrement(k)
	}

	result.SetInt64(int64(math.Round(sm)))
	return result, nil
}
