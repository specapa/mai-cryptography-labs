package number

import (
	"errors"
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

	bigComposite_0, _ := new(big.Int).SetString("123456789012345", 10)
	bigCompositeValue_0, _ := new(big.Int).SetString("65593911768000", 10)

	bigComposite_1, _ := new(big.Int).SetString("67676767676767676767", 10)
	bigCompositeValue_1, _ := new(big.Int).SetString("64131589399680000000", 10)

	return []eulerCase{
		{name: "Большое простое число", n: bigPrime, value: bigPrimeMinus1},
		{name: "Большое составное число 1", n: bigComposite_0, value: bigCompositeValue_0},
		{name: "Большое составное число 2", n: bigComposite_1, value: bigCompositeValue_1},
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
				t.Fatalf("%s return error", t.Name())
			}
			if got == nil {
				t.Fatalf("%s return nil", t.Name())
			}
			if got.Cmp(tc.value) != 0 {
				t.Fatalf("%s return wrong answer! Got: %s; Expected: %s", t.Name(), got, tc.value)
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
}

func TestEulerPhiFactorizationBigValues(t *testing.T) {
	svc := NewService()
	testEuler_TableDriven(t, svc.EulerPhiFactorization, eulerCasesBigValues())
}

func TestEulerPhiDFT(t *testing.T) {
	svc := NewService()
	testEuler_TableDriven(t, svc.EulerPhiDFT, eulerCases())
}

func TestEulerPhiInvalidInput(t *testing.T) {
	svc := NewService()
	methods := []struct {
		name string
		fn   func(*big.Int) (*big.Int, error)
	}{
		{name: "Definition", fn: svc.EulerPhiDefinition},
		{name: "Factorization", fn: svc.EulerPhiFactorization},
		{name: "DFT", fn: svc.EulerPhiDFT},
	}
	inputs := []struct {
		name string
		n    *big.Int
	}{
		{name: "nil", n: nil},
		{name: "ноль", n: big.NewInt(0)},
		{name: "отрицательное", n: big.NewInt(-15)},
	}

	for _, m := range methods {
		for _, in := range inputs {
			t.Run(m.name+"/"+in.name, func(t *testing.T) {
				got, err := m.fn(in.n)
				if !errors.Is(err, ErrInvalidInput) {
					t.Fatalf("got error %v, expected ErrInvalidInput", err)
				}
				if got != nil {
					t.Fatalf("got %v, expected nil", got)
				}
			})
		}
	}
}

func TestEulerPhiDFTRepresentation(t *testing.T) {
	svc := NewService()
	n, ok := new(big.Int).SetString("67676767676767676767", 10)
	if !ok {
		t.Fatal("failed to parse n")
	}

	got, err := svc.EulerPhiDFT(n)
	if !errors.Is(err, ErrRepresentation) {
		t.Fatalf("got error %v, expected ErrRepresentation", err)
	}
	if got != nil {
		t.Fatalf("got %v, expected nil", got)
	}
}
