package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

func (k Keeper) Create(ctx context.Context, creator sdk.AccAddress, asset string) error {
	return k.PushAsset(ctx, creator, asset)
}

func (k Keeper) Send(ctx context.Context, sender, recipient sdk.AccAddress, asset string) error {
	if err := k.PopAsset(ctx, sender, asset); err != nil {
		return err
	}

	return k.PushAsset(ctx, recipient, asset)
}

func (k Keeper) PushAsset(ctx context.Context, address sdk.AccAddress, asset string) error {
	err := k.HasAsset(ctx, address, asset)
	if !errorsmod.IsOf(err, testv1alpha1.ErrAssetNotFound) {
		if errorsmod.IsOf(err, testv1alpha1.ErrInvariantBroken) {
			return err
		}

		addrStr, err := k.addressBytesToString(address)
		if err != nil {
			return err
		}

		return errorsmod.Wrap(testv1alpha1.ErrAssetAlreadyExists.Wrap(asset), addrStr)
	}

	return k.setAsset(ctx, address, asset)
}

func (k Keeper) PopAsset(ctx context.Context, address sdk.AccAddress, asset string) error {
	if err := k.HasAsset(ctx, address, asset); err != nil {
		return err
	}

	return k.removeAsset(ctx, address, asset)
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
