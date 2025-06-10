package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyEpochIdentifier = []byte("EpochIdentifier")
	KeyFeeExcludeList  = []byte("FeeExcludeList")
)

func NewParams(epochIdentifier string, feeExcludeList []string) Params {
	return Params{
		EpochIdentifier: epochIdentifier,
		FeeExcludeList:  feeExcludeList,
	}
}

// default gamm module parameters.
func DefaultParams() Params {
	return Params{
		EpochIdentifier: "day",
		FeeExcludeList:  []string{},
	}
}

// ParamTable for gamm module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, validateString),
		paramtypes.NewParamSetPair(KeyFeeExcludeList, &p.FeeExcludeList, validateStringSlice),
	}
}

func validateString(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == "" {
		return fmt.Errorf("cannot be empty")
	}
	return nil
}

func validateStringSlice(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
