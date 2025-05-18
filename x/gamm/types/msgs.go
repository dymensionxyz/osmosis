package types

import (
	"fmt"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSwapExactAmountIn{}
	_ sdk.Msg = &MsgSwapExactAmountOut{}
	_ sdk.Msg = &MsgJoinPool{}
	_ sdk.Msg = &MsgExitPool{}
	_ sdk.Msg = &MsgJoinSwapExternAmountIn{}
	_ sdk.Msg = &MsgJoinSwapShareAmountOut{}
	_ sdk.Msg = &MsgExitSwapExternAmountOut{}
	_ sdk.Msg = &MsgExitSwapShareAmountIn{}
)

func ValidateFutureGovernor(governor string) error {
	// allow empty governor
	if governor == "" {
		return nil
	}

	// validation for future owner
	_, err := sdk.AccAddressFromBech32(governor)
	if err == nil {
		return nil
	}

	lockTimeStr := ""
	splits := strings.Split(governor, ",")
	if len(splits) > 2 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid future governor: %s", governor))
	}

	// token,100h
	if len(splits) == 2 {
		lpTokenStr := splits[0]
		if sdk.ValidateDenom(lpTokenStr) != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid future governor: %s", governor))
		}
		lockTimeStr = splits[1]
	}

	// 100h
	if len(splits) == 1 {
		lockTimeStr = splits[0]
	}

	// Note that a duration of 0 is allowed
	_, err = time.ParseDuration(lockTimeStr)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid future governor: %s", governor))
	}
	return nil
}

func (msg MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = SwapAmountInRoutes(msg.Routes).Validate()
	if err != nil {
		return err
	}

	if !msg.TokenIn.IsValid() || !msg.TokenIn.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenIn.String())
	}

	if !msg.TokenOutMinAmount.IsPositive() {
		return ErrNotPositiveCriteria
	}

	return nil
}

func (msg MsgSwapExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = SwapAmountOutRoutes(msg.Routes).Validate()
	if err != nil {
		return err
	}

	if !msg.TokenOut.IsValid() || !msg.TokenOut.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenOut.String())
	}

	if !msg.TokenInMaxAmount.IsPositive() {
		return ErrNotPositiveCriteria
	}

	return nil
}

func (msg MsgJoinPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.ShareOutAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveRequireAmount, msg.ShareOutAmount.String())
	}

	tokenInMaxs := sdk.Coins(msg.TokenInMaxs)
	if !tokenInMaxs.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, tokenInMaxs.String())
	}

	return nil
}

func (msg MsgExitPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.ShareInAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveRequireAmount, msg.ShareInAmount.String())
	}

	tokenOutMins := sdk.Coins(msg.TokenOutMins)
	if !tokenOutMins.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, tokenOutMins.String())
	}

	return nil
}

func (msg MsgJoinSwapExternAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.TokenIn.IsValid() || !msg.TokenIn.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenIn.String())
	}

	if !msg.ShareOutMinAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveCriteria, msg.ShareOutMinAmount.String())
	}

	return nil
}

func (msg MsgJoinSwapShareAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = sdk.ValidateDenom(msg.TokenInDenom)
	if err != nil {
		return err
	}

	if !msg.ShareOutAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveRequireAmount, msg.ShareOutAmount.String())
	}

	if !msg.TokenInMaxAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveCriteria, msg.TokenInMaxAmount.String())
	}

	return nil
}

func (msg MsgExitSwapExternAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !msg.TokenOut.IsValid() || !msg.TokenOut.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, msg.TokenOut.String())
	}

	if !msg.ShareInMaxAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveCriteria, msg.ShareInMaxAmount.String())
	}

	return nil
}

func (msg MsgExitSwapShareAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = sdk.ValidateDenom(msg.TokenOutDenom)
	if err != nil {
		return err
	}

	if !msg.ShareInAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveRequireAmount, msg.ShareInAmount.String())
	}

	if !msg.TokenOutMinAmount.IsPositive() {
		return errorsmod.Wrap(ErrNotPositiveCriteria, msg.TokenOutMinAmount.String())
	}

	return nil
}
