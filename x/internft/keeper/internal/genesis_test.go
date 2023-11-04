package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
	keeper "github.com/0tech/andromeda/x/internft/keeper/internal"
	"github.com/0tech/andromeda/x/internft/module"
)

func TestImportExportGenesis(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(internftv1alpha1.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	keeper := keeper.NewKeeper(encCfg.Codec, storeService, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})

	addr := createAddresses(1, "addr")[0]
	classIDs := createIDs(2, "class")
	nftIDs := createIDs(2, "nft42")

	testCases := map[string]struct {
		gs *internftv1alpha1.GenesisState
	}{
		"default": {
			gs: internftv1alpha1.DefaultGenesisState(),
		},
		"all features": {
			gs: &internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []internftv1alpha1.Trait{
							{
								Id: "color",
							},
							{
								Id:      "level",
								Variable: true,
							},
						},
						Nfts: []internftv1alpha1.GenesisNFT{
							{
								Id: nftIDs[0],
								Properties: []internftv1alpha1.Property{
									{
										Id:   "color",
										Fact: "white",
									},
									{
										Id:   "level",
										Fact: "42",
									},
								},
								Owner: addr.String(),
							},
							{
								Id:    nftIDs[1],
								Owner: addr.String(),
							},
						},
					},
					{
						Id:              classIDs[1],
						Nfts: []internftv1alpha1.GenesisNFT{
							{
								Id:    nftIDs[0],
								Owner: addr.String(),
							},
							{
								Id:    nftIDs[1],
								Owner: addr.String(),
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, _ := ctx.CacheContext()

			err := tc.gs.ValidateBasic()
			assert.NoError(t, err)

			err = keeper.InitGenesis(ctx, tc.gs)
			require.NoError(t, err)

			exported := keeper.ExportGenesis(ctx)
			require.Equal(t, tc.gs, exported)
		})
	}
}
