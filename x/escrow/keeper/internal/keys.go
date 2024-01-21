package internal

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

var (
	paramsKey = collections.NewPrefix(0x00)

	agentsSeqKey           = collections.NewPrefix(0x10)
	agentsKeyPrefix        = collections.NewPrefix(0x11)
	agentsKeyAddressPrefix = collections.NewPrefix(0x12)

	proposalsSeqKey      = collections.NewPrefix(0x20)
	proposalsKeyPrefix   = collections.NewPrefix(0x21)
	proposalsKeyIDPrefix = collections.NewPrefix(0x22)
)

func newAgentsIndexes(sb *collections.SchemaBuilder) agentsIndexes {
	return agentsIndexes{
		address: indexes.NewUnique(
			sb, agentsKeyAddressPrefix, "agents_by_address",
			sdk.AccAddressKey,
			collections.PairKeyCodec(sdk.AccAddressKey, sdk.AccAddressKey),
			func(key collections.Pair[sdk.AccAddress, sdk.AccAddress], _ escrowv1alpha1.Agent) (sdk.AccAddress, error) {
				address := key.K2()
				return address, nil
			},
		),
	}
}

type agentsIndexes struct {
	address *indexes.Unique[sdk.AccAddress, collections.Pair[sdk.AccAddress, sdk.AccAddress], escrowv1alpha1.Agent]
}

func (a agentsIndexes) IndexesList() []collections.Index[collections.Pair[sdk.AccAddress, sdk.AccAddress], escrowv1alpha1.Agent] {
	return []collections.Index[collections.Pair[sdk.AccAddress, sdk.AccAddress], escrowv1alpha1.Agent]{
		a.address,
	}
}

func newProposalsIndexes(sb *collections.SchemaBuilder) proposalsIndexes {
	return proposalsIndexes{
		id: indexes.NewUnique(
			sb, proposalsKeyIDPrefix, "proposals_by_id",
			collections.Uint64Key,
			collections.PairKeyCodec(sdk.AccAddressKey, collections.Uint64Key),
			func(key collections.Pair[sdk.AccAddress, uint64], _ escrowv1alpha1.Proposal) (uint64, error) {
				id := key.K2()
				return id, nil
			},
		),
	}
}

type proposalsIndexes struct {
	id *indexes.Unique[uint64, collections.Pair[sdk.AccAddress, uint64], escrowv1alpha1.Proposal]
}

func (a proposalsIndexes) IndexesList() []collections.Index[collections.Pair[sdk.AccAddress, uint64], escrowv1alpha1.Proposal] {
	return []collections.Index[collections.Pair[sdk.AccAddress, uint64], escrowv1alpha1.Proposal]{
		a.id,
	}
}
