package internal

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

type queryServer struct {
	keeper Keeper
}

var _ escrowv1alpha1.QueryServer = (*queryServer)(nil)

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) escrowv1alpha1.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

func (s queryServer) Params(ctx context.Context, req *escrowv1alpha1.QueryParamsRequest) (*escrowv1alpha1.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	_, err := s.keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &escrowv1alpha1.QueryParamsResponse{
		Params: &escrowv1alpha1.QueryParamsResponse_Params{},
	}, nil
}

func (s queryServer) Agent(ctx context.Context, req *escrowv1alpha1.QueryAgentRequest) (*escrowv1alpha1.QueryAgentResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil address")
	}

	address, err := s.keeper.addressStringToBytes(req.Address)
	if err != nil {
		return nil, errorsmod.Wrap(err, "address")
	}

	creator, _, err := s.keeper.GetAgent(ctx, address)
	if err != nil {
		return nil, err
	}

	creatorStr, err := s.keeper.addressBytesToString(creator)
	if err != nil {
		return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "creator")
	}

	return &escrowv1alpha1.QueryAgentResponse{
		Agent: &escrowv1alpha1.QueryAgentResponse_Agent{
			Address: req.Address,
			Creator: creatorStr,
		},
	}, nil
}

func (s queryServer) Agents(ctx context.Context, req *escrowv1alpha1.QueryAgentsRequest) (*escrowv1alpha1.QueryAgentsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	agents, pageRes, err := query.CollectionPaginate(ctx, s.keeper.agents, req.Pagination, func(key collections.Pair[sdk.AccAddress, sdk.AccAddress], _ escrowv1alpha1.Agent) (*escrowv1alpha1.QueryAgentsResponse_Agent, error) {
		creator := key.K1()
		address := key.K2()

		addressStr, err := s.keeper.addressBytesToString(address)
		if err != nil {
			return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "address")
		}

		creatorStr, err := s.keeper.addressBytesToString(creator)
		if err != nil {
			return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "creator")
		}

		return &escrowv1alpha1.QueryAgentsResponse_Agent{
			Address: addressStr,
			Creator: creatorStr,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return &escrowv1alpha1.QueryAgentsResponse{
		Agents:     agents,
		Pagination: pageRes,
	}, nil
}

func (s queryServer) Proposal(ctx context.Context, req *escrowv1alpha1.QueryProposalRequest) (*escrowv1alpha1.QueryProposalResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Id == 0 {
		return nil, escrowv1alpha1.ErrUnimplemented.Wrap("nil id")
	}

	proposer, proposal, err := s.keeper.GetProposal(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	proposerStr, err := s.keeper.addressBytesToString(proposer)
	if err != nil {
		return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "proposer")
	}

	agentStr, err := s.keeper.addressBytesToString(proposal.Agent)
	if err != nil {
		return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "agent")
	}

	return &escrowv1alpha1.QueryProposalResponse{
		Proposal: &escrowv1alpha1.QueryProposalResponse_Proposal{
			Id:          req.Id,
			Proposer:    proposerStr,
			Agent:       agentStr,
			PreActions:  proposal.PreActions,
			PostActions: proposal.PostActions,
		},
	}, nil
}

func (s queryServer) Proposals(ctx context.Context, req *escrowv1alpha1.QueryProposalsRequest) (*escrowv1alpha1.QueryProposalsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	proposals, pageRes, err := query.CollectionPaginate(ctx, s.keeper.proposals, req.Pagination, func(key collections.Pair[sdk.AccAddress, uint64], value escrowv1alpha1.Proposal) (*escrowv1alpha1.QueryProposalsResponse_Proposal, error) {
		proposer := key.K1()
		id := key.K2()
		proposal := value

		proposerStr, err := s.keeper.addressBytesToString(proposer)
		if err != nil {
			return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "proposer")
		}

		agentStr, err := s.keeper.addressBytesToString(value.Agent)
		if err != nil {
			return nil, errorsmod.Wrap(escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error()), "agent")
		}

		s.keeper.fixActions(&proposal)

		return &escrowv1alpha1.QueryProposalsResponse_Proposal{
			Id:          id,
			Proposer:    proposerStr,
			Agent:       agentStr,
			PreActions:  proposal.PreActions,
			PostActions: proposal.PostActions,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return &escrowv1alpha1.QueryProposalsResponse{
		Proposals:  proposals,
		Pagination: pageRes,
	}, nil
}
