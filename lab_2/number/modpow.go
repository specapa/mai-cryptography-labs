package number

import "math/big"

func (s *Service) ModPow(base, exp, mod *big.Int) *big.Int {
	if mod.Sign() <= 0 {
		panic("Mod must be over 0")
	}

	panic("Not implemented")

	// b := new(big.Int).Abs(base)
	// e := new(big.Int).Abs(exp)

	// return b
}
