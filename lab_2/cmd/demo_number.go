package main

import (
	"fmt"
	"lab2/number"
)

func main() {
	svc := number.NewService()

	gcd := svc.GCD(12, 18)
	fmt.Println(gcd)
}
