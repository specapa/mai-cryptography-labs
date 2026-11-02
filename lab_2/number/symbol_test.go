package number

import (
	"errors"
	"math/big"
	"testing"
)

type symbolCase struct {
	name       string
	a, n, want *big.Int
}

func legendreCases() []symbolCase {
	bigPrime, _ := new(big.Int).SetString("1000000000000037", 10)

	return []symbolCase{
		{name: "Квадратичный вычет", a: big.NewInt(4), n: big.NewInt(5), want: big.NewInt(1)},
		{name: "Квадратичный невычет", a: big.NewInt(3), n: big.NewInt(5), want: big.NewInt(-1)},
		{name: "Числитель делится на модуль", a: big.NewInt(5), n: big.NewInt(5), want: big.NewInt(0)},
		{name: "Числитель больше модуля", a: big.NewInt(6), n: big.NewInt(5), want: big.NewInt(1)},
		{name: "Отрицательный числитель, p = 1 mod 4", a: big.NewInt(-1), n: big.NewInt(5), want: big.NewInt(1)},
		{name: "Отрицательный числитель, p = 3 mod 4", a: big.NewInt(-1), n: big.NewInt(7), want: big.NewInt(-1)},
		{name: "Двойка при p = 7 mod 8", a: big.NewInt(2), n: big.NewInt(7), want: big.NewInt(1)},
		{name: "Двойка при p = 5 mod 8", a: big.NewInt(2), n: big.NewInt(5), want: big.NewInt(-1)},
		{name: "Нулевой числитель", a: big.NewInt(0), n: big.NewInt(7), want: big.NewInt(0)},
		{name: "Вычет по модулю 13", a: big.NewInt(10), n: big.NewInt(13), want: big.NewInt(1)},
		{name: "Большое простое число", a: big.NewInt(2), n: bigPrime, want: big.NewInt(-1)},
	}
}

func jacobiCases() []symbolCase {
	return []symbolCase{
		{name: "Модуль единица", a: big.NewInt(1), n: big.NewInt(1), want: big.NewInt(1)},
		{name: "Ноль по модулю единицы", a: big.NewInt(0), n: big.NewInt(1), want: big.NewInt(1)},
		{name: "Числитель больше единицы", a: big.NewInt(8), n: big.NewInt(1), want: big.NewInt(1)},
		{name: "Символ равен 1 при составном модуле", a: big.NewInt(2), n: big.NewInt(15), want: big.NewInt(1)},
		{name: "Символ равен -1", a: big.NewInt(7), n: big.NewInt(15), want: big.NewInt(-1)},
		{name: "Общий делитель с модулем", a: big.NewInt(3), n: big.NewInt(9), want: big.NewInt(0)},
		{name: "Степень простого числа", a: big.NewInt(2), n: big.NewInt(9), want: big.NewInt(1)},
		{name: "Произведение двух простых", a: big.NewInt(5), n: big.NewInt(21), want: big.NewInt(1)},
		{name: "Отрицательный числитель", a: big.NewInt(-1), n: big.NewInt(15), want: big.NewInt(-1)},
		{name: "Числитель больше модуля", a: big.NewInt(17), n: big.NewInt(15), want: big.NewInt(1)},
		{name: "Составной модуль", a: big.NewInt(19), n: big.NewInt(45), want: big.NewInt(1)},
		{name: "Ноль и составной модуль", a: big.NewInt(0), n: big.NewInt(15), want: big.NewInt(0)},
	}
}

func testSymbolTable(t *testing.T, symbol func(*big.Int, *big.Int) (*big.Int, error), cases []symbolCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := new(big.Int).Set(tc.a)
			n := new(big.Int).Set(tc.n)
			got, err := symbol(a, n)

			if err != nil {
				t.Fatalf("%s return error: %v", t.Name(), err)
			}
			if got == nil {
				t.Fatalf("%s return nil", t.Name())
			}
			if got.Cmp(tc.want) != 0 {
				t.Fatalf("%s return wrong answer! Got: %s; Expected: %s", t.Name(), got, tc.want)
			}
		})
	}
}

func TestLegendre(t *testing.T) {
	svc := NewService()
	testSymbolTable(t, svc.Legendre, legendreCases())
}

func TestJacobi(t *testing.T) {
	svc := NewService()
	testSymbolTable(t, svc.Jacobi, jacobiCases())
}

func TestSymbolInvalidArg(t *testing.T) {
	svc := NewService()
	methods := []struct {
		name string
		fn   func(*big.Int, *big.Int) (*big.Int, error)
	}{
		{name: "Legendre", fn: svc.Legendre},
		{name: "Jacobi", fn: svc.Jacobi},
	}
	inputs := []struct {
		name string
		a, n *big.Int
	}{
		{name: "nil числитель", a: nil, n: big.NewInt(5)},
		{name: "nil модуль", a: big.NewInt(2), n: nil},
		{name: "нулевой модуль", a: big.NewInt(2), n: big.NewInt(0)},
		{name: "отрицательный модуль", a: big.NewInt(2), n: big.NewInt(-5)},
		{name: "чётный модуль", a: big.NewInt(3), n: big.NewInt(2)},
	}

	for _, m := range methods {
		for _, in := range inputs {
			t.Run(m.name+"/"+in.name, func(t *testing.T) {
				got, err := m.fn(in.a, in.n)
				if !errors.Is(err, ErrInvalidSymbolArg) {
					t.Fatalf("got error %v, expected ErrInvalidSymbolArg", err)
				}
				if got != nil {
					t.Fatalf("got %v, expected nil", got)
				}
			})
		}
	}
}
