package keeper

import (
	"context"

	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
	"github.com/0tech/andromeda/x/internft/keeper/internal"
)

type Keeper struct {
	impl internal.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
) Keeper {
	return Keeper{
		impl: internal.NewKeeper(
			cdc,
			storeService,
			authority,
		),
	}
}

func NewMsgServer(keeper Keeper) internftv1alpha1.MsgServer {
	return internal.NewMsgServer(keeper.impl)
}

func NewQueryServer(keeper Keeper) internftv1alpha1.QueryServer {
	return internal.NewQueryServer(keeper.impl)
}

func (k Keeper) InitGenesis(ctx context.Context, gs *internftv1alpha1.GenesisState) error {
	return k.impl.InitGenesis(ctx, gs)
}

func (k Keeper) ExportGenesis(ctx context.Context) *internftv1alpha1.GenesisState {
	return k.impl.ExportGenesis(ctx)
}
