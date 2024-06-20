package app

import (
	"encoding/json"
	"os"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	simapp "github.com/cosmos/cosmos-sdk/testutil/sims"
	dymd "github.com/dymensionxyz/dymension/v3/app"
)

var defaultGenesisBz []byte

func getDefaultGenesisStateBytes(cdc codec.JSONCodec) []byte {
	if len(defaultGenesisBz) == 0 {
		genesisState := dymd.NewDefaultGenesisState(cdc)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}
		defaultGenesisBz = stateBytes
	}
	return defaultGenesisBz
}

// Setup initializes a new OsmosisApp.
func Setup(isCheckTx bool) *dymd.App {
	db := dbm.NewMemDB()
	encCdc := dymd.MakeEncodingConfig()
	app := dymd.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, dymd.DefaultNodeHome, 0, encCdc, simapp.EmptyAppOptions{})

	if !isCheckTx {
		stateBytes := getDefaultGenesisStateBytes(encCdc.Codec)

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// SetupTestingAppWithLevelDb initializes a new OsmosisApp intended for testing,
// with LevelDB as a db.
func SetupTestingAppWithLevelDb(isCheckTx bool) (app *dymd.App, cleanupFn func()) {
	dir, err := os.MkdirTemp(os.TempDir(), "osmosis_leveldb_testing")
	if err != nil {
		panic(err)
	}
	db, err := dbm.NewDB("osmosis_leveldb_testing", dbm.GoLevelDBBackend, dir)
	if err != nil {
		panic(err)
	}
	encCdc := dymd.MakeEncodingConfig()
	app = dymd.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, dymd.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})

	if !isCheckTx {
		genesisState := dymd.NewDefaultGenesisState(encCdc.Codec)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
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
