package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
		return nil, err
	}

	return owner, nil
}

func (k Keeper) getOwner(ctx context.Context, token *internftv1alpha1.Token) (*sdk.AccAddress, error) {
	owner, err := k.owners.Get(ctx, collections.Join(token.ClassId, token.Id))
	if err != nil {
		return nil, internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &owner, nil
}

func (k Keeper) setOwner(ctx context.Context, token *internftv1alpha1.Token, owner sdk.AccAddress) error {
	if err := k.owners.Set(ctx, collections.Join(token.ClassId, token.Id), owner); err != nil {
		return internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) deleteOwner(ctx context.Context, token *internftv1alpha1.Token) error {
	if err := k.owners.Remove(ctx, collections.Join(token.ClassId, token.Id)); err != nil {
		return internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}
