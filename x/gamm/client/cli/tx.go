package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	math "cosmossdk.io/math"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/osmosis-labs/osmosis/v15/osmoutils/osmocli"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/pool-models/balancer"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewTxCmd() *cobra.Command {
	txCmd := osmocli.TxIndexCmd(types.ModuleName)
	osmocli.AddTxCmd(txCmd, NewJoinPoolCmd)
	osmocli.AddTxCmd(txCmd, NewExitPoolCmd)
	osmocli.AddTxCmd(txCmd, NewSwapExactAmountInCmd)
	osmocli.AddTxCmd(txCmd, NewSwapExactAmountOutCmd)
	osmocli.AddTxCmd(txCmd, NewJoinSwapExternAmountIn)
	osmocli.AddTxCmd(txCmd, NewJoinSwapShareAmountOut)
	osmocli.AddTxCmd(txCmd, NewExitSwapExternAmountOut)
	osmocli.AddTxCmd(txCmd, NewExitSwapShareAmountIn)
	txCmd.AddCommand(
		NewCreatePoolCmd().BuildCommandCustomFn(),
	)
	return txCmd
}

var poolIdFlagOverride = map[string]string{
	"poolid": FlagPoolId,
}

func NewCreatePoolCmd() *osmocli.TxCliDesc {
	desc := osmocli.TxCliDesc{
		Use:   "create-pool [flags]",
		Short: "create a new pool and provide the liquidity to it",
		Long: `Must provide path to a pool JSON file (--pool-file) describing the pool to be created
Sample pool JSON file contents for balancer:
{
	"weights": "4uatom,4osmo,2uakt",
	"initial-deposit": "100uatom,5osmo,20uakt",
	"swap-fee": "0.01",
	"exit-fee": "0.01",
	"future-governor": "168h"
}
`,
		NumArgs:          0,
		ParseAndBuildMsg: BuildCreatePoolCmd,
		Flags: osmocli.FlagDesc{
			RequiredFlags: []*flag.FlagSet{FlagSetCreatePoolFile()},
			OptionalFlags: []*flag.FlagSet{FlagSetCreatePoolType()},
		},
	}
	return &desc
}

func NewJoinPoolCmd() (*osmocli.TxCliDesc, *types.MsgJoinPool) {
	return &osmocli.TxCliDesc{
		Use:   "join-pool",
		Short: "join a new pool and provide the liquidity to it",
		CustomFlagOverrides: map[string]string{
			"poolid": FlagPoolId,
		},
		CustomFieldParsers: map[string]osmocli.CustomFieldParserFn{
			"TokenInMaxs":    osmocli.FlagOnlyParser(maxAmountsInParser),
			"ShareOutAmount": osmocli.FlagOnlyParser(shareAmountOutParser),
		},
		Flags: osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetJoinPool()}},
	}, &types.MsgJoinPool{}
}

func NewExitPoolCmd() (*osmocli.TxCliDesc, *types.MsgExitPool) {
	return &osmocli.TxCliDesc{
		Use:   "exit-pool",
		Short: "exit a new pool and withdraw the liquidity from it",
		CustomFlagOverrides: map[string]string{
			"poolid": FlagPoolId,
		},
		CustomFieldParsers: map[string]osmocli.CustomFieldParserFn{
			"TokenOutMins":  osmocli.FlagOnlyParser(minAmountsOutParser),
			"ShareInAmount": osmocli.FlagOnlyParser(shareAmountInParser),
		},
		Flags: osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetExitPool()}},
	}, &types.MsgExitPool{}
}

func NewSwapExactAmountInCmd() (*osmocli.TxCliDesc, *types.MsgSwapExactAmountIn) {
	return &osmocli.TxCliDesc{
		Use:   "swap-exact-amount-in [token-in] [token-out-min-amount]",
		Short: "swap exact amount in",
		CustomFieldParsers: map[string]osmocli.CustomFieldParserFn{
			"Routes": osmocli.FlagOnlyParser(swapAmountInRoutes),
		},
		Flags: osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetMultihopSwapRoutes()}},
	}, &types.MsgSwapExactAmountIn{}
}

