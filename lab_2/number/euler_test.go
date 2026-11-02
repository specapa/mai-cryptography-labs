package number

import (
	"math/big"
	"testing"
)

type eulerCase struct {
	name     string
	n, value *big.Int
}

func eulerCasesBigValues() []eulerCase {
	bigPrime, _ := new(big.Int).SetString("1000000000000037", 10)
	bigPrimeMinus1, _ := new(big.Int).SetString("1000000000000036", 10)

	bigComposite, _ := new(big.Int).SetString("123456789012345", 10)
	bigCompositeValue, _ := new(big.Int).SetString("82304496112000", 10)

	return []eulerCase{
		{name: "Большое простое число", n: bigPrime, value: bigPrimeMinus1},
		{name: "Большое составное число", n: bigComposite, value: bigCompositeValue},
	}
}

func eulerCases() []eulerCase {
	return []eulerCase{
		{name: "Граничный случай единица", n: big.NewInt(1), value: big.NewInt(1)},
		{name: "Наименьшее простое число", n: big.NewInt(2), value: big.NewInt(1)},
		{name: "Небольшое простое число", n: big.NewInt(13), value: big.NewInt(12)},
		{name: "Степень двойки (2^10)", n: big.NewInt(1024), value: big.NewInt(512)},
		{name: "Степень тройки (3^5)", n: big.NewInt(243), value: big.NewInt(162)},
		{name: "Четное составное число", n: big.NewInt(10), value: big.NewInt(4)},
		{name: "Квадрат простого числа (7^2)", n: big.NewInt(49), value: big.NewInt(42)},
		{name: "Произведение двух простых (17 * 19)", n: big.NewInt(323), value: big.NewInt(288)},
		{name: "Число с большим количеством делителей", n: big.NewInt(100), value: big.NewInt(40)},
	}
}

func testEuler_TableDriven(t *testing.T, EulerPhi func(*big.Int) (*big.Int, error), cases []eulerCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			n := new(big.Int).Set(tc.n)
			got, err := EulerPhi(n)

			if err != nil {
				t.Fatal("EulerPhi return error")
			}
			if got == nil {
				t.Fatal("EulerPhi return nil")
			}
			if got.Cmp(tc.value) != 0 {
				t.Errorf("EulerPhi return wrong answer! Got: %s; Expected: %s", got, tc.value)
			}
		})
	}
}

func TestEulerPhiDefinition(t *testing.T) {
	svc := NewService()
	testEuler_TableDriven(t, svc.EulerPhiDefinition, eulerCases())
}

func TestEulerPhiFactorization(t *testing.T) {
	svc := NewService()
	testEuler_TableDriven(t, svc.EulerPhiFactorization, eulerCases())
	testEuler_TableDriven(t, svc.EulerPhiFactorization, eulerCasesBigValues())
}

func TestEulerPhiDFT(t *testing.T) {
	svc := NewService()
	testEuler_TableDriven(t, svc.EulerPhiDFT, eulerCases())
	testEuler_TableDriven(t, svc.EulerPhiDFT, eulerCasesBigValues())
}
