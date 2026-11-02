package number

import (
	"errors"
	"math/big"
	"testing"
)

func carmichaelCases() []eulerCase {
	return []eulerCase{
		{name: "Граничный случай единица", n: big.NewInt(1), value: big.NewInt(1)},
		{name: "Двойка", n: big.NewInt(2), value: big.NewInt(1)},
		{name: "Четыре", n: big.NewInt(4), value: big.NewInt(2)},
		{name: "Восемь", n: big.NewInt(8), value: big.NewInt(2)},
		{name: "Степень двойки (2^10)", n: big.NewInt(1024), value: big.NewInt(256)},
		{name: "Нечётное простое число", n: big.NewInt(7), value: big.NewInt(6)},
		{name: "Степень тройки (3^5)", n: big.NewInt(243), value: big.NewInt(162)},
		{name: "Степень пятерки (5^2)", n: big.NewInt(25), value: big.NewInt(20)},
		{name: "Произведение двух простых (3 * 5)", n: big.NewInt(15), value: big.NewInt(4)},
		{name: "Произведение двух простых (17 * 19)", n: big.NewInt(323), value: big.NewInt(144)},
		{name: "Число с большим количеством делителей", n: big.NewInt(100), value: big.NewInt(20)},
	}
}

func carmichaelCasesBigValues() []eulerCase {
	bigPrime, _ := new(big.Int).SetString("1000000000000037", 10)
	bigPrimeMinus1, _ := new(big.Int).SetString("1000000000000036", 10)

	bigComposite, _ := new(big.Int).SetString("123456789012345", 10)
	bigCompositeValue, _ := new(big.Int).SetString("819923897100", 10)

	bigSmooth, _ := new(big.Int).SetString("67676767676767676767", 10)
	bigSmoothValue, _ := new(big.Int).SetString("82473751800", 10)

	return []eulerCase{
		{name: "Большое простое число", n: bigPrime, value: bigPrimeMinus1},
		{name: "Большое составное число", n: bigComposite, value: bigCompositeValue},
		{name: "Произведение небольших простых", n: bigSmooth, value: bigSmoothValue},
	}
}

func TestCarmichael(t *testing.T) {
	svc := NewService()
	testEuler_TableDriven(t, svc.Carmichael, carmichaelCases())
}

func TestCarmichaelBigValues(t *testing.T) {
	svc := NewService()
	testEuler_TableDriven(t, svc.Carmichael, carmichaelCasesBigValues())
}

func TestCarmichaelInvalidInput(t *testing.T) {
	svc := NewService()
	inputs := []struct {
		name string
		n    *big.Int
	}{
		{name: "nil", n: nil},
		{name: "ноль", n: big.NewInt(0)},
		{name: "отрицательное", n: big.NewInt(-15)},
	}

	for _, in := range inputs {
		t.Run(in.name, func(t *testing.T) {
			got, err := svc.Carmichael(in.n)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("got error %v, expected ErrInvalidInput", err)
			}
			if got != nil {
				t.Fatalf("got %v, expected nil", got)
			}
		})
	}
}
