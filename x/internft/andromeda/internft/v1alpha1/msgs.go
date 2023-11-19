package internftv1alpha1

import (
	errorsmod "cosmossdk.io/errors"
)

func (m MsgSend) ValidateCompatibility() error {
	if m.Sender == "" {
		return ErrUnimplemented.Wrap("nil sender")
	}

	if m.Recipient == "" {
		return ErrUnimplemented.Wrap("nil recipient")
	}

	if m.Token == nil {
		return ErrUnimplemented.Wrap("nil token")
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
		return ErrUnimplemented.Wrap("nil operator")
	}

	if m.Class == nil {
		return ErrUnimplemented.Wrap("nil class")
	}

	if m.Traits == nil {
		return ErrUnimplemented.Wrap("nil traits")
	}

	// TODO(@0Tech): data
	// if m.Data == nil {
	// 	return ErrUnimplemented.Wrap("nil data")
	// }

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

	return nil
}

func (m MsgUpdateClass) ValidateCompatibility() error {
	if m.Operator == "" {
		return ErrUnimplemented.Wrap("nil operator")
	}

	if m.Class == nil {
		return ErrUnimplemented.Wrap("nil class")
	}

	// TODO(@0Tech): data
	// if m.Data == nil {
	// 	return ErrUnimplemented.Wrap("nil data")
	// }

	return nil
}

func (m MsgUpdateClass) ValidateBasic() error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateAddress(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := m.Class.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "class")
	}

	return nil
}

func (m MsgNewToken) ValidateCompatibility() error {
	if m.Operator == "" {
		return ErrUnimplemented.Wrap("nil operator")
	}

	if m.Recipient == "" {
		return ErrUnimplemented.Wrap("nil recipient")
	}

	if m.Token == nil {
		return ErrUnimplemented.Wrap("nil token")
	}

	if m.Properties == nil {
		return ErrUnimplemented.Wrap("nil properties")
	}

	return nil
}

func (m MsgNewToken) ValidateBasic() error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

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

	return nil
}

func (m MsgBurnToken) ValidateCompatibility() error {
	if m.Owner == "" {
		return ErrUnimplemented.Wrap("nil owner")
	}

	if m.Token == nil {
		return ErrUnimplemented.Wrap("nil token")
	}

	return nil
}

func (m MsgBurnToken) ValidateBasic() error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateAddress(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := m.Token.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (m MsgUpdateToken) ValidateCompatibility() error {
	if m.Owner == "" {
		return ErrUnimplemented.Wrap("nil owner")
	}

	if m.Token == nil {
		return ErrUnimplemented.Wrap("nil token")
	}

	if m.Properties == nil {
		return ErrUnimplemented.Wrap("nil properties")
	}

	return nil
}

func (m MsgUpdateToken) ValidateBasic() error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateAddress(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := m.Token.ValidateBasic(); err != nil {
		return err
	}

	if err := Properties(m.Properties).ValidateBasic(); err != nil {
		return err
	}

	return nil
}
