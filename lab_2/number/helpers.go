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

func (s *Service) fastDecrement(z *big.Int) {
	if bits := z.Bits(); len(bits) > 0 && bits[0] > 0 {
		bits[0]--
	} else {
		bigOne := new(big.Int)
		bigOne.SetInt64(1)
		z.Sub(z, bigOne)
	}
}

func (s *Service) fastAddUint(z *big.Int, n big.Word) {
	if n == 0 {
		return
	}
	if bits := z.Bits(); len(bits) > 0 && bits[0] <= ^big.Word(0)-n {
		bits[0] += n
		return
	}
	z.Add(z, new(big.Int).SetUint64(uint64(n)))
}
