package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
	"github.com/0tech/andromeda/x/test/keeper/internal"
)

type Keeper struct {
	impl internal.Keeper
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	prefix []byte,
) (*Keeper, error) {
	impl, err := internal.NewKeeper(
		cdc,
		storeService,
		prefix,
	)
	if err != nil {
		return nil, err
	}

	return &Keeper{impl: *impl}, nil
}

func NewMsgServer(keeper Keeper) testv1alpha1.MsgServer {
	return internal.NewMsgServer(keeper.impl)
}

func NewQueryServer(keeper Keeper) testv1alpha1.QueryServer {
	return internal.NewQueryServer(keeper.impl)
}
