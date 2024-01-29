package internal

import (
	"context"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) DefaultGenesis() *escrowv1alpha1.GenesisState {
	return &escrowv1alpha1.GenesisState{
		Params:       DefaultGenesisParams(),
		NextAgent:    1,
		Agents:       []*escrowv1alpha1.GenesisState_Agent{},
		NextProposal: 1,
		Proposals:    []*escrowv1alpha1.GenesisState_Proposal{},
	}
}

func DefaultGenesisParams() *escrowv1alpha1.GenesisState_Params {
	return &escrowv1alpha1.GenesisState_Params{
		MaxMetadataLength: 255,
	}
}

func (k Keeper) ValidateGenesis(gs *escrowv1alpha1.GenesisState) error {
	if gs.Params == nil {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil params")
	}

	if gs.NextAgent == 0 {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil next_agent")
	}

	if gs.Agents == nil {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil agents")
	}

	if gs.NextProposal == 0 {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil next_proposal")
	}

	if gs.Proposals == nil {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil proposals")
	}

	if err := k.validateGenesisParams(gs.Params); err != nil {
		return errors.Wrap(err, "params")
	}

	if err := k.validateGenesisAgents(gs.Agents); err != nil {
		return errors.Wrap(err, "agents")
	}

	if err := k.validateGenesisProposals(gs.Proposals); err != nil {
		return errors.Wrap(err, "proposals")
	}

	return nil
}

func (k Keeper) validateGenesisParams(params *escrowv1alpha1.GenesisState_Params) error {
	if params.MaxMetadataLength == 0 {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil max_metadata_length")
	}

	return nil
}

func (k Keeper) validateGenesisAgents(agents []*escrowv1alpha1.GenesisState_Agent) error {
	seen := map[string]bool{}
	for i, agent := range agents {
		if agent == nil {
			return indexedError(escrowv1alpha1.ErrUnimplemented.Wrap("nil agent"), i)
		}

		if err := k.validateGenesisAgent(agent); err != nil {
			return indexedError(err, i)
		}

		if seen[agent.Address] {
			return indexedError(escrowv1alpha1.ErrDuplicateEntry, i)
		}
		seen[agent.Address] = true
	}

	return nil
}

func (k Keeper) validateGenesisAgent(agent *escrowv1alpha1.GenesisState_Agent) error {
	if agent.Address == "" {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil address")
	}

	if agent.Creator == "" {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil creator")
	}

	if _, err := k.addressStringToBytes(agent.Address); err != nil {
		return errors.Wrap(err, "address")
	}

	if _, err := k.addressStringToBytes(agent.Creator); err != nil {
		return errors.Wrap(err, "creator")
	}

	return nil
}

func (k Keeper) validateGenesisProposals(proposals []*escrowv1alpha1.GenesisState_Proposal) error {
	seen := map[uint64]bool{}
	for i, proposal := range proposals {
		if proposal == nil {
			return indexedError(escrowv1alpha1.ErrUnimplemented.Wrap("nil proposal"), i)
		}

		if err := k.validateGenesisProposal(proposal); err != nil {
			return indexedError(err, i)
		}

		if seen[proposal.Id] {
			return indexedError(escrowv1alpha1.ErrDuplicateEntry, i)
		}
		seen[proposal.Id] = true
	}

	return nil
}

func (k Keeper) validateGenesisProposal(proposal *escrowv1alpha1.GenesisState_Proposal) error {
	if proposal.Id == 0 {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil id")
	}

	if proposal.Proposer == "" {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil proposer")
	}

	if proposal.Agent == "" {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil agent")
	}

	if proposal.PreActions == nil {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil pre_actions")
	}

	if proposal.PostActions == nil {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil post_actions")
	}

	if proposal.Metadata == "" {
		return escrowv1alpha1.ErrUnimplemented.Wrap("nil metadata")
	}

	if _, err := k.addressStringToBytes(proposal.Proposer); err != nil {
		return errors.Wrap(err, "proposer")
	}

	if _, err := k.addressStringToBytes(proposal.Agent); err != nil {
		return errors.Wrap(err, "agent")
	}

	return nil
}

func (k Keeper) InitGenesis(ctx context.Context, gs *escrowv1alpha1.GenesisState) error {
	if err := k.initGenesisParams(ctx, gs.Params); err != nil {
		return errors.Wrap(err, "params")
	}

	if err := k.nextAgent.Set(ctx, gs.NextAgent); err != nil {
		return errors.Wrap(err, "next_agent")
	}

	if err := k.initGenesisAgents(ctx, gs.Agents); err != nil {
		return errors.Wrap(err, "agents")
	}

	if err := k.nextProposal.Set(ctx, gs.NextProposal); err != nil {
		return errors.Wrap(err, "next_proposal")
	}

	if err := k.initGenesisProposals(ctx, gs.Proposals); err != nil {
		return errors.Wrap(err, "proposals")
	}

	return nil
}

