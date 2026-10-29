package number

import (
	"math/big"
	"testing"
)

type gcdCase struct {
	name string
	a, b *big.Int
}

func gcdCases() []gcdCase {
	bigA, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
	bigB, _ := new(big.Int).SetString("987654321098765432109876543210", 10)

	return []gcdCase{
		{name: "Простые взаимно простые числа", a: big.NewInt(67), b: big.NewInt(52)},
		{name: "Одно и то же простое число", a: big.NewInt(67), b: big.NewInt(67)},
		{name: "Обычные числа с общим делителем", a: big.NewInt(48), b: big.NewInt(18)},
		{name: "Одно из чисел ноль (первое)", a: big.NewInt(0), b: big.NewInt(45)},
		{name: "Одно из чисел ноль (второе)", a: big.NewInt(45), b: big.NewInt(0)},
		{name: "Оба числа нули", a: big.NewInt(0), b: big.NewInt(0)},
		{name: "Первое число отрицательное", a: big.NewInt(-48), b: big.NewInt(18)},
		{name: "Второе число отрицательное", a: big.NewInt(48), b: big.NewInt(-18)},
		{name: "Оба числа отрицательные", a: big.NewInt(-48), b: big.NewInt(-18)},
		{name: "Одно число сильно больше другого", a: big.NewInt(1000000000), b: big.NewInt(2)},
		{name: "Большие многозначные числа", a: bigA, b: bigB},
		{name: "Сбой на чётном остатке", a: big.NewInt(84), b: big.NewInt(18)},
	}
}

func referenceGCD(a, b *big.Int) *big.Int {
	return new(big.Int).GCD(nil, nil, new(big.Int).Set(a), new(big.Int).Set(b))
}

func TestGCD_TableDriven(t *testing.T) {
	svc := NewService()

	for _, tc := range gcdCases() {
		t.Run(tc.name, func(t *testing.T) {
			expected := referenceGCD(tc.a, tc.b)
			actual := svc.GCD(tc.a, tc.b)

			if actual.Cmp(expected) != 0 {
				t.Errorf("GCD(%s, %s): expected %s, got %s", tc.a, tc.b, expected, actual)
			}
		})
	}
}

func TestEGCD_TableDriven(t *testing.T) {
	svc := NewService()

	for _, tc := range gcdCases() {
		t.Run(tc.name, func(t *testing.T) {
			expected := referenceGCD(tc.a, tc.b)
			g, x, y := svc.EGCD(tc.a, tc.b)

			if g.Cmp(expected) != 0 {
				t.Errorf("EGCD(%s, %s): expected gcd %s, got %s", tc.a, tc.b, expected, g)
			}

			combo := new(big.Int).Mul(tc.a, x)
			combo.Add(combo, new(big.Int).Mul(tc.b, y))
			if combo.Cmp(expected) != 0 {
				t.Errorf("EGCD(%s, %s): %s*%s + %s*%s = %s, want %s", tc.a, tc.b, tc.a, x, tc.b, y, combo, expected)
			}
		})
	}
}

func TestBGCD_TableDriven(t *testing.T) {
	svc := NewService()

	for _, tc := range gcdCases() {
		t.Run(tc.name, func(t *testing.T) {
			expected := referenceGCD(tc.a, tc.b)
			actual := svc.BinaryGCD(tc.a, tc.b)

			if actual.Cmp(expected) != 0 {
				t.Errorf("GCD(%s, %s): expected %s, got %s", tc.a, tc.b, expected, actual)
			}
		})
	}
}
