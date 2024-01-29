package internal

import (
	"context"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

type msgServer struct {
	keeper Keeper
}

var _ escrowv1alpha1.MsgServer = (*msgServer)(nil)

// NewMsgServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) escrowv1alpha1.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func (s msgServer) UpdateParams(ctx context.Context, req *escrowv1alpha1.MsgUpdateParams) (*escrowv1alpha1.MsgUpdateParamsResponse, error) {
	if req.Authority == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil authority")
	}

	if req.MaxMetadataLength == 0 {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil max_metadata_length")
	}

	authority, err := s.keeper.addressStringToBytes(req.Authority)
	if err != nil {
		return nil, errors.Wrap(err, "authority")
	}

	if err := s.keeper.validateAuthority(authority); err != nil {
		return nil, err
	}

	if err := s.keeper.UpdateParams(ctx, &escrowv1alpha1.Params{
		MaxMetadataLength: req.MaxMetadataLength,
	}); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&escrowv1alpha1.EventUpdateParams{
		Authority:         req.Authority,
		MaxMetadataLength: req.MaxMetadataLength,
	}); err != nil {
		return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &escrowv1alpha1.MsgUpdateParamsResponse{}, nil
}

func (s msgServer) CreateAgent(ctx context.Context, req *escrowv1alpha1.MsgCreateAgent) (*escrowv1alpha1.MsgCreateAgentResponse, error) {
	if req.Creator == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil creator")
	}

	creator, err := s.keeper.addressStringToBytes(req.Creator)
	if err != nil {
		return nil, errors.Wrap(err, "creator")
	}

	agent, err := s.keeper.CreateAgent(ctx, creator)
	if err != nil {
		return nil, err
	}

	agentStr, err := s.keeper.addressBytesToString(agent)
	if err != nil {
		return nil, errors.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "agent")
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&escrowv1alpha1.EventCreateAgent{
		Agent:   agentStr,
		Creator: req.Creator,
	}); err != nil {
		return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &escrowv1alpha1.MsgCreateAgentResponse{Agent: agentStr}, nil
}

func (s msgServer) SubmitProposal(ctx context.Context, req *escrowv1alpha1.MsgSubmitProposal) (*escrowv1alpha1.MsgSubmitProposalResponse, error) {
	if req.Proposer == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil proposer")
	}

	if req.Agent == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil agent")
	}

	if req.PreActions == nil {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil pre_actions")
	}

	if req.PostActions == nil {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil post_actions")
	}

	if req.Metadata == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil metadata")
	}

	proposer, err := s.keeper.addressStringToBytes(req.Proposer)
	if err != nil {
		return nil, errors.Wrap(err, "proposer")
	}

	agent, err := s.keeper.addressStringToBytes(req.Agent)
	if err != nil {
		return nil, errors.Wrap(err, "agent")
	}

	signers := []sdk.AccAddress{proposer, agent}

	if err := s.keeper.validateActions(req.PreActions, signers); err != nil {
		return nil, errors.Wrap(err, "pre_actions")
	}

	if err := s.keeper.validateActions(req.PostActions, signers); err != nil {
		return nil, errors.Wrap(err, "post_actions")
	}

	if err := s.keeper.SubmitProposal(ctx, proposer, agent, req.PreActions, req.PostActions, req.Metadata); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&escrowv1alpha1.EventSubmitProposal{
		Proposer:    req.Proposer,
		Agent:       req.Agent,
		PreActions:  req.PreActions,
		PostActions: req.PostActions,
		Metadata:    req.Metadata,
	}); err != nil {
		return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &escrowv1alpha1.MsgSubmitProposalResponse{}, nil
}

func (s msgServer) Exec(ctx context.Context, req *escrowv1alpha1.MsgExec) (*escrowv1alpha1.MsgExecResponse, error) {
	if req.Executor == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil executor")
	}

	if req.Agent == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil agent")
	}

	if req.Actions == nil {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil actions")
	}

	executor, err := s.keeper.addressStringToBytes(req.Executor)
	if err != nil {
		return nil, errors.Wrap(err, "executor")
	}

	agent, err := s.keeper.addressStringToBytes(req.Agent)
	if err != nil {
		return nil, errors.Wrap(err, "agent")
	}

	signers := []sdk.AccAddress{executor, agent}

	if err := s.keeper.validateActions(req.Actions, signers); err != nil {
		return nil, errors.Wrap(err, "actions")
	}

	if err := s.keeper.Exec(ctx, executor, agent, req.Actions); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&escrowv1alpha1.EventExec{
		Executor: req.Executor,
		Agent:    req.Agent,
		Actions:  req.Actions,
	}); err != nil {
		return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &escrowv1alpha1.MsgExecResponse{}, nil
}