func (k Keeper) initGenesisParams(ctx context.Context, params *escrowv1alpha1.GenesisState_Params) error {
	return k.setParams(ctx, &escrowv1alpha1.Params{
		MaxMetadataLength: params.MaxMetadataLength,
	})
}

func (k Keeper) initGenesisAgents(ctx context.Context, agents []*escrowv1alpha1.GenesisState_Agent) error {
	for i, agent := range agents {
		if err := k.initGenesisAgent(ctx, agent); err != nil {
			return indexedError(err, i)
		}
	}

	return nil
}

func (k Keeper) initGenesisAgent(ctx context.Context, agent *escrowv1alpha1.GenesisState_Agent) error {
	address, err := k.addressStringToBytes(agent.Address)
	if err != nil {
		return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "address")
	}

	creator, err := k.addressStringToBytes(agent.Creator)
	if err != nil {
		return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "creator")
	}

	return k.setAgent(ctx, address, &escrowv1alpha1.Agent{
		Creator: creator,
	})
}

func (k Keeper) initGenesisProposals(ctx context.Context, proposals []*escrowv1alpha1.GenesisState_Proposal) error {
	for i, proposal := range proposals {
		if err := k.initGenesisProposal(ctx, proposal); err != nil {
			return indexedError(err, i)
		}
	}

	return nil
}

func (k Keeper) initGenesisProposal(ctx context.Context, proposal *escrowv1alpha1.GenesisState_Proposal) error {
	proposer, err := k.addressStringToBytes(proposal.Proposer)
	if err != nil {
		return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "proposer")
	}

	agent, err := k.addressStringToBytes(proposal.Agent)
	if err != nil {
		return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "agent")
	}

	return k.setProposal(ctx, proposal.Id, &escrowv1alpha1.Proposal{
		Proposer:    proposer,
		Agent:       agent,
		PreActions:  proposal.PreActions,
		PostActions: proposal.PostActions,
		Metadata:    proposal.Metadata,
	})
}

func (k Keeper) ExportGenesis(ctx context.Context) (*escrowv1alpha1.GenesisState, error) {
	params, err := k.exportGenesisParams(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "params")
	}

	nextAgent, err := k.nextAgent.Peek(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "next_agent")
	}

	agents, err := k.exportGenesisAgents(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "agents")
	}

	nextProposal, err := k.nextProposal.Peek(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "next_proposal")
	}

	proposals, err := k.exportGenesisProposals(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "proposals")
	}

	return &escrowv1alpha1.GenesisState{
		Params:       params,
		NextAgent:    nextAgent,
		Agents:       agents,
		NextProposal: nextProposal,
		Proposals:    proposals,
	}, nil
}

func (k Keeper) exportGenesisParams(ctx context.Context) (*escrowv1alpha1.GenesisState_Params, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &escrowv1alpha1.GenesisState_Params{
		MaxMetadataLength: params.MaxMetadataLength,
	}, nil
}

func (k Keeper) exportGenesisAgents(ctx context.Context) ([]*escrowv1alpha1.GenesisState_Agent, error) {
	agents := []*escrowv1alpha1.GenesisState_Agent{}
	if err := k.iterateAgents(ctx, func(address sdk.AccAddress, agent escrowv1alpha1.Agent) error {
		addressStr, err := k.addressBytesToString(address)
		if err != nil {
			return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "address")
		}

		creatorStr, err := k.addressBytesToString(agent.Creator)
		if err != nil {
			return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "creator")
		}

		agents = append(agents, &escrowv1alpha1.GenesisState_Agent{
			Address: addressStr,
			Creator: creatorStr,
		})

		return nil
	}); err != nil {
		return nil, err
	}

	return agents, nil
}

func (k Keeper) exportGenesisProposals(ctx context.Context) ([]*escrowv1alpha1.GenesisState_Proposal, error) {
	proposals := []*escrowv1alpha1.GenesisState_Proposal{}
	if err := k.iterateProposals(ctx, func(id uint64, proposal escrowv1alpha1.Proposal) error {
		proposerStr, err := k.addressBytesToString(proposal.Proposer)
		if err != nil {
			return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "proposer")
		}

		agentStr, err := k.addressBytesToString(proposal.Agent)
		if err != nil {
			return errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "agent")
		}

		proposals = append(proposals, &escrowv1alpha1.GenesisState_Proposal{
			Id:          id,
			Proposer:    proposerStr,
			Agent:       agentStr,
			PreActions:  proposal.PreActions,
			PostActions: proposal.PostActions,
			Metadata:    proposal.Metadata,
		})

		return nil
	}); err != nil {
		return nil, err
	}

	return proposals, nil
}
