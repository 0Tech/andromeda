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
	agentsKeyCreatorPrefix = collections.NewPrefix(0x12)

	proposalsKeyPrefix         = collections.NewPrefix(0x20)
	proposalsKeyProposerPrefix = collections.NewPrefix(0x21)
)

func newAgentsIndexes(sb *collections.SchemaBuilder) agentsIndexes {
	return agentsIndexes{
		creator: indexes.NewMulti(
			sb, agentsKeyCreatorPrefix, "agents_by_creator",
			sdk.AccAddressKey,
			sdk.AccAddressKey,
			func(_ sdk.AccAddress, value escrowv1alpha1.Agent) (sdk.AccAddress, error) {
				creator := value.Creator
				return creator, nil
			},
		),
	}
}

type agentsIndexes struct {
	creator *indexes.Multi[sdk.AccAddress, sdk.AccAddress, escrowv1alpha1.Agent]
}

func (a agentsIndexes) IndexesList() []collections.Index[sdk.AccAddress, escrowv1alpha1.Agent] {
	return []collections.Index[sdk.AccAddress, escrowv1alpha1.Agent]{
		a.creator,
	}
}

func newProposalsIndexes(sb *collections.SchemaBuilder) proposalsIndexes {
	return proposalsIndexes{
		proposer: indexes.NewMulti(
			sb, proposalsKeyProposerPrefix, "proposals_by_proposer",
			sdk.AccAddressKey,
			sdk.AccAddressKey,
			func(_ sdk.AccAddress, value escrowv1alpha1.Proposal) (sdk.AccAddress, error) {
				proposer := value.Proposer
				return proposer, nil
			},
		),
	}
}

type proposalsIndexes struct {
	proposer *indexes.Multi[sdk.AccAddress, sdk.AccAddress, escrowv1alpha1.Proposal]
}

func (a proposalsIndexes) IndexesList() []collections.Index[sdk.AccAddress, escrowv1alpha1.Proposal] {
	return []collections.Index[sdk.AccAddress, escrowv1alpha1.Proposal]{
		a.proposer,
	}
}