func NewSwapExactAmountOutCmd() (*osmocli.TxCliDesc, *types.MsgSwapExactAmountOut) {
	// Can't get rid of this parser without a break, because the args are out of order.
	return &osmocli.TxCliDesc{
		Use:              "swap-exact-amount-out [token-out] [token-in-max-amount]",
		Short:            "swap exact amount out",
		NumArgs:          2,
		ParseAndBuildMsg: NewBuildSwapExactAmountOutMsg,
		Flags:            osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetMultihopSwapRoutes()}},
	}, &types.MsgSwapExactAmountOut{}
}

func NewJoinSwapExternAmountIn() (*osmocli.TxCliDesc, *types.MsgJoinSwapExternAmountIn) {
	return &osmocli.TxCliDesc{
		Use:                 "join-swap-extern-amount-in [token-in] [share-out-min-amount]",
		Short:               "join swap extern amount in",
		CustomFlagOverrides: poolIdFlagOverride,
		Flags:               osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetJustPoolId()}},
	}, &types.MsgJoinSwapExternAmountIn{}
}

func NewJoinSwapShareAmountOut() (*osmocli.TxCliDesc, *types.MsgJoinSwapShareAmountOut) {
	return &osmocli.TxCliDesc{
		Use:                 "join-swap-share-amount-out [token-in-denom] [share-out-amount] [token-in-max-amount] ",
		Short:               "join swap share amount out",
		CustomFlagOverrides: poolIdFlagOverride,
		Flags:               osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetJustPoolId()}},
	}, &types.MsgJoinSwapShareAmountOut{}
}

func NewExitSwapExternAmountOut() (*osmocli.TxCliDesc, *types.MsgExitSwapExternAmountOut) {
	return &osmocli.TxCliDesc{
		Use:                 "exit-swap-extern-amount-out [token-out] [share-in-max-amount]",
		Short:               "exit swap extern amount out",
		CustomFlagOverrides: poolIdFlagOverride,
		Flags:               osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetJustPoolId()}},
	}, &types.MsgExitSwapExternAmountOut{}
}

func NewExitSwapShareAmountIn() (*osmocli.TxCliDesc, *types.MsgExitSwapShareAmountIn) {
	return &osmocli.TxCliDesc{
		Use:                 "exit-swap-share-amount-in [token-out-denom] [share-in-amount] [token-out-min-amount]",
		Short:               "exit swap share amount in",
		CustomFlagOverrides: poolIdFlagOverride,
		Flags:               osmocli.FlagDesc{RequiredFlags: []*flag.FlagSet{FlagSetJustPoolId()}},
	}, &types.MsgExitSwapShareAmountIn{}
}

func BuildCreatePoolCmd(clientCtx client.Context, args []string, fs *flag.FlagSet) (sdk.Msg, error) {
	poolType, err := fs.GetString(FlagPoolType)
	if err != nil {
		return nil, err
	}
	poolType = strings.ToLower(poolType)

	var msg sdk.Msg
	if poolType == "balancer" || poolType == "uniswap" {
		msg, err = NewBuildCreateBalancerPoolMsg(clientCtx, fs)
		if err != nil {
			return nil, err
		}
	}

	//TODO: validate poolType

	return msg, nil
}

