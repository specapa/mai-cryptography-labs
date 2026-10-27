package main

import (
	"fmt"
	"lab2/number"
	"math/big"
)

func main() {
	svc := number.NewService()
	gcd := svc.GCD(big.NewInt(16), big.NewInt(127))
	fmt.Println(gcd)
}
