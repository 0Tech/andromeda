package internftv1alpha1

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	"github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1/codec"
)

var (
	_ sdk.Msg            = (*MsgSend)(nil)
	_ legacytx.LegacyMsg = (*MsgSend)(nil)
)

// ValidateBasic implements Msg.
func (m MsgSend) ValidateBasic() error {
	if err := ValidateAddress(m.Sender); err != nil {
		return errorsmod.Wrap(err, "sender")
	}

	if err := ValidateAddress(m.Recipient); err != nil {
		return errorsmod.Wrap(err, "recipient")
	}

	if err := m.Nft.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgSend) GetSigners() []sdk.AccAddress {
	signer := sdk.MustAccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

var (
	_ sdk.Msg            = (*MsgNewClass)(nil)
	_ legacytx.LegacyMsg = (*MsgNewClass)(nil)
)

// ValidateBasic implements Msg.
func (m MsgNewClass) ValidateBasic() error {
	if err := ValidateAddress(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := Traits(m.Traits).ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgNewClass) GetSigners() []sdk.AccAddress {
	signer := sdk.MustAccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgNewClass) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

var (
	_ sdk.Msg            = (*MsgUpdateClass)(nil)
	_ legacytx.LegacyMsg = (*MsgUpdateClass)(nil)
)

// ValidateBasic implements Msg.
func (m MsgUpdateClass) ValidateBasic() error {
	if err := ValidateClassID(m.ClassId); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgUpdateClass) GetSigners() []sdk.AccAddress {
	signer := ClassOwner(m.ClassId)
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgUpdateClass) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

var (
	_ sdk.Msg            = (*MsgMintNFT)(nil)
	_ legacytx.LegacyMsg = (*MsgMintNFT)(nil)
)

// ValidateBasic implements Msg.
func (m MsgMintNFT) ValidateBasic() error {
	if err := ValidateClassID(m.ClassId); err != nil {
		return err
	}

	if err := Properties(m.Properties).ValidateBasic(); err != nil {
		return err
	}

	if err := ValidateAddress(m.Recipient); err != nil {
		return errorsmod.Wrap(err, "recipient")
	}

	return nil
}

// GetSigners implements Msg
func (m MsgMintNFT) GetSigners() []sdk.AccAddress {
	signer := ClassOwner(m.ClassId)
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgMintNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

var (
	_ sdk.Msg            = (*MsgBurnNFT)(nil)
	_ legacytx.LegacyMsg = (*MsgBurnNFT)(nil)
)

// ValidateBasic implements Msg.
func (m MsgBurnNFT) ValidateBasic() error {
	if err := ValidateAddress(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := m.Nft.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgBurnNFT) GetSigners() []sdk.AccAddress {
	signer := sdk.MustAccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgBurnNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}

var (
	_ sdk.Msg            = (*MsgUpdateNFT)(nil)
	_ legacytx.LegacyMsg = (*MsgUpdateNFT)(nil)
)

// ValidateBasic implements Msg.
func (m MsgUpdateNFT) ValidateBasic() error {
	if err := m.Nft.ValidateBasic(); err != nil {
		return err
	}

	if len(m.Properties) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("empty properties")
	}

	if err := Properties(m.Properties).ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgUpdateNFT) GetSigners() []sdk.AccAddress {
	signer := ClassOwner(m.Nft.ClassId)
	return []sdk.AccAddress{signer}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgUpdateNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(codec.ModuleCdc.MustMarshalJSON(&m))
}
