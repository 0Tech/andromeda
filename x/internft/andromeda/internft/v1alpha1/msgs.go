package internftv1alpha1

import (
	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (m MsgSend) ValidateCompatibility() error {
	if m.Sender == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil sender")
	}

	if m.Recipient == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil recipient")
	}

	if m.Token == nil {
		return sdkerrors.ErrNotSupported.Wrap("nil token")
	}

	return nil
}

func (m MsgSend) ValidateBasic() error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateAddress(m.Sender); err != nil {
		return errorsmod.Wrap(err, "sender")
	}

	if err := ValidateAddress(m.Recipient); err != nil {
		return errorsmod.Wrap(err, "recipient")
	}

	if err := m.Token.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "token")
	}

	return nil
}

func (m MsgNewClass) ValidateCompatibility() error {
	if m.Operator == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil operator")
	}

	if m.Class == nil {
		return sdkerrors.ErrNotSupported.Wrap("nil class")
	}

	if m.Traits == nil {
		return sdkerrors.ErrNotSupported.Wrap("nil traits")
	}

	// TODO: data
	return nil
}

func (m MsgNewClass) ValidateBasic() error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateAddress(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := m.Class.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "class")
	}

	if err := Traits(m.Traits).ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "traits")
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
		return errorsmod.Wrap(err, "class")
	}

	if err := ValidateOperator(m.Operator, m.Class.Id); err != nil {
		return err
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgNewToken) ValidateBasic() error {
	if err := ValidateAddress(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := ValidateAddress(m.Recipient); err != nil {
		return errorsmod.Wrap(err, "recipient")
	}

	if err := m.Token.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "token")
	}

	if err := Properties(m.Properties).ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "properties")
	}

	if err := ValidateOperator(m.Operator, m.Token.ClassId); err != nil {
		return err
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgBurnToken) ValidateBasic() error {
	if err := ValidateAddress(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := m.Token.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// ValidateBasic implements Msg.
func (m MsgUpdateToken) ValidateBasic() error {
	if err := ValidateAddress(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := m.Token.ValidateBasic(); err != nil {
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
