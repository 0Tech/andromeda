package internal

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

type msgServer struct {
	keeper Keeper
}

var _ internftv1alpha1.MsgServer = (*msgServer)(nil)

// NewMsgServer returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) internftv1alpha1.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func (s msgServer) Send(ctx context.Context, req *internftv1alpha1.MsgSend) (*internftv1alpha1.MsgSendResponse, error) {
	sender := sdk.MustAccAddressFromBech32(req.Sender)
	recipient := sdk.MustAccAddressFromBech32(req.Recipient)

	if err := s.keeper.Send(ctx, sender, recipient, req.Token); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventSend{
		Sender:   req.Sender,
		Receiver: req.Recipient,
		Token:      req.Token,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgSendResponse{}, nil
}

func (s msgServer) NewClass(ctx context.Context, req *internftv1alpha1.MsgNewClass) (*internftv1alpha1.MsgNewClassResponse, error) {
	if err := s.keeper.NewClass(ctx, req.Class, req.Traits); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventNewClass{
		Class:  req.Class,
		Traits: req.Traits, // TODO: sort
		Data:   req.Data,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgNewClassResponse{}, nil
}

func (s msgServer) UpdateClass(ctx context.Context, req *internftv1alpha1.MsgUpdateClass) (*internftv1alpha1.MsgUpdateClassResponse, error) {
	// TODO: data

	if err := s.keeper.UpdateClass(ctx, req.Class); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventUpdateClass{
		Class: req.Class,
		Data:  req.Data,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgUpdateClassResponse{}, nil
}

func (s msgServer) NewToken(ctx context.Context, req *internftv1alpha1.MsgNewToken) (*internftv1alpha1.MsgNewTokenResponse, error) {
	recipient := sdk.MustAccAddressFromBech32(req.Recipient)

	if err := s.keeper.NewToken(ctx, recipient, req.Token, req.Properties); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventNewToken{
		Token: req.Token,
		Properties: req.Properties, // TODO: sort
		Recipient:  req.Recipient,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgNewTokenResponse{}, nil
}

func (s msgServer) BurnToken(ctx context.Context, req *internftv1alpha1.MsgBurnToken) (*internftv1alpha1.MsgBurnTokenResponse, error) {
	owner := sdk.MustAccAddressFromBech32(req.Owner)

	if err := s.keeper.BurnToken(ctx, owner, req.Token); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventBurnToken{
		Owner: req.Owner,
		Token:   req.Token,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgBurnTokenResponse{}, nil
}

func (s msgServer) UpdateToken(ctx context.Context, req *internftv1alpha1.MsgUpdateToken) (*internftv1alpha1.MsgUpdateTokenResponse, error) {
	owner := sdk.MustAccAddressFromBech32(req.Owner)

	if err := s.keeper.validateOwner(ctx, req.Token, owner); err != nil {
		return nil, err
	}

	if err := s.keeper.UpdateToken(ctx, req.Token, req.Properties); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventUpdateToken{
		Owner: req.Owner,
		Token:        req.Token,
		Properties: req.Properties,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgUpdateTokenResponse{}, nil
}
