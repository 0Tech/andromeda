package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) SendToken(ctx context.Context, sender, recipient sdk.AccAddress, token *internftv1alpha1.Token) error {
	if err := k.validateOwner(ctx, token, sender); err != nil {
		return err
	}
	k.setOwner(ctx, token, recipient)

	return nil
}

func (k Keeper) validateOwner(ctx context.Context, token *internftv1alpha1.Token, owner sdk.AccAddress) error {
	if actual, err := k.getOwner(ctx, token); err != nil || !owner.Equals(actual) {
		return errorsmod.Wrap(internftv1alpha1.ErrPermissionDenied.Wrapf("not owner of %s", token), owner.String())
	}

	return nil
}

func (k Keeper) GetOwner(ctx context.Context, token *internftv1alpha1.Token) (*sdk.AccAddress, error) {
	if err := k.hasToken(ctx, token); err != nil {
		return nil, err
	}

	owner, err := k.getOwner(ctx, token)
	if err != nil {
		panic(err)
	}

	return owner, nil
}

// func (k Keeper) hasOwner(ctx context.Context, token internftv1alpha1.Token) error {
// 	_, err := k.getOwner(ctx, token)
// 	return err
// }

func (k Keeper) getOwner(ctx context.Context, token *internftv1alpha1.Token) (*sdk.AccAddress, error) {
	owner, err := k.owners.Get(ctx, collections.Join(token.ClassId, token.Id))
	if err != nil {
		if sdkerrors.ErrNotFound.Is(err) {
			err = errorsmod.Wrap(sdkerrors.ErrNotFound.Wrap("owner"), token.String())
		}

		return nil, err
	}

	return &owner, nil
}

func (k Keeper) setOwner(ctx context.Context, token *internftv1alpha1.Token, owner sdk.AccAddress) {
	if err := k.owners.Set(ctx, collections.Join(token.ClassId, token.Id), owner); err != nil {
		panic(err)
	}
}

func (k Keeper) deleteOwner(ctx context.Context, token *internftv1alpha1.Token) {
	if err := k.owners.Remove(ctx, collections.Join(token.ClassId, token.Id)); err != nil {
		panic(err)
	}
}
