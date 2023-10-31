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

// Send defines a method to send an nft from one account to another account.
func (s msgServer) Send(ctx context.Context, req *internftv1alpha1.MsgSend) (*internftv1alpha1.MsgSendResponse, error) {
	sender := sdk.MustAccAddressFromBech32(req.Sender)
	recipient := sdk.MustAccAddressFromBech32(req.Recipient)

	if err := s.keeper.Send(ctx, sender, recipient, req.Nft); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventSend{
		Sender:   req.Sender,
		Receiver: req.Recipient,
		Nft:      req.Nft,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgSendResponse{}, nil
}

// NewClass defines a method to create a class.
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

// UpdateClass defines a method to update a class.
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

// MintNFT defines a method to mint an nft.
func (s msgServer) MintNFT(ctx context.Context, req *internftv1alpha1.MsgMintNFT) (*internftv1alpha1.MsgMintNFTResponse, error) {
	recipient := sdk.MustAccAddressFromBech32(req.Recipient)

	if err := s.keeper.MintNFT(ctx, recipient, req.Nft, req.Properties); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventMintNFT{
		Nft: req.Nft,
		Properties: req.Properties, // TODO: sort
		Recipient:  req.Recipient,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgMintNFTResponse{}, nil
}

// BurnNFT defines a method to burn an nft.
func (s msgServer) BurnNFT(ctx context.Context, req *internftv1alpha1.MsgBurnNFT) (*internftv1alpha1.MsgBurnNFTResponse, error) {
	owner := sdk.MustAccAddressFromBech32(req.Owner)

	if err := s.keeper.BurnNFT(ctx, owner, req.Nft); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventBurnNFT{
		Owner: req.Owner,
		Nft:   req.Nft,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgBurnNFTResponse{}, nil
}

// UpdateNFT defines a method to update an nft.
func (s msgServer) UpdateNFT(ctx context.Context, req *internftv1alpha1.MsgUpdateNFT) (*internftv1alpha1.MsgUpdateNFTResponse, error) {
	if err := s.keeper.UpdateNFT(ctx, req.Nft, req.Properties); err != nil {
		return nil, err
	}

	if err := sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&internftv1alpha1.EventUpdateNFT{
		Nft:        req.Nft,
		Properties: req.Properties,
	}); err != nil {
		panic(err)
	}

	return &internftv1alpha1.MsgUpdateNFTResponse{}, nil
}
