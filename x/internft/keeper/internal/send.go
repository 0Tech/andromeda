package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) Send(ctx context.Context, sender, recipient sdk.AccAddress, nft internftv1alpha1.NFT) error {
	if err := k.validateOwner(ctx, nft, sender); err != nil {
		return err
	}
	k.setOwner(ctx, nft, recipient)

	return nil
}

func (k Keeper) validateOwner(ctx context.Context, nft internftv1alpha1.NFT, owner sdk.AccAddress) error {
	if actual, err := k.getOwner(ctx, nft); err != nil || !owner.Equals(actual) {
		return errorsmod.Wrap(internftv1alpha1.ErrInsufficientNFT.Wrap("not owns nft"), nft.String())
	}

	return nil
}

func (k Keeper) GetOwner(ctx context.Context, nft internftv1alpha1.NFT) (*sdk.AccAddress, error) {
	if err := k.hasNFT(ctx, nft); err != nil {
		return nil, err
	}

	owner, err := k.getOwner(ctx, nft)
	if err != nil {
		panic(err)
	}

	return owner, nil
}

// func (k Keeper) hasOwner(ctx context.Context, nft internftv1alpha1.NFT) error {
// 	_, err := k.getOwner(ctx, nft)
// 	return err
// }

func (k Keeper) getOwner(ctx context.Context, nft internftv1alpha1.NFT) (*sdk.AccAddress, error) {
	owner, err := k.owners.Get(ctx, collections.Join(nft.ClassId, nft.Id))
	if err != nil {
		if sdkerrors.ErrNotFound.Is(err) {
			err = errorsmod.Wrap(sdkerrors.ErrNotFound.Wrap("owner"), nft.String())
		}

		return nil, err
	}

	return &owner, nil
}

func (k Keeper) setOwner(ctx context.Context, nft internftv1alpha1.NFT, owner sdk.AccAddress) {
	if err := k.owners.Set(ctx, collections.Join(nft.ClassId, nft.Id), owner); err != nil {
		panic(err)
	}
}

func (k Keeper) deleteOwner(ctx context.Context, nft internftv1alpha1.NFT) {
	if err := k.owners.Remove(ctx, collections.Join(nft.ClassId, nft.Id)); err != nil {
		panic(err)
	}
}
