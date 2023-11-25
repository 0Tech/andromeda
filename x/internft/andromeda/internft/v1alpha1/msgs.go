package internftv1alpha1

import (
	errorsmod "cosmossdk.io/errors"
)

func (m MsgSendToken) ValidateCompatibility() error {
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

type MsgSendTokenInternal struct {
	Sender Address
	Recipient Address
	Token TokenInternal
}

func (mi *MsgSendTokenInternal) Parse(m MsgSendToken) error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := mi.Sender.Parse(m.Sender); err != nil {
		return errorsmod.Wrap(err, "sender")
	}

	if err := mi.Recipient.Parse(m.Recipient); err != nil {
		return errorsmod.Wrap(err, "recipient")
	}

	if err := mi.Token.Parse(*m.Token); err != nil {
		return errorsmod.Wrap(err, "token")
	}

	return nil
}

func (m MsgCreateClass) ValidateCompatibility() error {
	if m.Operator == "" {
		return ErrUnimplemented.Wrap("nil operator")
	}

	if m.Class == nil {
		return ErrUnimplemented.Wrap("nil class")
	}

	return nil
}

type MsgCreateClassInternal struct {
	Operator Address
	Class ClassInternal
}

func (mi *MsgCreateClassInternal) Parse(m MsgCreateClass) error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := mi.Operator.Parse(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := mi.Class.Parse(*m.Class); err != nil {
		return errorsmod.Wrap(err, "class")
	}

	return nil
}

func (m MsgUpdateTrait) ValidateCompatibility() error {
	if m.Operator == "" {
		return ErrUnimplemented.Wrap("nil operator")
	}

	if m.Class == nil {
		return ErrUnimplemented.Wrap("nil class")
	}

	if m.Trait == nil {
		return ErrUnimplemented.Wrap("nil trait")
	}

	return nil
}

type MsgUpdateTraitInternal struct {
	Operator Address
	Class ClassInternal
	Trait TraitInternal
}

func (mi *MsgUpdateTraitInternal) Parse(m MsgUpdateTrait) error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := mi.Operator.Parse(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := mi.Class.Parse(*m.Class); err != nil {
		return errorsmod.Wrap(err, "class")
	}

	if err := mi.Trait.Parse(*m.Trait); err != nil {
		return errorsmod.Wrap(err, "trait")
	}

	return nil
}

func (m MsgMintToken) ValidateCompatibility() error {
	if m.Operator == "" {
		return ErrUnimplemented.Wrap("nil operator")
	}

	if m.Token == nil {
		return ErrUnimplemented.Wrap("nil token")
	}

	return nil
}

type MsgMintTokenInternal struct {
	Operator Address
	Token TokenInternal
}

func (mi *MsgMintTokenInternal) Parse(m MsgMintToken) error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := mi.Operator.Parse(m.Operator); err != nil {
		return errorsmod.Wrap(err, "operator")
	}

	if err := mi.Token.Parse(*m.Token); err != nil {
		return errorsmod.Wrap(err, "token")
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

type MsgBurnTokenInternal struct {
	Owner Address
	Token TokenInternal
}

func (mi *MsgBurnTokenInternal) Parse(m MsgBurnToken) error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := mi.Owner.Parse(m.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := mi.Token.Parse(*m.Token); err != nil {
		return errorsmod.Wrap(err, "token")
	}

	return nil
}

func (m MsgUpdateProperty) ValidateCompatibility() error {
	if m.Operator == "" {
		return ErrUnimplemented.Wrap("nil operator")
	}

	if m.Token == nil {
		return ErrUnimplemented.Wrap("nil token")
	}

	if m.Property == nil {
		return ErrUnimplemented.Wrap("nil property")
	}

	return nil
}

type MsgUpdatePropertyInternal struct {
	Operator Address
	Token TokenInternal
	Property PropertyInternal
}

func (mi *MsgUpdatePropertyInternal) Parse(m MsgUpdateProperty) error {
	if err := m.ValidateCompatibility(); err != nil {
		return err
	}

	if err := mi.Operator.Parse(m.Operator); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	if err := mi.Token.Parse(*m.Token); err != nil {
		return errorsmod.Wrap(err, "token")
	}

	if err := mi.Property.Parse(*m.Property); err != nil {
		return errorsmod.Wrap(err, "property")
	}

	return nil
}
