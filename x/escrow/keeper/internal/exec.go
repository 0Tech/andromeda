package internal

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) Exec(ctx context.Context, id uint64, executor, agent sdk.AccAddress, actions []*codectypes.Any) error {
	_, proposal, err := k.GetProposal(ctx, id)
	if err != nil {
		return err
	}

	if !agent.Equals(sdk.AccAddress(proposal.Agent)) {
		return escrowv1alpha1.ErrPermissionDenied.Wrap("agent differs")
	}

	for _, phase := range []struct {
		name    string
		actions []*codectypes.Any
	}{
		{
			name:    "actions",
			actions: actions,
		},
		{
			name:    "post_actions",
			actions: proposal.PostActions,
		},
	} {
		if err := k.executeActions(ctx, phase.actions); err != nil {
			return errorsmod.Wrap(err, phase.name)
		}
	}

	if err := k.removeProposal(ctx, id); err != nil {
		// TODO: invariant broken error handler
		return err
	}

	return nil
}

func (k Keeper) executeActions(ctx context.Context, actions []*codectypes.Any) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for i, action := range actions {
		addIndex := func(err error) error {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		msg, err := k.anyToMsg(*action)
		if err != nil {
			return addIndex(err)
		}

		handler := k.router.Handler(msg)
		if handler == nil {
			panic("TODO: error")
		}

		result, err := handler(sdkCtx, msg)
		if err != nil {
			return addIndex(err)
		}

		sdkCtx.EventManager().EmitEvents(result.GetEvents())
	}

	return nil
}

func (k Keeper) anyToMsg(any codectypes.Any) (sdk.Msg, error) {
	var msg sdk.Msg
	if err := k.cdc.UnpackAny(&any, &msg); err != nil {
		// TODO: define error
		return nil, err
	}

	return msg, nil
}