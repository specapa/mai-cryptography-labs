package number

import "math/big"

func (s *Service) GCD(a, b *big.Int) *big.Int {
	x := new(big.Int).Abs(a)
	y := new(big.Int).Abs(b)
	temp := new(big.Int)

	for y.Sign() != 0 {
		temp.Mod(x, y)
		x.Set(y)
		y.Set(temp)
	}

	return x
}

func (s *Service) EGCD(a, b *big.Int) (g, x, y *big.Int) {
	rem0 := new(big.Int).Abs(a)
	rem1 := new(big.Int).Abs(b)

	x0 := big.NewInt(1)
	x1 := big.NewInt(0)
	y0 := big.NewInt(0)
	y1 := big.NewInt(1)

	q := new(big.Int)
	r := new(big.Int)
	tmp := new(big.Int)

	for rem1.Sign() != 0 {
		q.DivMod(rem0, rem1, r)

		rem0.Set(rem1)
		rem1.Set(r)

		tmp.Mul(q, x1)
		tmp.Sub(x0, tmp)
		x0.Set(x1)
		x1.Set(tmp)

		tmp.Mul(q, y1)
		tmp.Sub(y0, tmp)
		y0.Set(y1)
		y1.Set(tmp)
	}

	if a.Sign() < 0 {
		x0.Neg(x0)
	}
	if b.Sign() < 0 {
		y0.Neg(y0)
	}
	return rem0, x0, y0
}

func isEven(a *big.Int) bool {
	return a.Bit(0) == 0
}

func (s *Service) BinaryGCD(a, b *big.Int) *big.Int {
	x := new(big.Int).Abs(a)
	y := new(big.Int).Abs(b)
	k := uint(0)

	for {
		if x == y {
			return x.Lsh(x, k)
		} else if x.Sign() == 0 {
			return y.Lsh(y, k)
		} else if y.Sign() == 0 {
			return x.Lsh(x, k)
		}

		for isEven(x) && isEven(y) {
			k++
			x.Rsh(x, 1)
			y.Rsh(y, 1)
		}

		for isEven(x) {
			x.Rsh(x, 1)
		}

		for isEven(y) {
			y.Rsh(y, 1)
		}

		if x.Cmp(y) > 0 {
			x.Sub(x, y)
		} else {
			y.Sub(y, x)
		}
	}
}
