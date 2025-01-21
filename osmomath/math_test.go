package osmomath

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestAbsDifferenceWithSign(t *testing.T) {
	decA, err := math.LegacyNewDecFromStr("3.2")
	require.NoError(t, err)
	decB, err := math.LegacyNewDecFromStr("4.3432389")
	require.NoError(t, err)

	s, b := AbsDifferenceWithSign(decA, decB)
	require.True(t, b)

	expectedDec, err := math.LegacyNewDecFromStr("1.1432389")
	require.NoError(t, err)
	require.Equal(t, expectedDec, s)
}

func TestPowApprox(t *testing.T) {
	testCases := []struct {
		base           math.LegacyDec
		exp            math.LegacyDec
		powPrecision   math.LegacyDec
		expectedResult math.LegacyDec
		expectPanic    bool
	}{
		{
			// medium base, small exp
			base:           math.LegacyMustNewDecFromStr("0.8"),
			exp:            math.LegacyMustNewDecFromStr("0.32"),
			powPrecision:   math.LegacyMustNewDecFromStr("0.00000001"),
			expectedResult: math.LegacyMustNewDecFromStr("0.93108385"),
		},
		{
			// zero exp
			base:           math.LegacyMustNewDecFromStr("0.8"),
			exp:            sdk.ZeroDec(),
			powPrecision:   math.LegacyMustNewDecFromStr("0.00001"),
			expectedResult: math.LegacyOneDec(),
		},
		{
			// zero base, this should panic
			base:           sdk.ZeroDec(),
			exp:            math.LegacyOneDec(),
			powPrecision:   math.LegacyMustNewDecFromStr("0.00001"),
			expectedResult: sdk.ZeroDec(),
			expectPanic:    true,
		},
		{
			// large base, small exp
			base:           math.LegacyMustNewDecFromStr("1.9999"),
			exp:            math.LegacyMustNewDecFromStr("0.23"),
			powPrecision:   math.LegacyMustNewDecFromStr("0.000000001"),
			expectedResult: math.LegacyMustNewDecFromStr("1.172821461"),
		},
		{
			// large base, large integer exp
			base:           math.LegacyMustNewDecFromStr("1.777"),
			exp:            math.LegacyMustNewDecFromStr("20"),
			powPrecision:   math.LegacyMustNewDecFromStr("0.000000000001"),
			expectedResult: math.LegacyMustNewDecFromStr("98570.862372081602"),
		},
		{
			// medium base, large exp, high precision
			base:           math.LegacyMustNewDecFromStr("1.556"),
			exp:            math.LegacyMustNewDecFromStr("20.9123"),
			powPrecision:   math.LegacyMustNewDecFromStr("0.0000000000000001"),
			expectedResult: math.LegacyMustNewDecFromStr("10360.058421529811344618"),
		},
		{
			// high base, large exp, high precision
			base:           math.LegacyMustNewDecFromStr("1.886"),
			exp:            math.LegacyMustNewDecFromStr("31.9123"),
			powPrecision:   math.LegacyMustNewDecFromStr("0.00000000000001"),
			expectedResult: math.LegacyMustNewDecFromStr("621110716.84727942280335811"),
		},
		{
			// base equal one
			base:           math.LegacyMustNewDecFromStr("1"),
			exp:            math.LegacyMustNewDecFromStr("123"),
			powPrecision:   math.LegacyMustNewDecFromStr("0.00000001"),
			expectedResult: math.LegacyOneDec(),
		},
		{
			// base equal one
			base:           math.LegacyMustNewDecFromStr("1"),
			exp:            math.LegacyMustNewDecFromStr("123"),
			powPrecision:   math.LegacyMustNewDecFromStr("0.00000001"),
			expectedResult: math.LegacyOneDec(),
		},
		{
			// base equal one
			base:         math.LegacyMustNewDecFromStr("1.99999"),
			exp:          math.LegacyMustNewDecFromStr("0.1"),
			powPrecision: powPrecision,
			expectPanic:  true,
		},
		{
			// base equal one
			base:         math.LegacyMustNewDecFromStr("1.999999999999999999"),
			exp:          math.LegacyMustNewDecFromStr("0.1"),
			powPrecision: powPrecision,
			expectPanic:  true,
		},
	}

	for i, tc := range testCases {
		var actualResult math.LegacyDec
		ConditionalPanic(t, tc.expectPanic, func() {
			fmt.Println(tc.base)
			actualResult = PowApprox(tc.base, tc.exp, tc.powPrecision)
			require.True(
				t,
				tc.expectedResult.Sub(actualResult).Abs().LTE(tc.powPrecision),
				fmt.Sprintf("test %d failed: expected value & actual value's difference should be less than precision, expected: %s, actual: %s, precision: %s", i, tc.expectedResult, actualResult, tc.powPrecision),
			)
		})
	}
}

func TestPow(t *testing.T) {
	testCases := []struct {
		base           math.LegacyDec
		exp            math.LegacyDec
		expectedResult math.LegacyDec
	}{
		{
			// medium base, small exp
			base:           math.LegacyMustNewDecFromStr("0.8"),
			exp:            math.LegacyMustNewDecFromStr("0.32"),
			expectedResult: math.LegacyMustNewDecFromStr("0.93108385"),
		},
		{
			// zero exp
			base:           math.LegacyMustNewDecFromStr("0.8"),
			exp:            sdk.ZeroDec(),
			expectedResult: math.LegacyOneDec(),
		},
		{
			// zero base, this should panic
			base:           sdk.ZeroDec(),
			exp:            math.LegacyOneDec(),
			expectedResult: sdk.ZeroDec(),
		},
		{
			// large base, small exp
			base:           math.LegacyMustNewDecFromStr("1.9999"),
			exp:            math.LegacyMustNewDecFromStr("0.23"),
			expectedResult: math.LegacyMustNewDecFromStr("1.172821461"),
		},
		{
			// small base, large exp
			base:           math.LegacyMustNewDecFromStr("0.0000123"),
			exp:            math.LegacyMustNewDecFromStr("123"),
			expectedResult: sdk.ZeroDec(),
		},
		{
			// large base, large exp
			base:           math.LegacyMustNewDecFromStr("1.777"),
			exp:            math.LegacyMustNewDecFromStr("20"),
			expectedResult: math.LegacyMustNewDecFromStr("98570.862372081602"),
		},
		{
			// base equal one
			base:           math.LegacyMustNewDecFromStr("1"),
			exp:            math.LegacyMustNewDecFromStr("123"),
			expectedResult: math.LegacyOneDec(),
		},
	}

	for i, tc := range testCases {
		var actualResult math.LegacyDec
		ConditionalPanic(t, tc.base.IsZero(), func() {
			actualResult = Pow(tc.base, tc.exp)
			require.True(
				t,
				tc.expectedResult.Sub(actualResult).Abs().LTE(powPrecision),
				fmt.Sprintf("test %d failed: expected value & actual value's difference should be less than precision", i),
			)
		})
	}
}
