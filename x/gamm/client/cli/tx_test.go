package cli_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/osmosis-labs/osmosis/v15/x/gamm/client/cli"
)

func TestParseCoinsNoSort(t *testing.T) {
	const (
		a = "aaa"
		b = "bbb"
		c = "ccc"
		d = "ddd"
	)

	var (
		ten = math.NewInt(10)

		coinA = sdk.NewCoin(a, ten)
		coinB = sdk.NewCoin(b, ten)
		coinC = sdk.NewCoin(c, ten)
		coinD = sdk.NewCoin(d, ten)
	)

	tests := map[string]struct {
		coinsStr      string
		expectedCoins sdk.Coins
	}{
		"ascending": {
			coinsStr: "10aaa,10bbb,10ccc,10ddd",
			expectedCoins: sdk.Coins{
				coinA,
				coinB,
				coinC,
				coinD,
			},
		},
		"descending": {
			coinsStr: "10ddd,10ccc,10bbb,10aaa",
			expectedCoins: sdk.Coins{
				coinD,
				coinC,
				coinB,
				coinA,
			},
		},
		"mixed with different values.": {
			coinsStr: "100ddd,20bbb,300aaa,40ccc",
			expectedCoins: sdk.Coins{
				sdk.NewCoin(d, math.NewInt(100)),
				sdk.NewCoin(b, math.NewInt(20)),
				sdk.NewCoin(a, math.NewInt(300)),
				sdk.NewCoin(c, math.NewInt(40)),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			coins, err := cli.ParseCoinsNoSort(tc.coinsStr)

			require.NoError(t, err)
			require.Equal(t, tc.expectedCoins, coins)
		})
	}
}
