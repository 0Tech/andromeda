package internal

import (
	"context"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) GetParams(ctx context.Context) *internftv1alpha1.Params {
	params, err := k.params.Get(ctx)
	if err != nil {
		panic(err)
	}

	return &params
}

func (k Keeper) SetParams(ctx context.Context, params *internftv1alpha1.Params) {
	if err := k.params.Set(ctx, *params); err != nil {
		panic(err)
	}
}
