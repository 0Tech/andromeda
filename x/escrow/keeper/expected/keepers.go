package expected

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	MessageRouter interface {
		Handler(msg sdk.Msg) MsgServiceHandler
	}
	MsgServiceHandler = func(ctx sdk.Context, req sdk.Msg) (*sdk.Result, error)

	AuthKeeper interface {
		HasAccount(context.Context, sdk.AccAddress) bool
		NewAccount(context.Context, sdk.AccountI) sdk.AccountI
		SetAccount(context.Context, sdk.AccountI)
	}
)
