package testutils

import (
	"encoding/json"
	"os"
	"time"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simapp "github.com/cosmos/cosmos-sdk/testutil/sims"

	osmod "github.com/osmosis-labs/osmosis/v15/app"
	"github.com/osmosis-labs/osmosis/v15/app/params"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// SimApp testing.
var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

var TestChainID = "dymension_100-1"

var defaultGenesisBz []byte

func getDefaultGenesisState(encCdc params.EncodingConfig) []byte {
	if len(defaultGenesisBz) == 0 {
		genesisState := osmod.NewDefaultGenesisState(encCdc.Codec)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}
		defaultGenesisBz = stateBytes
	}
	return defaultGenesisBz
}

// Setup initializes a new OsmosisApp.
func Setup(isCheckTx bool, chainID string) *osmod.App {
	db := dbm.NewMemDB()
	encCdc := osmod.MakeEncodingConfig()
	if chainID == "" {
		chainID = TestChainID
	}
	app := osmod.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, osmod.DefaultNodeHome, 0, encCdc, simapp.EmptyAppOptions{}, baseapp.SetChainID(chainID))

	if !isCheckTx {
		genesisState := osmod.NewDefaultGenesisState(encCdc.Codec)
		// set EnableCreate to false
		if evmGenesisStateJson, found := genesisState[evmtypes.ModuleName]; found {
			// force disable Enable Create of x/evm
			var evmGenesisState evmtypes.GenesisState
			encCdc.Codec.MustUnmarshalJSON(evmGenesisStateJson, &evmGenesisState)
			evmGenesisState.Params.EnableCreate = false
			genesisState[evmtypes.ModuleName] = encCdc.Codec.MustMarshalJSON(&evmGenesisState)
		}

		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		_ = app.InitChain(
			abci.RequestInitChain{
				ChainId:         chainID,
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// SetupTestingAppWithLevelDb initializes a new OsmosisApp intended for testing,
// with LevelDB as a db.
func SetupTestingAppWithLevelDb(isCheckTx bool) (app *osmod.App, cleanupFn func()) {
	dir, err := os.MkdirTemp(os.TempDir(), "osmosis_leveldb_testing")
	if err != nil {
		panic(err)
	}
	db, err := dbm.NewDB("osmosis_leveldb_testing", dbm.GoLevelDBBackend, dir)
	if err != nil {
		panic(err)
	}
	encCdc := osmod.MakeEncodingConfig()
	app = osmod.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, osmod.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{}, baseapp.SetChainID(TestChainID))

	if !isCheckTx {
		genesisState := osmod.NewDefaultGenesisState(encCdc.Codec)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				ChainId:         TestChainID,
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	cleanupFn = func() {
		db.Close()
		err = os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}

	return app, cleanupFn
}
