package internal

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) SubmitProposal(ctx context.Context, proposer, agent sdk.AccAddress, preActions, postActions []*codectypes.Any, metadata string) error {
	if err := k.validateMetadata(ctx, metadata); err != nil {
		return err
	}

	agentInfo, err := k.GetAgent(ctx, agent)
	if err != nil {
		return err
	}

	if !proposer.Equals(sdk.AccAddress(agentInfo.Creator)) {
		return escrowv1alpha1.ErrPermissionDenied.Wrap("proposer differs from creator")
	}

	if err := k.setProposal(ctx, agent, &escrowv1alpha1.Proposal{
		Proposer:    proposer,
		PreActions:  preActions,
		PostActions: postActions,
		Metadata:    metadata,
	}); err != nil {
		return err
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
			return errors.Wrap(err, phase.name)
		}
	}

	return k.removeAgent(ctx, agent)
}

func (k Keeper) validateActions(actions []*codectypes.Any, signers []sdk.AccAddress) error {
	signerMap := map[string]bool{}
	for _, signer := range signers {
		signerMap[string(signer)] = true
	}

	for i, action := range actions {
		addIndex := func(err error) error {
			return errors.Wrapf(err, "index %d", i)
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
				return addIndex(errors.Wrap(escrowv1alpha1.ErrPermissionDenied.Wrap("wrong signer"), signerStr))
			}
		}
	}

	return nil
}

func (k Keeper) GetProposal(ctx context.Context, agent sdk.AccAddress) (*escrowv1alpha1.Proposal, error) {
	proposal, err := k.proposals.Get(ctx, agent)
	if err != nil {
		if !errors.IsOf(err, collections.ErrNotFound) {
			return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		return nil, escrowv1alpha1.ErrProposalNotFound.Wrap(err.Error())
	}
	k.fixActions(&proposal)

	return &proposal, nil
}

func (k Keeper) setProposal(ctx context.Context, agent sdk.AccAddress, proposal *escrowv1alpha1.Proposal) error {
	if err := k.proposals.Set(ctx, agent, *proposal); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) removeProposal(ctx context.Context, agent sdk.AccAddress) error {
	if err := k.proposals.Remove(ctx, agent); err != nil {
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

func (k Keeper) iterateProposals(ctx context.Context, fn func(agent sdk.AccAddress, proposal escrowv1alpha1.Proposal) error) error {
	iter, err := k.proposals.Iterate(ctx, nil)
	if err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		value, err := iter.Value()
		if err != nil {
			return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
		}

		agent := key
		proposal := value
		k.fixActions(&proposal)

		if err := fn(agent, proposal); err != nil {
			return err
		}
	}

	return nil
}
