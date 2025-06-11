package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	// DYM represents 1 DYM. Equals to 10^18 base denom.
	DYM = math.NewIntWithDecimal(1, 18)

	KeyPoolCreationFee           = []byte("PoolCreationFee")
	KeyEnabledGlobalFees         = []byte("EnabledGlobalFees")
	KeyGlobalFees                = []byte("GlobalPoolFees")
	KeyTakerFees                 = []byte("TakerFees")
	KeyAllowedPoolCreationDenoms = []byte("AllowedPoolCreationDenoms")
	KeyMinSwapAmount             = []byte("MinSwapAmount")
)

// ParamTable for gamm module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// default gamm module parameters.
func DefaultParams() Params {
	return Params{
		PoolCreationFee:           sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000_000_000)},
		EnableGlobalPoolFees:      false,
		GlobalFees:                GlobalFees{math.LegacyMustNewDecFromStr("0.02"), math.LegacyZeroDec()},
		TakerFee:                  math.LegacyMustNewDecFromStr("0.01"),
		AllowedPoolCreationDenoms: []string{sdk.DefaultBondDenom},
		MinSwapAmount:             DYM, // Default to 1 DYM (10^18 base denom)
	}
}

// validate params.
func (p Params) Validate() error {
	if err := validatePoolCreationFee(p.PoolCreationFee); err != nil {
		return err
	}
	if err := validateGlobalFees(p.GlobalFees); err != nil {
		return err
	}
	if err := validateTakerFees(p.TakerFee); err != nil {
		return err
	}
	if err := validateAllowedPoolCreationDenoms(p.AllowedPoolCreationDenoms); err != nil {
		return err
	}
	if err := validateMinSwapAmount(p.MinSwapAmount); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		paramtypes.NewParamSetPair(KeyEnabledGlobalFees, &p.EnableGlobalPoolFees, func(value interface{}) error { return nil }),
		paramtypes.NewParamSetPair(KeyGlobalFees, &p.GlobalFees, validateGlobalFees),
		paramtypes.NewParamSetPair(KeyTakerFees, &p.TakerFee, validateTakerFees),
		paramtypes.NewParamSetPair(KeyAllowedPoolCreationDenoms, &p.AllowedPoolCreationDenoms, validateAllowedPoolCreationDenoms),
		paramtypes.NewParamSetPair(KeyMinSwapAmount, &p.MinSwapAmount, validateMinSwapAmount),
	}
}

func validatePoolCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid pool creation fee: %+v", i)
	}

	return nil
}

func validateGlobalFees(i interface{}) error {
	v, ok := i.(GlobalFees)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.ExitFee.IsNil() || v.SwapFee.IsNil() {
		return fmt.Errorf("invalid global pool params: %+v", i)
	}
	if v.ExitFee.IsNegative() {
		return ErrNegativeExitFee
	}

	if v.ExitFee.GTE(math.LegacyOneDec()) {
		return ErrTooMuchExitFee
	}

	if v.SwapFee.IsNegative() {
		return ErrNegativeSwapFee
	}

	if v.SwapFee.GTE(math.LegacyOneDec()) {
		return ErrTooMuchSwapFee
	}

	return nil
}

func validateTakerFees(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNil() {
		return fmt.Errorf("invalid global pool params: %+v", i)
	}
	if v.IsNegative() {
		return ErrNegativeExitFee
	}

	if v.GTE(math.LegacyOneDec()) {
		return ErrTooMuchExitFee
	}

	return nil
}

func validateAllowedPoolCreationDenoms(i interface{}) error {
	denoms, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, denom := range denoms {
		err := sdk.ValidateDenom(denom)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateMinSwapAmount(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("min swap amount cannot be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("min swap amount cannot be negative: %s", v)
	}

	return nil
}
