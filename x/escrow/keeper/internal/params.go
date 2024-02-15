package internal

import (
	"context"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) UpdateParams(ctx context.Context, newParams *escrowv1alpha1.Params) error {
	// no transition rules yet.
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
