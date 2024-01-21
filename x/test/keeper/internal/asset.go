package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

func (k Keeper) Create(ctx context.Context, creator sdk.AccAddress, asset string) error {
	if err := k.PushAsset(ctx, creator, asset); err != nil {
		return err
	}

	return nil
}

func (k Keeper) Send(ctx context.Context, sender, recipient sdk.AccAddress, asset string) error {
	if err := k.PopAsset(ctx, sender, asset); err != nil {
		return err
	}

	if err := k.PushAsset(ctx, recipient, asset); err != nil {
		return err
	}

	return nil
}

func (k Keeper) PushAsset(ctx context.Context, address sdk.AccAddress, asset string) error {
	if err := k.HasAsset(ctx, address, asset); err == nil {
		addrStr, err := k.addressBytesToString(address)
		if err != nil {
			return err
		}

		return errorsmod.Wrap(testv1alpha1.ErrAssetAlreadyExists.Wrap(asset), addrStr)
	}

	if err := k.setAsset(ctx, address, asset); err != nil {
		return err
	}

	return nil
}

func (k Keeper) PopAsset(ctx context.Context, address sdk.AccAddress, asset string) error {
	if err := k.HasAsset(ctx, address, asset); err != nil {
		return err
	}

	if err := k.removeAsset(ctx, address, asset); err != nil {
		return err
	}

	return nil
}

func (k Keeper) HasAsset(ctx context.Context, address sdk.AccAddress, asset string) error {
	has, err := k.assets.Has(ctx, collections.Join(address, asset))
	if err != nil {
		return testv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	if !has {
		addrStr, err := k.addressBytesToString(address)
		if err != nil {
			return err
		}

		return errorsmod.Wrap(testv1alpha1.ErrAssetNotFound.Wrap(asset), addrStr)
	}

	return nil
}

func (k Keeper) setAsset(ctx context.Context, address sdk.AccAddress, asset string) error {
	if err := k.assets.Set(ctx, collections.Join(address, asset), testv1alpha1.Asset{}); err != nil {
		return testv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}

func (k Keeper) removeAsset(ctx context.Context, address sdk.AccAddress, asset string) error {
	if err := k.assets.Remove(ctx, collections.Join(address, asset)); err != nil {
		return testv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}
