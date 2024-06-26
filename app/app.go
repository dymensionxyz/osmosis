package app

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"

	appparams "github.com/osmosis-labs/osmosis/v15/app/params"

	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"

	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/keeper"

	/* ------------------------------ ethermint imports ----------------------------- */

	ethermint "github.com/evmos/ethermint/types"
	"github.com/evmos/ethermint/x/evm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/feemarket"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"

	/* ----------------------------- osmosis imports ---------------------------- */

	"github.com/osmosis-labs/osmosis/v15/x/epochs"
	epochskeeper "github.com/osmosis-labs/osmosis/v15/x/epochs/keeper"
	"github.com/osmosis-labs/osmosis/v15/x/lockup"
	lockupkeeper "github.com/osmosis-labs/osmosis/v15/x/lockup/keeper"

	"github.com/osmosis-labs/osmosis/v15/x/gamm"
	gammkeeper "github.com/osmosis-labs/osmosis/v15/x/gamm/keeper"
	"github.com/osmosis-labs/osmosis/v15/x/incentives"
	incentiveskeeper "github.com/osmosis-labs/osmosis/v15/x/incentives/keeper"
	"github.com/osmosis-labs/osmosis/v15/x/poolmanager"
	poolmanagerkeeper "github.com/osmosis-labs/osmosis/v15/x/poolmanager/keeper"

	"github.com/osmosis-labs/osmosis/v15/x/txfees"
	txfeeskeeper "github.com/osmosis-labs/osmosis/v15/x/txfees/keeper"
	/* ----------------------------- test imports ---------------------------- */)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+appparams.Name)

	sdk.DefaultPowerReduction = ethermint.PowerReduction
}

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper                 authkeeper.AccountKeeper
	AuthzKeeper                   authzkeeper.Keeper
	BankKeeper                    bankkeeper.Keeper
	CapabilityKeeper              *capabilitykeeper.Keeper
	StakingKeeper                 stakingkeeper.Keeper
	SlashingKeeper                slashingkeeper.Keeper
	MintKeeper                    mintkeeper.Keeper
	DistrKeeper                   distrkeeper.Keeper
	GovKeeper                     govkeeper.Keeper
	CrisisKeeper                  *crisiskeeper.Keeper
	UpgradeKeeper                 *upgradekeeper.Keeper
	ParamsKeeper                  paramskeeper.Keeper
	IBCKeeper                     *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper                evidencekeeper.Keeper
	TransferKeeper                ibctransferkeeper.Keeper
	FeeGrantKeeper                feegrantkeeper.Keeper
	PacketForwardMiddlewareKeeper *packetforwardkeeper.Keeper
	ConsensusParamsKeeper         consensusparamkeeper.Keeper

	// Ethermint keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// Osmosis keepers
	GAMMKeeper        *gammkeeper.Keeper
	PoolManagerKeeper *poolmanagerkeeper.Keeper
	LockupKeeper      *lockupkeeper.Keeper
	EpochsKeeper      *epochskeeper.Keeper
	IncentivesKeeper  *incentiveskeeper.Keeper
	TxFeesKeeper      *txfeeskeeper.Keeper

	// the module manager
	mm *module.Manager

	// module configurator
	configurator module.Configurator
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		consensus.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(nil),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic

		// Ethermint modules
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},

		// Osmosis modules
		lockup.AppModuleBasic{},
		epochs.AppModuleBasic{},
		gamm.AppModuleBasic{},
		poolmanager.AppModuleBasic{},
		incentives.AppModuleBasic{},
		txfees.AppModuleBasic{},
	)
)

func New(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool, homePath string, invCheckPeriod uint,
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	return nil
}

func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// GetBaseApp returns the base app of the application
func (app App) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *App) ExportState(ctx sdk.Context) map[string]json.RawMessage {
	return app.mm.ExportGenesis(ctx, app.AppCodec())
}

// GetTxConfig implements ibctesting.TestingApp
func (app *App) GetTxConfig() client.TxConfig {
	return moduletestutil.MakeTestEncodingConfig().TxConfig
}

func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	defaultGenesisState := ModuleBasics.DefaultGenesis(cdc)
	return defaultGenesisState
}
