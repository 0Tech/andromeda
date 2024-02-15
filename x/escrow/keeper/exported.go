package keeper

import (
	"context"

	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/keeper/expected"
	"github.com/0tech/andromeda/x/escrow/keeper/internal"
)

type (
	MessageRouter = expected.MessageRouter
	AuthKeeper    = expected.AuthKeeper
)

type Keeper struct {
	impl internal.Keeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	authority sdk.AccAddress,
	router expected.MessageRouter,
	authKeeper expected.AuthKeeper,
) (*Keeper, error) {
	impl, err := internal.NewKeeper(
		cdc,
		storeService,
		authority,
		router,
		authKeeper,
	)
	if err != nil {
		return nil, err
	}

	return &Keeper{impl: *impl}, nil
}

func NewMsgServer(keeper Keeper) escrowv1alpha1.MsgServer {
	return internal.NewMsgServer(keeper.impl)
}

func NewQueryServer(keeper Keeper) escrowv1alpha1.QueryServer {
	return internal.NewQueryServer(keeper.impl)
}

func (k Keeper) DefaultGenesis() *escrowv1alpha1.GenesisState {
	return k.impl.DefaultGenesis()
}

func (k Keeper) ValidateGenesis(gs *escrowv1alpha1.GenesisState) error {
	return k.impl.ValidateGenesis(gs)
}

func (k Keeper) InitGenesis(ctx context.Context, gs *escrowv1alpha1.GenesisState) error {
	return k.impl.InitGenesis(ctx, gs)
}

func (k Keeper) ExportGenesis(ctx context.Context) (*escrowv1alpha1.GenesisState, error) {
	return k.impl.ExportGenesis(ctx)
}
