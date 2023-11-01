package internftv1alpha1

import (
	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// ValidateBasic implements Msg.
func (m MsgNewClass) ValidateBasic() error {
	if err := ValidateAddress(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := m.Class.ValidateBasic(); err != nil {
		return err
	}

	if err := Traits(m.Traits).ValidateBasic(); err != nil {
		return err
	}

	if err := ValidateOperator(m.Operator, m.Class.Id); err != nil {
		return err
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgUpdateClass) ValidateBasic() error {
	if err := ValidateAddress(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := m.Class.ValidateBasic(); err != nil {
		return err
	}

	if err := ValidateOperator(m.Operator, m.Class.Id); err != nil {
		return err
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgMintNFT) ValidateBasic() error {
	if err := ValidateAddress(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := ValidateAddress(m.Recipient); err != nil {
		return errorsmod.Wrap(err, "recipient")
	}

	if err := m.Nft.ValidateBasic(); err != nil {
		return err
	}

	if err := Properties(m.Properties).ValidateBasic(); err != nil {
		return err
	}

	if err := ValidateOperator(m.Operator, m.Nft.ClassId); err != nil {
		return err
	}

	return nil
}

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

// ValidateBasic implements Msg.
func (m MsgUpdateNFT) ValidateBasic() error {
	if err := ValidateAddress(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

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
