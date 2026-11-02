package number

import (
	"errors"
	"math/big"
)

var ErrInvalidSymbolArg = errors.New("symbol arguments are out of domain")

func (s *Service) Legendre(a, p *big.Int) (*big.Int, error) {
	panic("Undefined")
}

func (s *Service) Jacobi(a, n *big.Int) (*big.Int, error) {
	panic("Undefined")
}
