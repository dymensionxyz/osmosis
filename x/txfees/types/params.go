package types

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
