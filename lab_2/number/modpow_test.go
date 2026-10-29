package number

import (
	"math/big"
	"testing"
)

type modPowCase struct {
	name      string
	base, exp *big.Int
	mod       *big.Int
	want      *big.Int
}

func modPowCases() []modPowCase {
	bigBase, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
	bigExp, _ := new(big.Int).SetString("12345678901234567890", 10)
	bigMod, _ := new(big.Int).SetString("987654321098765432109876543211", 10)

	return []modPowCase{
		{
			name: "Обычное возведение в степень по модулю",
			base: big.NewInt(2), exp: big.NewInt(10), mod: big.NewInt(1000),
			want: big.NewInt(24),
		},
		{
			name: "Отрицательное основание",
			base: big.NewInt(-2), exp: big.NewInt(5), mod: big.NewInt(13),
			want: big.NewInt(7),
		},
		{
			name: "Нулевая степень",
			base: big.NewInt(9), exp: big.NewInt(0), mod: big.NewInt(17),
			want: big.NewInt(1),
		},
		{
			name: "Нулевое основание",
			base: big.NewInt(0), exp: big.NewInt(5), mod: big.NewInt(17),
			want: big.NewInt(0),
		},
		{
			name: "Модуль равен единице",
			base: big.NewInt(123), exp: big.NewInt(456), mod: big.NewInt(1),
			want: big.NewInt(0),
		},
		{
			name: "Большие числа",
			base: bigBase, exp: bigExp, mod: bigMod,
			want: new(big.Int).Exp(bigBase, bigExp, bigMod),
		},
	}
}

func TestModPow_TableDriven(t *testing.T) {
	svc := NewService()

	for _, tc := range modPowCases() {
		t.Run(tc.name, func(t *testing.T) {
			base := new(big.Int).Set(tc.base)
			exp := new(big.Int).Set(tc.exp)
			mod := new(big.Int).Set(tc.mod)

			got := svc.ModPow(base, exp, mod)
			if got == nil {
				t.Fatal("ModPow() returned nil")
			}
			if got.Cmp(tc.want) != 0 {
				t.Errorf("ModPow(%s, %s, %s) = %s, want %s", base, exp, mod, got, tc.want)
			}
			if got.Sign() < 0 || got.Cmp(mod) >= 0 {
				t.Errorf("ModPow() result %s is outside [0, %s)", got, mod)
			}

			if base.Cmp(tc.base) != 0 || exp.Cmp(tc.exp) != 0 || mod.Cmp(tc.mod) != 0 {
				t.Error("ModPow() modified an input argument")
			}
		})
	}
}

func TestModPow_PanicsForNonPositiveModulus(t *testing.T) {
	svc := NewService()

	for _, mod := range []*big.Int{big.NewInt(0), big.NewInt(-7)} {
		t.Run(mod.String(), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("ModPow() did not panic for modulus %s", mod)
				}
			}()

			svc.ModPow(big.NewInt(2), big.NewInt(3), mod)
		})
	}
}
