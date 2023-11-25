package internftv1alpha1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	for msg, name := range map[sdk.Msg]string{
		&MsgSendToken{}:        "MsgSendToken",
		&MsgCreateClass{}:    "MsgCreateClass",
		&MsgUpdateTrait{}: "MsgUpdateTrait",
		&MsgMintToken{}:     "MsgMintToken",
		&MsgBurnToken{}:     "MsgBurnToken",
		&MsgUpdateProperty{}:   "MsgUpdateProperty",
	} {
		const prefix = "andromeda/x/internft/"
		legacy.RegisterAminoMsg(cdc, msg, prefix+name)
	}
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSendToken{},
		&MsgCreateClass{},
		&MsgUpdateTrait{},
		&MsgMintToken{},
		&MsgBurnToken{},
		&MsgUpdateProperty{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
