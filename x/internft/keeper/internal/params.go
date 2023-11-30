package internal

import (
	"context"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) UpdateParams(ctx context.Context, params *internftv1alpha1.Params) error {
	return k.setParams(ctx, params)
}

func (k Keeper) GetParams(ctx context.Context) (*internftv1alpha1.Params, error) {
	params, err := k.params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

func (k Keeper) setParams(ctx context.Context, params *internftv1alpha1.Params) error {
	if err := k.params.Set(ctx, *params); err != nil {
		return internftv1alpha1.ErrInvariantBroken.Wrap(err.Error())
	}

	return nil
}
