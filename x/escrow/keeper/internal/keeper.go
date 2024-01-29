package internal

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/keeper/expected"
)

type Keeper struct {
	cdc codec.Codec

	authority sdk.AccAddress
	router    expected.MessageRouter

	authKeeper expected.AuthKeeper

	schema collections.Schema

	params collections.Item[escrowv1alpha1.Params]

	nextAgent collections.Sequence
	agents    *collections.IndexedMap[sdk.AccAddress, escrowv1alpha1.Agent, agentsIndexes]

	nextProposal collections.Sequence
	proposals    *collections.IndexedMap[collections.Pair[sdk.AccAddress, uint64], escrowv1alpha1.Proposal, proposalsIndexes]
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	authority sdk.AccAddress,
	router expected.MessageRouter,
	authKeeper expected.AuthKeeper,
) (*Keeper, error) {
	sb := collections.NewSchemaBuilder(storeService)
	k := &Keeper{
		cdc:        cdc,
		authority:  authority,
		router:     router,
		authKeeper: authKeeper,
		params: collections.NewItem(sb, paramsKey, "params",
			codec.CollValue[escrowv1alpha1.Params](cdc)),
		nextAgent: collections.NewSequence(sb, agentsSeqKey, "next_agent"),
		agents: collections.NewIndexedMap(sb, agentsKeyPrefix, "agents",
			sdk.AccAddressKey,
			codec.CollValue[escrowv1alpha1.Agent](cdc),
			newAgentsIndexes(sb),
		),
		nextProposal: collections.NewSequence(sb, proposalsSeqKey, "next_proposal"),
		proposals: collections.NewIndexedMap(sb, proposalsKeyPrefix, "proposals",
			collections.PairKeyCodec(sdk.AccAddressKey, collections.Uint64Key),
			codec.CollValue[escrowv1alpha1.Proposal](cdc),
			newProposalsIndexes(sb),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		return nil, err
	}
	k.schema = schema

	return k, nil
}

func (k Keeper) addressBytesToString(addr []byte) (string, error) {
	addressCodec := k.cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addrStr, err := addressCodec.BytesToString(addr)
	if err != nil {
		return "", escrowv1alpha1.ErrInvalidAddress.Wrap(err.Error())
	}

	return addrStr, nil
}

func (k Keeper) addressStringToBytes(addr string) (sdk.AccAddress, error) {
	addressCodec := k.cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addrBytes, err := addressCodec.StringToBytes(addr)
	if err != nil {
		return nil, escrowv1alpha1.ErrInvalidAddress.Wrap(err.Error())
	}

	return addrBytes, nil
}

func (k Keeper) GetAuthority() sdk.AccAddress {
	return k.authority
}

func (k Keeper) validateAuthority(candidate sdk.AccAddress) error {
	if !candidate.Equals(k.authority) {
		return escrowv1alpha1.ErrPermissionDenied.Wrap("not authority")
	}

	return nil
}

func (k Keeper) validateMetadata(ctx context.Context, metadata string) error {
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	if length := uint64(len(metadata)); length > params.MaxMetadataLength {
		return errors.Wrapf(escrowv1alpha1.ErrLargeMetadata.Wrapf("over limit of %d", params.MaxMetadataLength), "%d", length)
	}

	return nil
}

func indexedError(err error, index int) error {
	return errors.Wrapf(err, "index %d", index)
}
