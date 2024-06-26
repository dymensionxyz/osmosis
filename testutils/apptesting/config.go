package apptesting

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// BaseDenomUnit is the base denom unit for the Hub
	BaseDenom            = "adym"
	AccountAddressPrefix = "dym"
)

func SetAddressPrefixes() {
	// Set prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	// Set config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

	config.SetAddressVerifier(func(bytes []byte) error {
		if len(bytes) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
		}

		if len(bytes) > address.MaxAddrLen {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bytes))
		}

		if len(bytes) != 20 && len(bytes) != 32 {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "address length must be 20 or 32 bytes, got %d", len(bytes))
		}

		return nil
	})
}
