package internal

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
		return errorsmod.Wrap(err, "params")
	}

	if err := k.validateGenesisAgents(gs.Agents); err != nil {
		return errorsmod.Wrap(err, "agents")
	}

	if err := k.validateGenesisProposals(gs.Proposals); err != nil {
		return errorsmod.Wrap(err, "proposals")
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
	seenAddress := sdk.AccAddress{}
	for i, agent := range agents {
		addIndex := func(err error) error {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if agent == nil {
			return addIndex(escrowv1alpha1.ErrUnimplemented.Wrap("nil agent"))
		}

		address, err := k.validateGenesisAgent(agent)
		if err != nil {
			return addIndex(err)
		}

		if !(bytes.Compare(address, seenAddress) > 0) {
			return addIndex(sdkerrors.ErrInvalidRequest.Wrap("unsorted agent"))
		}
		seenAddress = address
	}

	return nil
}

func (k Keeper) validateGenesisAgent(agent *escrowv1alpha1.GenesisState_Agent) (sdk.AccAddress, error) {
	if agent.Address == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil address")
	}

	if agent.Creator == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil creator")
	}

	address, err := k.addressStringToBytes(agent.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "address")
	}

	if _, err := k.addressStringToBytes(agent.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "creator")
	}

	return address, nil
}

func (k Keeper) validateGenesisProposals(proposals []*escrowv1alpha1.GenesisState_Proposal) error {
	seenID := uint64(0)
	for i, proposal := range proposals {
		addIndex := func(err error) error {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if proposal == nil {
			return addIndex(escrowv1alpha1.ErrUnimplemented.Wrap("nil proposal"))
		}

		if err := k.validateGenesisProposal(proposal); err != nil {
			return addIndex(err)
		}

		if !(proposal.Id > seenID) {
			return addIndex(sdkerrors.ErrInvalidRequest.Wrap("unsorted proposal"))
		}
		seenID = proposal.Id
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
		return errorsmod.Wrap(err, "proposer")
	}

	if _, err := k.addressStringToBytes(proposal.Agent); err != nil {
		return errorsmod.Wrap(err, "agent")
	}

	return nil
}

func (k Keeper) InitGenesis(ctx context.Context, gs *escrowv1alpha1.GenesisState) error {
	if err := k.initGenesisParams(ctx, gs.Params); err != nil {
		return errorsmod.Wrap(err, "params")
	}

	if err := k.nextAgent.Set(ctx, gs.NextAgent); err != nil {
		return errorsmod.Wrap(err, "next_agent")
	}

	if err := k.initGenesisAgents(ctx, gs.Agents); err != nil {
		return errorsmod.Wrap(err, "agents")
	}

	if err := k.nextProposal.Set(ctx, gs.NextProposal); err != nil {
		return errorsmod.Wrap(err, "next_proposal")
	}

	if err := k.initGenesisProposals(ctx, gs.Proposals); err != nil {
		return errorsmod.Wrap(err, "proposals")
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
		addIndex := func(err error) error {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if err := k.initGenesisAgent(ctx, agent); err != nil {
			return addIndex(err)
		}
	}

	return nil
}

func (k Keeper) initGenesisAgent(ctx context.Context, agent *escrowv1alpha1.GenesisState_Agent) error {
	address, err := k.addressStringToBytes(agent.Address)
	if err != nil {
		return errorsmod.Wrap(err, "address")
	}

	creator, err := k.addressStringToBytes(agent.Creator)
	if err != nil {
		return errorsmod.Wrap(err, "creator")
	}

	return k.setAgent(ctx, address, creator, &escrowv1alpha1.Agent{})
}

func (k Keeper) initGenesisProposals(ctx context.Context, proposals []*escrowv1alpha1.GenesisState_Proposal) error {
	for i, proposal := range proposals {
		addIndex := func(err error) error {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if err := k.initGenesisProposal(ctx, proposal); err != nil {
			return addIndex(err)
		}
	}

	return nil
}

func (k Keeper) initGenesisProposal(ctx context.Context, proposal *escrowv1alpha1.GenesisState_Proposal) error {
	proposer, err := k.addressStringToBytes(proposal.Proposer)
	if err != nil {
		return errorsmod.Wrap(err, "proposer")
	}

	agent, err := k.addressStringToBytes(proposal.Agent)
	if err != nil {
		return errorsmod.Wrap(err, "agent")
	}

	return k.setProposal(ctx, proposal.Id, proposer, &escrowv1alpha1.Proposal{
		Agent:       agent,
		PreActions:  proposal.PreActions,
		PostActions: proposal.PostActions,
		Metadata:    proposal.Metadata,
	})
}

func (k Keeper) ExportGenesis(ctx context.Context) (*escrowv1alpha1.GenesisState, error) {
	params, err := k.exportGenesisParams(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "params")
	}

	nextAgent, err := k.nextAgent.Peek(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "next_agent")
	}

	agents, err := k.exportGenesisAgents(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "agents")
	}

	nextProposal, err := k.nextProposal.Peek(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "next_proposal")
	}

	proposals, err := k.exportGenesisProposals(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "proposals")
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
	if err := k.iterateAgents(ctx, func(address, creator sdk.AccAddress, agent escrowv1alpha1.Agent) error {
		addressStr, err := k.addressBytesToString(address)
		if err != nil {
			return errorsmod.Wrap(err, "address")
		}

		creatorStr, err := k.addressBytesToString(creator)
		if err != nil {
			return errorsmod.Wrap(err, "creator")
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
	if err := k.iterateProposals(ctx, func(id uint64, proposer sdk.AccAddress, proposal escrowv1alpha1.Proposal) error {
		proposerStr, err := k.addressBytesToString(proposer)
		if err != nil {
			return errorsmod.Wrap(err, "proposer")
		}

		agentStr, err := k.addressBytesToString(proposal.Agent)
		if err != nil {
			return errorsmod.Wrap(err, "agent")
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
