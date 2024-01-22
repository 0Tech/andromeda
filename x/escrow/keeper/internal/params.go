package internal

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) UpdateParams(ctx context.Context, newParams *escrowv1alpha1.Params) error {
	oldParams, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	if newParams.MaxMetadataLength < oldParams.MaxMetadataLength {
		// TODO: define error
		return sdkerrors.ErrInvalidRequest.Wrap("cannot lower max_metadata_length")
	}

	return k.setParams(ctx, newParams)
}

func (k Keeper) GetParams(ctx context.Context) (*escrowv1alpha1.Params, error) {
	params, err := k.params.Get(ctx)
	if err != nil {
		return nil, escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return &params, nil
}

func (k Keeper) setParams(ctx context.Context, params *escrowv1alpha1.Params) error {
	if err := k.params.Set(ctx, *params); err != nil {
		return escrowv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}
