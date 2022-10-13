package ibc_rate_limit

import (
	"encoding/json"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/osmosis-labs/osmosis/v12/x/ibc-rate-limit/types"
)

var (
	msgSend = "send_packet"
	msgRecv = "recv_packet"
)

func CheckAndUpdateRateLimits(ctx sdk.Context, contractKeeper *wasmkeeper.PermissionedKeeper,
	msgType, contract string,
	channelValue sdk.Int, sourceChannel, denom string,
	amount string,
) error {
	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return err
	}

	sendPacketMsg, err := BuildWasmExecMsg(
		msgType,
		sourceChannel,
		denom,
		channelValue,
		amount,
	)
	if err != nil {
		return err
	}

	_, err = contractKeeper.Sudo(ctx, contractAddr, sendPacketMsg)
	if err != nil {
		return sdkerrors.Wrap(types.ErrRateLimitExceeded, err.Error())
	}

	return nil
}

type UndoSendMsg struct {
	UndoSend UndoSendMsgContent `json:"undo_send"`
}

type UndoSendMsgContent struct {
	ChannelId string `json:"channel_id"`
	Denom     string `json:"denom"`
	Funds     string `json:"funds"`
}

func UndoSendRateLimit(ctx sdk.Context, contractKeeper *wasmkeeper.PermissionedKeeper,
	contract string,
	sourceChannel, denom string,
	amount string,
) error {
	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return err
	}
	msg := UndoSendMsg{UndoSend: UndoSendMsgContent{ChannelId: sourceChannel, Denom: denom, Funds: amount}}
	asJson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = contractKeeper.Sudo(ctx, contractAddr, asJson)
	if err != nil {
		return sdkerrors.Wrap(types.ErrContractError, err.Error())
	}

	return nil
}

type SendPacketMsg struct {
	SendPacket RateLimitExecMsg `json:"send_packet"`
}

type RecvPacketMsg struct {
	RecvPacket RateLimitExecMsg `json:"recv_packet"`
}

type RateLimitExecMsg struct {
	ChannelId    string  `json:"channel_id"`
	Denom        string  `json:"denom"`
	ChannelValue sdk.Int `json:"channel_value"`
	Funds        string  `json:"funds"`
}

func BuildWasmExecMsg(msgType, sourceChannel, denom string, channelValue sdk.Int, amount string) ([]byte, error) {
	content := RateLimitExecMsg{
		ChannelId:    sourceChannel,
		Denom:        denom,
		ChannelValue: channelValue,
		Funds:        amount,
	}

	var (
		asJson []byte
		err    error
	)
	switch {
	case msgType == msgSend:
		msg := SendPacketMsg{SendPacket: content}
		asJson, err = json.Marshal(msg)
	case msgType == msgRecv:
		msg := RecvPacketMsg{RecvPacket: content}
		asJson, err = json.Marshal(msg)
	default:
		return []byte{}, types.ErrBadMessage
	}

	if err != nil {
		return []byte{}, err
	}

	return asJson, nil
}

func GetLocalDenom(packet exported.PacketI, data transfertypes.FungibleTokenPacketData) string {
	var denomTrace transfertypes.DenomTrace
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := data.Denom[len(voucherPrefix):]
		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace = transfertypes.ParseDenomTrace(unprefixedDenom)
	} else {
		// since SendPacket did not prefix the denomination, we must prefix denomination here
		sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedDenom := sourcePrefix + data.Denom

		// construct the denomination trace from the full raw denomination
		denomTrace = transfertypes.ParseDenomTrace(prefixedDenom)
	}

	return denomTrace.IBCDenom()
}

func GetFundsFromPacket(packet exported.PacketI, local bool) (string, string, error) {
	var packetData transfertypes.FungibleTokenPacketData
	err := json.Unmarshal(packet.GetData(), &packetData)
	if err != nil {
		return "", "", err
	}
	if local {
		return packetData.Amount, packetData.Denom, nil
	}
	return packetData.Amount, GetLocalDenom(packet, packetData), nil
}
