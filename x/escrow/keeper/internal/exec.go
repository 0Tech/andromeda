package internal

import (
	"context"

	"cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
)

func (k Keeper) Exec(ctx context.Context, agents []sdk.AccAddress, actions []*codectypes.Any) error {
	var postActions []*codectypes.Any
	for _, agent := range agents {
		proposal, err := k.GetProposal(ctx, agent)
		if err != nil {
			return err
		}

		postActions = append(postActions, proposal.PostActions...)

		if err := k.removeProposal(ctx, agent); err != nil {
			return err
		}
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
			actions: postActions,
		},
	} {
		if err := k.executeActions(ctx, phase.actions); err != nil {
			return errors.Wrap(err, phase.name)
		}
	}

	return nil
}

func (k Keeper) executeActions(ctx context.Context, actions []*codectypes.Any) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for i, action := range actions {
		msg, err := k.anyToMsg(*action)
		if err != nil {
			return indexedError(err, i)
		}

		handler := k.router.Handler(msg)
		if handler == nil {
			return indexedError(escrowv1alpha1.ErrInvalidMessage.Wrap("handler not found"), i)
		}

		result, err := handler(sdkCtx, msg)
		if err != nil {
			return indexedError(err, i)
		}

		sdkCtx.EventManager().EmitEvents(result.GetEvents())
	}

	return nil
}

func (k Keeper) anyToMsg(any codectypes.Any) (sdk.Msg, error) {
	var msg sdk.Msg
	if err := k.cdc.UnpackAny(&any, &msg); err != nil {
		return nil, escrowv1alpha1.ErrInvalidMessage.Wrap(err.Error())
	}

	return msg, nil
}
