package number

import "math/big"

func (s *Service) fastIncrement(z *big.Int) {
	if bits := z.Bits(); len(bits) > 0 && bits[0] < ^big.Word(0) {
		bits[0]++
	} else {
		bigOne := new(big.Int)
		bigOne.SetInt64(1)
		z.Add(z, bigOne)
	}
}
