package internal

import (
	"context"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

type msgServer struct {
	keeper Keeper
}

var _ testv1alpha1.MsgServer = (*msgServer)(nil)

// NewMsgServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) testv1alpha1.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func (s msgServer) Create(ctx context.Context, req *testv1alpha1.MsgCreate) (*testv1alpha1.MsgCreateResponse, error) {
	if req.Creator == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil creator")
	}

	if req.Asset == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil asset")
	}

	creator, err := s.keeper.addressStringToBytes(req.Creator)
	if err != nil {
		return nil, errors.Wrap(err, "creator")
	}

	if err := s.keeper.Create(ctx, creator, req.Asset); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&testv1alpha1.EventCreate{
		Creator: req.Creator,
		Asset:   req.Asset,
	}); err != nil {
		return nil, testv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &testv1alpha1.MsgCreateResponse{}, nil
}

func (s msgServer) Send(ctx context.Context, req *testv1alpha1.MsgSend) (*testv1alpha1.MsgSendResponse, error) {
	if req.Sender == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil sender")
	}

	if req.Recipient == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil recipient")
	}

	if req.Asset == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil asset")
	}

	sender, err := s.keeper.addressStringToBytes(req.Sender)
	if err != nil {
		return nil, errors.Wrap(err, "sender")
	}

	recipient, err := s.keeper.addressStringToBytes(req.Recipient)
	if err != nil {
		return nil, errors.Wrap(err, "recipient")
	}

	if err := s.keeper.Send(ctx, sender, recipient, req.Asset); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&testv1alpha1.EventSend{
		Sender:    req.Sender,
		Recipient: req.Recipient,
		Asset:     req.Asset,
	}); err != nil {
		return nil, testv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &testv1alpha1.MsgSendResponse{}, nil
}
