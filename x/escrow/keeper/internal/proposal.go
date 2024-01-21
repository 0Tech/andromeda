package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) SubmitProposal(ctx context.Context, proposer, agent sdk.AccAddress, preActions, postActions []*codectypes.Any) (uint64, error) {
	if err := k.HasAgent(ctx, agent, proposer); err != nil {
		return 0, err
	}

	id, err := k.nextProposal.Next(ctx)
	if err != nil {
		return 0, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	if err := k.setProposal(ctx, id, proposer, &escrowv1alpha1.Proposal{
		Agent:       agent,
		PreActions:  preActions,
		PostActions: postActions,
	}); err != nil {
		return 0, err
	}

	for _, phase := range []struct {
		name    string
		actions []*codectypes.Any
	}{
		{
			name:    "pre_actions",
			actions: preActions,
		},
	} {
		if err := k.executeActions(ctx, phase.actions); err != nil {
			return 0, errorsmod.Wrap(err, phase.name)
		}
	}

	if err := k.removeAgent(ctx, agent); err != nil {
		return 0, err
	}

	return id, nil
}

func (k Keeper) validateActions(actions []*codectypes.Any, signers []sdk.AccAddress) error {
	signerMap := map[string]bool{}
	for _, signer := range signers {
		signerMap[string(signer)] = true
	}

	for i, action := range actions {
		addIndex := func(err error) error {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		msg, err := k.anyToMsg(*action)
		if err != nil {
			return addIndex(err)
		}

		actionSigners, _, err := k.cdc.GetMsgV1Signers(msg)
		if err != nil {
			return addIndex(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()))
		}

		for _, signer := range actionSigners {
			if !signerMap[string(signer)] {
				signerStr, err := k.addressBytesToString(signer)
				if err != nil {
					return addIndex(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()))
				}
				return addIndex(errorsmod.Wrap(escrowv1alpha1.ErrPermissionDenied.Wrap("wrong signer"), signerStr))
			}
		}
	}

	return nil
}

func (k Keeper) getProposalKey(ctx context.Context, id uint64) (*collections.Pair[sdk.AccAddress, uint64], error) {
	key, err := k.proposals.Indexes.id.MatchExact(ctx, id)
	if err != nil {
		if !errorsmod.IsOf(err, collections.ErrNotFound) {
			return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		return nil, escrowv1alpha1.ErrProposalNotFound.Wrap(err.Error())
	}

	return &key, nil
}

func (k Keeper) GetProposal(ctx context.Context, id uint64) (sdk.AccAddress, *escrowv1alpha1.Proposal, error) {
	key, err := k.getProposalKey(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	proposer := key.K1()

	proposal, err := k.proposals.Get(ctx, *key)
	if err != nil {
		return nil, nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	k.fixActions(&proposal)

	return proposer, &proposal, nil
}

func (k Keeper) setProposal(ctx context.Context, id uint64, proposer sdk.AccAddress, proposal *escrowv1alpha1.Proposal) error {
	if err := k.proposals.Set(ctx, collections.Join(proposer, id), *proposal); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) removeProposal(ctx context.Context, id uint64) error {
	key, err := k.getProposalKey(ctx, id)
	if err != nil {
		return err
	}

	if err := k.proposals.Remove(ctx, *key); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) fixActions(proposal *escrowv1alpha1.Proposal) {
	// TODO: find better solution
	if proposal.PreActions == nil {
		proposal.PreActions = []*codectypes.Any{}
	}
	if proposal.PostActions == nil {
		proposal.PostActions = []*codectypes.Any{}
	}
}

func (k Keeper) iterateProposals(ctx context.Context, fn func(id uint64, proposer sdk.AccAddress, proposal escrowv1alpha1.Proposal) error) error {
	iter, err := k.proposals.Indexes.id.Iterate(ctx, nil)
	if err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, err := iter.PrimaryKey()
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		value, err := k.proposals.Get(ctx, key)
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		proposer := key.K1()
		id := key.K2()
		proposal := value

		if err := fn(id, proposer, proposal); err != nil {
			return err
		}
	}

	return nil
}
