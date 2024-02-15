package escrowv1alpha1

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgCreateAgent{},
		&MsgSubmitProposal{},
		&MsgExec{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