func NewBuildCreateBalancerPoolMsg(clientCtx client.Context, fs *flag.FlagSet) (sdk.Msg, error) {
	pool, err := parseCreateBalancerPoolFlags(fs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool: %w", err)
	}

	deposit, err := sdk.ParseCoinsNormalized(pool.InitialDeposit)
	if err != nil {
		return nil, err
	}

	poolAssetCoins, err := sdk.ParseDecCoins(pool.Weights)
	if err != nil {
		return nil, err
	}

	if len(deposit) != len(poolAssetCoins) {
		return nil, errors.New("deposit tokens and token weights should have same length")
	}

	swapFee, err := math.LegacyNewDecFromStr(pool.SwapFee)
	if err != nil {
		return nil, err
	}

	exitFee, err := math.LegacyNewDecFromStr(pool.ExitFee)
	if err != nil {
		return nil, err
	}

	var poolAssets []balancer.PoolAsset
	for i := 0; i < len(poolAssetCoins); i++ {
		if poolAssetCoins[i].Denom != deposit[i].Denom {
			return nil, errors.New("deposit tokens and token weights should have same denom order")
		}

		poolAssets = append(poolAssets, balancer.PoolAsset{
			Weight: poolAssetCoins[i].Amount.RoundInt(),
			Token:  deposit[i],
		})
	}

	poolParams := &balancer.PoolParams{
		SwapFee: swapFee,
		ExitFee: exitFee,
	}

	msg := &balancer.MsgCreateBalancerPool{
		Sender:             clientCtx.GetFromAddress().String(),
		PoolParams:         poolParams,
		PoolAssets:         poolAssets,
		FuturePoolGovernor: pool.FutureGovernor,
	}

	if (pool.SmoothWeightChangeParams != smoothWeightChangeParamsInputs{}) {
		duration, err := time.ParseDuration(pool.SmoothWeightChangeParams.Duration)
		if err != nil {
			return nil, fmt.Errorf("could not parse duration: %w", err)
		}

		targetPoolAssetCoins, err := sdk.ParseDecCoins(pool.SmoothWeightChangeParams.TargetPoolWeights)
		if err != nil {
			return nil, err
		}

		if len(targetPoolAssetCoins) != len(poolAssetCoins) {
			return nil, errors.New("initial pool weights and target pool weights should have same length")
		}

		var targetPoolAssets []balancer.PoolAsset
		for i := 0; i < len(targetPoolAssetCoins); i++ {
			if targetPoolAssetCoins[i].Denom != poolAssetCoins[i].Denom {
				return nil, errors.New("initial pool weights and target pool weights should have same denom order")
			}

			targetPoolAssets = append(targetPoolAssets, balancer.PoolAsset{
				Weight: targetPoolAssetCoins[i].Amount.RoundInt(),
				Token:  deposit[i],
				// TODO: This doesn't make sense. Should only use denom, not an sdk.Coin
			})
		}

		smoothWeightParams := balancer.SmoothWeightChangeParams{
			Duration:           duration,
			InitialPoolWeights: poolAssets,
			TargetPoolWeights:  targetPoolAssets,
		}

		if pool.SmoothWeightChangeParams.StartTime != "" {
			startTime, err := time.Parse(time.RFC3339, pool.SmoothWeightChangeParams.StartTime)
			if err != nil {
				return nil, fmt.Errorf("could not parse time: %w", err)
			}

			smoothWeightParams.StartTime = startTime
		}

		msg.PoolParams.SmoothWeightChangeParams = &smoothWeightParams
	}

	return msg, nil
}

func shareAmountInParser(fs *flag.FlagSet) (math.Int, error) {
	return sdkIntParser(FlagShareAmountIn, fs)
}

func shareAmountOutParser(fs *flag.FlagSet) (math.Int, error) {
	return sdkIntParser(FlagShareAmountOut, fs)
}

func sdkIntParser(flagName string, fs *flag.FlagSet) (math.Int, error) {
	amountStr, err := fs.GetString(flagName)
	if err != nil {
		return math.ZeroInt(), err
	}

	res, ok := math.NewIntFromString(amountStr)
	if !ok {
		return math.ZeroInt(), errors.New("invalid share amount")
	}
	return res, nil
}

func maxAmountsInParser(fs *flag.FlagSet) (sdk.Coins, error) {
	return stringArrayCoinsParser(FlagMaxAmountsIn, fs)
}

func minAmountsOutParser(fs *flag.FlagSet) (sdk.Coins, error) {
	return stringArrayCoinsParser(FlagMinAmountsOut, fs)
}

func stringArrayCoinsParser(flagName string, fs *flag.FlagSet) (sdk.Coins, error) {
	amountsArr, err := fs.GetStringArray(flagName)
	if err != nil {
		return nil, err
	}

	coins := sdk.Coins{}
	for i := 0; i < len(amountsArr); i++ {
		parsed, err := sdk.ParseCoinsNormalized(amountsArr[i])
		if err != nil {
			return nil, err
		}
		coins = coins.Add(parsed...)
	}
	return coins, nil
}

