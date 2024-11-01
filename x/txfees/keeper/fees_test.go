package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/osmosis-labs/osmosis/v15/testutils/apptesting"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

func (s *KeeperTestSuite) TestChargeFees() {
	accs := apptesting.CreateRandomAccounts(2)

	testCases := map[string]struct {
		payer             sdk.AccAddress
		takerFee          sdk.Coin
		beneficiary       *sdk.AccAddress
		expTakerFee       sdk.Coins
		expBeneficiaryRev sdk.Coins
	}{
		"with beneficiary": {
			payer:             accs[0],
			takerFee:          sdk.NewCoin("adym", sdk.NewInt(100)),
			beneficiary:       &accs[1],
			expTakerFee:       sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(50))),
			expBeneficiaryRev: sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(50))),
		},
		"without beneficiary": {
			payer:             accs[0],
			takerFee:          sdk.NewCoin("adym", sdk.NewInt(100)),
			beneficiary:       nil,
			expTakerFee:       sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(100))),
			expBeneficiaryRev: nil,
		},
		"base denom fee": {
			payer:             accs[0],
			takerFee:          sdk.NewCoin("adym", sdk.NewInt(100)),
			beneficiary:       &accs[1],
			expTakerFee:       sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(50))),
			expBeneficiaryRev: sdk.NewCoins(sdk.NewCoin("adym", sdk.NewInt(50))),
		},
		"fee with fee token": {
			payer:             accs[0],
			takerFee:          sdk.NewCoin("foo", sdk.NewInt(100)),
			beneficiary:       &accs[1],
			expTakerFee:       sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(50))),
			expBeneficiaryRev: sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(50))),
		},
		"fee without fee token": {
			payer:             accs[0],
			takerFee:          sdk.NewCoin("baz", sdk.NewInt(100)),
			beneficiary:       &accs[1],
			expTakerFee:       sdk.NewCoins(sdk.NewCoin("baz", sdk.NewInt(50))),
			expBeneficiaryRev: sdk.NewCoins(sdk.NewCoin("baz", sdk.NewInt(50))),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			s.SetupTest()

			// Create base denom and prepare pools
			//
			// Base denom: adym
			// Fee denoms: foo, bar
			// Pools:
			//  - adym <-> foo
			//  - bar  <-> foo
			//  - bar  <-> adym
			//  - bar  <-> baz

			s.FundAcc(tc.payer, sdk.NewCoins(tc.takerFee))

			err := s.App.TxFeesKeeper.SetBaseDenom(s.Ctx, "adym")
			s.Require().NoError(err)

			pool1coins := []sdk.Coin{sdk.NewCoin("adym", sdk.NewInt(100000)), sdk.NewCoin("foo", sdk.NewInt(100000))}
			s.PrepareBalancerPoolWithCoins(pool1coins...)

			pool2coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000)), sdk.NewCoin("foo", sdk.NewInt(100000))}
			s.PrepareBalancerPoolWithCoins(pool2coins...)

			pool3coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000)), sdk.NewCoin("adym", sdk.NewInt(100000))}
			s.PrepareBalancerPoolWithCoins(pool3coins...)

			pool4coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000)), sdk.NewCoin("baz", sdk.NewInt(100000))}
			s.PrepareBalancerPoolWithCoins(pool4coins...)

			initialTxFeesBalance := s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.AccountKeeper.GetModuleAddress(types.ModuleName))

			// Reset event counts to 0 by creating a new manager.
			s.Ctx = s.Ctx.WithEventManager(sdk.NewEventManager())

			// Charge fees
			err = s.App.TxFeesKeeper.ChargeFeesFromPayer(s.Ctx, tc.payer, tc.takerFee, tc.beneficiary)
			s.Require().NoError(err)

			// Verify results

			// Verify charge fee event
			eventName := proto.MessageName(new(types.EventChargeFee))
			s.AssertEventEmitted(s.Ctx, eventName, 1)
			event := s.ExtractChargeFeeEvent(s.Ctx.EventManager().Events(), eventName)
			s.Require().Equal(tc.payer.String(), event.Payer)
			s.Require().Equal(tc.expTakerFee.String(), event.TakerFee)
			if tc.beneficiary != nil {
				s.Require().Equal(tc.beneficiary.String(), event.Beneficiary)
				s.Require().Equal(tc.expBeneficiaryRev.String(), event.BeneficiaryRevenue)
			}

			// The fee is either burned or not applied if case of error
			actualTxFeesBalance := s.App.BankKeeper.GetAllBalances(s.Ctx, s.App.AccountKeeper.GetModuleAddress(types.ModuleName))
			s.Require().Equal(initialTxFeesBalance, actualTxFeesBalance)

			// Check beneficiary balance
			var actualBeneficiaryBalance sdk.Coins
			if tc.beneficiary != nil {
				actualBeneficiaryBalance = s.App.BankKeeper.GetAllBalances(s.Ctx, *tc.beneficiary)
			}
			s.Require().Equal(tc.expBeneficiaryRev, actualBeneficiaryBalance)
		})
	}
}

func (s *KeeperTestSuite) ExtractChargeFeeEvent(events []sdk.Event, eventName string) types.EventChargeFee {
	event, found := s.FindLastEventOfType(events, eventName)
	s.Require().True(found)
	chargeFee := types.EventChargeFee{}
	attrs := s.ExtractAttributes(event)
	for key, value := range attrs {
		switch key {
		case "payer":
			chargeFee.Payer = value
		case "taker_fee":
			chargeFee.TakerFee = value
		case "beneficiary":
			chargeFee.Beneficiary = value
		case "beneficiary_revenue":
			chargeFee.BeneficiaryRevenue = value
		}
	}
	return chargeFee
}
