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

func (s msgServer) SendToken(ctx context.Context, req *internftv1alpha1.MsgSendToken) (*internftv1alpha1.MsgSendTokenResponse, error) {
	var parsed internftv1alpha1.MsgSendTokenInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	if err := s.keeper.SendToken(ctx, parsed.Sender.AccAddress(), parsed.Recipient.AccAddress(), req.Token); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventSendToken{
		Sender:   req.Sender,
		Recipient: req.Recipient,
		Token:      req.Token,
	}); err != nil {
		return nil, internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &internftv1alpha1.MsgSendTokenResponse{}, nil
}

func (s msgServer) CreateClass(ctx context.Context, req *internftv1alpha1.MsgCreateClass) (*internftv1alpha1.MsgCreateClassResponse, error) {
	var parsed internftv1alpha1.MsgCreateClassInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	if err := internftv1alpha1.ValidateOperator(req.Operator, req.Class.Id); err != nil {
		return nil, err
	}

	if err := s.keeper.CreateClass(ctx, req.Class); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventCreateClass{
		Operator: req.Operator,
		Class:  req.Class,
	}); err != nil {
		return nil, internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &internftv1alpha1.MsgCreateClassResponse{}, nil
}

func (s msgServer) UpdateTrait(ctx context.Context, req *internftv1alpha1.MsgUpdateTrait) (*internftv1alpha1.MsgUpdateTraitResponse, error) {
	var parsed internftv1alpha1.MsgUpdateTraitInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	if err := internftv1alpha1.ValidateOperator(req.Operator, req.Class.Id); err != nil {
		return nil, err
	}

	if err := s.keeper.UpdateTrait(ctx, req.Class, req.Trait); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventUpdateTrait{
		Operator: req.Operator,
		Class: req.Class,
		Trait: req.Trait,
	}); err != nil {
		return nil, internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &internftv1alpha1.MsgUpdateTraitResponse{}, nil
}

func (s msgServer) MintToken(ctx context.Context, req *internftv1alpha1.MsgMintToken) (*internftv1alpha1.MsgMintTokenResponse, error) {
	var parsed internftv1alpha1.MsgMintTokenInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	if err := internftv1alpha1.ValidateOperator(req.Operator, req.Token.ClassId); err != nil {
		return nil, err
	}

	if err := s.keeper.MintToken(ctx, parsed.Operator.AccAddress(), req.Token); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventMintToken{
		Operator: req.Operator,
		Token: req.Token,
	}); err != nil {
		return nil, internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &internftv1alpha1.MsgMintTokenResponse{}, nil
}

func (s msgServer) BurnToken(ctx context.Context, req *internftv1alpha1.MsgBurnToken) (*internftv1alpha1.MsgBurnTokenResponse, error) {
	var parsed internftv1alpha1.MsgBurnTokenInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	if err := s.keeper.BurnToken(ctx, parsed.Owner.AccAddress(), req.Token); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventBurnToken{
		Owner: req.Owner,
		Token:   req.Token,
	}); err != nil {
		return nil, internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &internftv1alpha1.MsgBurnTokenResponse{}, nil
}

func (s msgServer) UpdateProperty(ctx context.Context, req *internftv1alpha1.MsgUpdateProperty) (*internftv1alpha1.MsgUpdatePropertyResponse, error) {
	var parsed internftv1alpha1.MsgUpdatePropertyInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	if err := internftv1alpha1.ValidateOperator(req.Operator, req.Token.ClassId); err != nil {
		return nil, err
	}

	if err := s.keeper.UpdateProperty(ctx, req.Token, req.Property); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventUpdateProperty{
		Operator: req.Operator,
		Token:        req.Token,
		Property: req.Property,
	}); err != nil {
		return nil, internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &internftv1alpha1.MsgUpdatePropertyResponse{}, nil
}
