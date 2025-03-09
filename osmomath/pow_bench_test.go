package osmomath

import (
	"testing"

	"cosmossdk.io/math"
)

func BenchmarkPow(b *testing.B) {
	tests := []struct {
		base math.LegacyDec
		exp  math.LegacyDec
	}{
		// TODO: Choose selection here more robustly
		{
			base: math.LegacyMustNewDecFromStr("1.2"),
			exp:  math.LegacyMustNewDecFromStr("1.2"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.5"),
			exp:  math.LegacyMustNewDecFromStr("11.122"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.1"),
			exp:  math.LegacyMustNewDecFromStr("0.00000492"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.0002423"),
			exp:  math.LegacyMustNewDecFromStr("0.1234"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.493"),
			exp:  math.LegacyMustNewDecFromStr("0.00000121"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.000249"),
			exp:  math.LegacyMustNewDecFromStr("2.304"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.2342"),
			exp:  math.LegacyMustNewDecFromStr("32.2"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.000999"),
			exp:  math.LegacyMustNewDecFromStr("142.4"),
		},
		{
			base: math.LegacyMustNewDecFromStr("1.234"),
			exp:  math.LegacyMustNewDecFromStr("120.3"),
		},
		{
			base: math.LegacyMustNewDecFromStr("0.00122"),
			exp:  math.LegacyMustNewDecFromStr("123.2"),
		},
	}

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			Pow(test.base, test.exp)
		}
	}
}

func BenchmarkSqrtPow(b *testing.B) {
	tests := []struct {
		base math.LegacyDec
	}{
		// TODO: Choose selection here more robustly
		{
			base: math.LegacyMustNewDecFromStr("1.29847"),
		},
		{
			base: math.LegacyMustNewDecFromStr("1.313135"),
		},
		{
			base: math.LegacyMustNewDecFromStr("1.65976735939"),
		},
	}
	one_half := math.LegacyMustNewDecFromStr("0.5")

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			Pow(test.base, one_half)
		}
	}
}
