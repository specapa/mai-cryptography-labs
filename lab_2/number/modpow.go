package number

import (
	"errors"
	"math/big"
)

var ErrInvalidModules = errors.New("mod must ve greater than 0")

func (s *Service) ModPowRL(base, exp, mod *big.Int) (*big.Int, error) {
	if mod.Sign() <= 0 {
		return nil, ErrInvalidModules
	}

	result := new(big.Int)
	result.SetInt64(1)

	if exp.Sign() == 0 {
		result.SetInt64(1)
		return result, nil
	}

	b := new(big.Int)
	b.Set(base)

	e := new(big.Int).Abs(exp)
	m := new(big.Int).Abs(mod)

	for e.Sign() > 0 {
		if e.Bit(0) == 1 {
			result.Mul(result, b)
			result.Mod(result, m)
		}
		b.Mul(b, b)
		b.Mod(b, m)
		e.Rsh(e, 1)
	}

	return result, nil
}

func (s *Service) ModPowLR(base, exp, mod *big.Int) (*big.Int, error) {
	if mod.Sign() <= 0 {
		return nil, ErrInvalidModules
	}

	result := new(big.Int)
	result.SetInt64(1)

	if exp.Sign() == 0 {
		result.SetInt64(1)
		return result, nil
	}

	b := new(big.Int)
	b.Set(base)

	e := new(big.Int).Abs(exp)
	m := new(big.Int).Abs(mod)

	n := exp.BitLen()
	for i := range n {
		result.Mul(result, result)
		result.Mod(result, m)

		if e.Bit(i) == 1 {
			result.Mul(result, b)
			b.Mod(b, m)
		}
	}

	return result, nil
}

func (s *Service) ModPow(base, exp, mod *big.Int) (*big.Int, error) {
	return s.ModPowLR(base, exp, mod)
}
