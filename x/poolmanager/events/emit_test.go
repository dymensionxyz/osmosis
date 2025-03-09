package events_test

import (
	"strconv"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v15/testutils/apptesting"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	"github.com/osmosis-labs/osmosis/v15/x/poolmanager/events"
)

type PoolManagerEventsTestSuite struct {
	apptesting.KeeperTestHelper
}

const (
	addressString = "addr1---------------"
	testDenomA    = "denoma"
	testDenomB    = "denomb"
	testDenomC    = "denomc"
	testDenomD    = "denomd"
)

func TestPoolManagerEventsTestSuite(t *testing.T) {
	suite.Run(t, new(PoolManagerEventsTestSuite))
}

func (suite *PoolManagerEventsTestSuite) TestEmitSwapEvent() {
	testcases := map[string]struct {
		ctx             sdk.Context
		testAccountAddr sdk.AccAddress
		poolId          uint64
		tokensIn        sdk.Coins
		tokensOut       sdk.Coins
		closingPrice    math.LegacyDec
	}{
		"basic valid": {
			ctx:             suite.CreateTestContext(),
			testAccountAddr: sdk.AccAddress([]byte(addressString)),
			poolId:          1,
			tokensIn:        sdk.NewCoins(sdk.NewCoin(testDenomA, math.NewInt(1234))),
			tokensOut:       sdk.NewCoins(sdk.NewCoin(testDenomB, math.NewInt(5678))),
			closingPrice:    math.LegacyNewDec(123),
		},
		"valid with multiple tokens in and out": {
			ctx:             suite.CreateTestContext(),
			testAccountAddr: sdk.AccAddress([]byte(addressString)),
			poolId:          200,
			tokensIn:        sdk.NewCoins(sdk.NewCoin(testDenomA, math.NewInt(12)), sdk.NewCoin(testDenomB, math.NewInt(99))),
			tokensOut:       sdk.NewCoins(sdk.NewCoin(testDenomC, math.NewInt(88)), sdk.NewCoin(testDenomD, math.NewInt(34))),
			closingPrice:    math.LegacyNewDec(123),
		},
	}

	for name, tc := range testcases {
		suite.Run(name, func() {
			params := types.DefaultParams()
			expectedEvents := sdk.Events{
				sdk.NewEvent(
					types.TypeEvtTokenSwapped,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, tc.testAccountAddr.String()),
					sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(tc.poolId, 10)),
					sdk.NewAttribute(types.AttributeKeyTokensIn, tc.tokensIn.String()),
					sdk.NewAttribute(types.AttributeKeyTokensOut, tc.tokensOut.String()),
					sdk.NewAttribute(types.AttributeKeyClosingPrice, tc.closingPrice.String()),
					sdk.NewAttribute(types.AttributeKeyTakerFee, params.TakerFee.String()),
					sdk.NewAttribute(types.AttributeKeySwapFee, params.GlobalFees.SwapFee.String()),
				),
			}

			hasNoEventManager := tc.ctx.EventManager() == nil

			// System under test.
			events.EmitSwapEvent(tc.ctx, tc.testAccountAddr, tc.poolId, tc.tokensIn, tc.tokensOut, tc.closingPrice, params.TakerFee, params.GlobalFees.SwapFee)

			// Assertions
			if hasNoEventManager {
				// If there is no event manager on context, this is a no-op.
				return
			}

			eventManager := tc.ctx.EventManager()
			actualEvents := eventManager.Events()
			suite.Equal(expectedEvents, actualEvents)
		})
	}
}

func (suite *PoolManagerEventsTestSuite) TestEmitAddLiquidityEvent() {
	testcases := map[string]struct {
		ctx             sdk.Context
		testAccountAddr sdk.AccAddress
		poolId          uint64
		tokensIn        sdk.Coins
	}{
		"basic valid": {
			ctx:             suite.CreateTestContext(),
			testAccountAddr: sdk.AccAddress([]byte(addressString)),
			poolId:          1,
			tokensIn:        sdk.NewCoins(sdk.NewCoin(testDenomA, math.NewInt(1234))),
		},
		"valid with multiple tokens in": {
			ctx:             suite.CreateTestContext(),
			testAccountAddr: sdk.AccAddress([]byte(addressString)),
			poolId:          200,
			tokensIn:        sdk.NewCoins(sdk.NewCoin(testDenomA, math.NewInt(12)), sdk.NewCoin(testDenomB, math.NewInt(99))),
		},
	}

	for name, tc := range testcases {
		suite.Run(name, func() {
			expectedEvents := sdk.Events{
				sdk.NewEvent(
					types.TypeEvtPoolJoined,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, tc.testAccountAddr.String()),
					sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(tc.poolId, 10)),
					sdk.NewAttribute(types.AttributeKeyTokensIn, tc.tokensIn.String()),
				),
			}

			hasNoEventManager := tc.ctx.EventManager() == nil

			// System under test.
			events.EmitAddLiquidityEvent(tc.ctx, tc.testAccountAddr, tc.poolId, tc.tokensIn)

			// Assertions
			if hasNoEventManager {
				// If there is no event manager on context, this is a no-op.
				return
			}

			eventManager := tc.ctx.EventManager()
			actualEvents := eventManager.Events()
			suite.Equal(expectedEvents, actualEvents)
		})
	}
}

func (suite *PoolManagerEventsTestSuite) TestEmitRemoveLiquidityEvent() {
	testcases := map[string]struct {
		ctx             sdk.Context
		testAccountAddr sdk.AccAddress
		poolId          uint64
		tokensOut       sdk.Coins
	}{
		"basic valid": {
			ctx:             suite.CreateTestContext(),
			testAccountAddr: sdk.AccAddress([]byte(addressString)),
			poolId:          1,
			tokensOut:       sdk.NewCoins(sdk.NewCoin(testDenomA, math.NewInt(1234))),
		},
		"valid with multiple tokens out": {
			ctx:             suite.CreateTestContext(),
			testAccountAddr: sdk.AccAddress([]byte(addressString)),
			poolId:          200,
			tokensOut:       sdk.NewCoins(sdk.NewCoin(testDenomA, math.NewInt(12)), sdk.NewCoin(testDenomB, math.NewInt(99))),
		},
	}

	for name, tc := range testcases {
		suite.Run(name, func() {
			expectedEvents := sdk.Events{
				sdk.NewEvent(
					types.TypeEvtPoolExited,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, tc.testAccountAddr.String()),
					sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(tc.poolId, 10)),
					sdk.NewAttribute(types.AttributeKeyTokensOut, tc.tokensOut.String()),
				),
			}

			hasNoEventManager := tc.ctx.EventManager() == nil

			// System under test.
			events.EmitRemoveLiquidityEvent(tc.ctx, tc.testAccountAddr, tc.poolId, tc.tokensOut)

			// Assertions
			if hasNoEventManager {
				// If there is no event manager on context, this is a no-op.
				return
			}

			eventManager := tc.ctx.EventManager()
			actualEvents := eventManager.Events()
			suite.Equal(expectedEvents, actualEvents)
		})
	}
}
