package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyEpochIdentifier = []byte("EpochIdentifier")
)

// ParamTable for gamm module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(epochIdentifier string) Params {
	return Params{
		EpochIdentifier: epochIdentifier,
	}
}

// default gamm module parameters.
func DefaultParams() Params {
	return Params{
		EpochIdentifier: "day",
	}
}

// validate params.
func (p Params) Validate() error {
	return validateString(p.EpochIdentifier)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEpochIdentifier, &p.EpochIdentifier, validateString),
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
