package internftv1alpha1

import (
	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		Classes: []*GenesisClass{},
	}
}

func DefaultParams() *Params {
	return &Params{}
}

type Properties []*Property

func (ps Properties) ValidateBasic() error {
	seenID := ""
	for i, property := range ps {
		if property == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("nil property"), "index %d", i)
		}

		if err := (&PropertyInternal{}).Parse(*property); err != nil {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if !(property.TraitId > seenID) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("unsorted property"), "index %d", i)
		}
		seenID = property.TraitId
	}

	return nil
}

func (t GenesisToken) ValidateCompatibility() error {
	if t.Id == "" {
		return ErrUnimplemented.Wrap("nil id")
	}

	if t.Properties == nil {
		return ErrUnimplemented.Wrap("nil properties")
	}

	if t.Owner == "" {
		return ErrUnimplemented.Wrap("nil owner")
	}

	return nil
}

func (t GenesisToken) ValidateBasic(traitIDs map[string]struct{}) error {
	if err := t.ValidateCompatibility(); err != nil {
		return err
	}

	var id TokenID
	if err := id.Parse(t.Id); err != nil {
		return err
	}

	if err := Properties(t.Properties).ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "properties")
	}

	for i, property := range t.Properties {
		if _, hasTrait := traitIDs[property.TraitId]; !hasTrait {
			return errorsmod.Wrap(errorsmod.Wrapf(ErrTraitNotFound.Wrap(property.TraitId), "index %d", i), "properties")
		}
	}

	if err := (&Address{}).Parse(t.Owner); err != nil {
		return errorsmod.Wrap(err, "owner")
	}

	return nil
}

type GenesisTokens []*GenesisToken

func (ts GenesisTokens) ValidateBasic(traitIDs map[string]struct{}) error {
	seenID := ""
	for i, token := range ts {
		if token == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("nil token"), "index %d", i)
		}

		if err := token.ValidateBasic(traitIDs); err != nil {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if !(token.Id > seenID) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("unsorted token"), "index %d", i)
		}
		seenID = token.Id
	}

	return nil
}

type Traits []*Trait

func (ts Traits) ValidateBasic() error {
	seenID := ""
	for i, trait := range ts {
		if trait == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("nil trait"), "index %d", i)
		}

		if err := (&TraitInternal{}).Parse(*trait); err != nil {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if !(trait.Id > seenID) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("unsorted trait"), "index %d", i)
		}
		seenID = trait.Id
	}

	return nil
}

func (c GenesisClass) ValidateCompatibility() error {
	if c.Id == "" {
		return ErrUnimplemented.Wrap("nil id")
	}

	if c.Traits == nil {
		return ErrUnimplemented.Wrap("nil traits")
	}

	if c.Tokens == nil {
		return ErrUnimplemented.Wrap("nil tokens")
	}

	return nil
}

func (c GenesisClass) ValidateBasic() error {
	if err := c.ValidateCompatibility(); err != nil {
		return err
	}

	var id Reference
	if err := id.Parse(c.Id); err != nil {
		return err
	}

	if err := Traits(c.Traits).ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "traits")
	}

	traitIDs := map[string]struct{}{}
	for _, trait := range c.Traits {
		traitIDs[trait.Id] = struct{}{}
	}

	if err := GenesisTokens(c.Tokens).ValidateBasic(traitIDs); err != nil {
		return errorsmod.Wrap(err, "tokens")
	}

	return nil
}

type GenesisClasses []*GenesisClass

func (cs GenesisClasses) ValidateBasic() error {
	seenID := ""
	for i, class := range cs {
		if class == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("nil class"), "index %d", i)
		}

		if err := class.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if !(class.Id > seenID) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("unsorted class"), "index %d", i)
		}
		seenID = class.Id
	}

	return nil
}

func (s GenesisState) ValidateCompatibility() error {
	if s.Params == nil {
		return ErrUnimplemented.Wrap("nil params")
	}

	if s.Classes == nil {
		return ErrUnimplemented.Wrap("nil classes")
	}

	return nil
}

// ValidateBasic check the given genesis state has no integrity issues
func (s GenesisState) ValidateBasic() error {
	if err := s.ValidateCompatibility(); err != nil {
		return err
	}

	if err := (&ParamsInternal{}).Parse(*s.Params); err != nil {
		return errorsmod.Wrap(err, "params")
	}

	if err := GenesisClasses(s.Classes).ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "classes")
	}

	return nil
}
