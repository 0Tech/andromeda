package keeper

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	internft "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
	"github.com/0tech/andromeda/x/internft/keeper/internal"
)

type Keeper struct {
	impl internal.Keeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	authority string,
) Keeper {
	return Keeper{
		impl: internal.NewKeeper(
			storeKey,
			cdc,
			authority,
		),
	}
}

func NewMsgServer(keeper Keeper) internft.MsgServer {
	return internal.NewMsgServer(keeper.impl)
}

func NewQueryServer(keeper Keeper) internft.QueryServer {
	return internal.NewQueryServer(keeper.impl)
}

func (k Keeper) InitGenesis(ctx sdk.Context, gs *internft.GenesisState) error {
	return k.impl.InitGenesis(ctx, gs)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *internft.GenesisState {
	return k.impl.ExportGenesis(ctx)
}