func swapAmountInRoutes(fs *flag.FlagSet) ([]poolmanagertypes.SwapAmountInRoute, error) {
	swapRoutePoolIds, err := fs.GetString(FlagSwapRoutePoolIds)
	swapRoutePoolIdsArray := strings.Split(swapRoutePoolIds, ",")
	if err != nil {
		return nil, err
	}

	swapRouteDenoms, err := fs.GetString(FlagSwapRouteDenoms)
	swapRouteDenomsArray := strings.Split(swapRouteDenoms, ",")
	if err != nil {
		return nil, err
	}

	if len(swapRoutePoolIdsArray) != len(swapRouteDenomsArray) {
		return nil, errors.New("swap route pool ids and denoms mismatch")
	}

	routes := []poolmanagertypes.SwapAmountInRoute{}
	for index, poolIDStr := range swapRoutePoolIdsArray {
		pID, err := strconv.Atoi(poolIDStr)
		if err != nil {
			return nil, err
		}
		routes = append(routes, poolmanagertypes.SwapAmountInRoute{
			PoolId:        uint64(pID),
			TokenOutDenom: swapRouteDenomsArray[index],
		})
	}
	return routes, nil
}

func swapAmountOutRoutes(fs *flag.FlagSet) ([]poolmanagertypes.SwapAmountOutRoute, error) {
	swapRoutePoolIds, err := fs.GetString(FlagSwapRoutePoolIds)
	swapRoutePoolIdsArray := strings.Split(swapRoutePoolIds, ",")
	if err != nil {
		return nil, err
	}

	swapRouteDenoms, err := fs.GetString(FlagSwapRouteDenoms)
	swapRouteDenomsArray := strings.Split(swapRouteDenoms, ",")
	if err != nil {
		return nil, err
	}

	if len(swapRoutePoolIdsArray) != len(swapRouteDenomsArray) {
		return nil, errors.New("swap route pool ids and denoms mismatch")
	}

	routes := []poolmanagertypes.SwapAmountOutRoute{}
	for index, poolIDStr := range swapRoutePoolIdsArray {
		pID, err := strconv.Atoi(poolIDStr)
		if err != nil {
			return nil, err
		}
		routes = append(routes, poolmanagertypes.SwapAmountOutRoute{
			PoolId:       uint64(pID),
			TokenInDenom: swapRouteDenomsArray[index],
		})
	}
	return routes, nil
}

func NewBuildSwapExactAmountOutMsg(clientCtx client.Context, args []string, fs *flag.FlagSet) (sdk.Msg, error) {
	tokenOutStr, tokenInMaxAmountStr := args[0], args[1]
	routes, err := swapAmountOutRoutes(fs)
	if err != nil {
		return nil, err
	}

	tokenOut, err := sdk.ParseCoinNormalized(tokenOutStr)
	if err != nil {
		return nil, err
	}

	tokenInMaxAmount, ok := math.NewIntFromString(tokenInMaxAmountStr)
	if !ok {
		return nil, errors.New("invalid token in max amount")
	}
	return &types.MsgSwapExactAmountOut{
		Sender:           clientCtx.GetFromAddress().String(),
		Routes:           routes,
		TokenInMaxAmount: tokenInMaxAmount,
		TokenOut:         tokenOut,
	}, nil
}

// ParseCoinsNoSort parses coins from coinsStr but does not sort them.
// Returns error if parsing fails.
func ParseCoinsNoSort(coinsStr string) (sdk.Coins, error) {
	coinStrs := strings.Split(coinsStr, ",")
	decCoins := make(sdk.DecCoins, len(coinStrs))
	for i, coinStr := range coinStrs {
		coin, err := sdk.ParseDecCoin(coinStr)
		if err != nil {
			return sdk.Coins{}, err
		}

		decCoins[i] = coin
	}
	return sdk.NormalizeCoins(decCoins), nil
}
