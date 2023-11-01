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
		&MsgSend{}:        "MsgSend",
		&MsgNewClass{}:    "MsgNewClass",
		&MsgUpdateClass{}: "MsgUpdateClass",
		&MsgMintNFT{}:     "MsgMintNFT",
		&MsgBurnNFT{}:     "MsgBurnNFT",
		&MsgUpdateNFT{}:   "MsgUpdateNFT",
	} {
		const prefix = "andromeda/x/internft/"
		legacy.RegisterAminoMsg(cdc, msg, prefix+name)
	}
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
		&MsgNewClass{},
		&MsgUpdateClass{},
		&MsgMintNFT{},
		&MsgBurnNFT{},
		&MsgUpdateNFT{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
