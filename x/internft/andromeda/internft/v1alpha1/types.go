package internftv1alpha1

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateAddress(address string) error {
	if _, err := sdk.AccAddressFromBech32(address); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap(address)
	}

	return nil
}

func ValidateClassID(id string) error {
	if _, err := sdk.AccAddressFromBech32(id); err != nil {
		return ErrInvalidClassID.Wrap(id)
	}

	return nil
}

func (p Params) ValidateCompatibility() error {
	return nil
}

func (p Params) ValidateBasic() error {
	if err := p.ValidateCompatibility(); err != nil {
		return err
	}

	return nil
}

func (c Class) ValidateCompatibility() error {
	if c.Id == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil id")
	}

	return nil
}

func (c Class) ValidateBasic() error {
	if err := c.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(c.Id); err != nil {
		return err
	}

	return nil
}

func ValidateTraitID(id string) error {
	return nil
}

func (t Trait) ValidateCompatibility() error {
	if t.Id == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil id")
	}

	if t.Mutability == Trait_MUTABILITY_UNSPECIFIED {
		return sdkerrors.ErrNotSupported.Wrap("nil mutability")
	}

	return nil
}

func (t Trait) ValidateBasic() error {
	if err := t.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateTraitID(t.Id); err != nil {
		return err
	}

	return nil
}

type Traits []*Trait

func (ts Traits) ValidateBasic() error {
	seenIDs := map[string]struct{}{}
	for i, trait := range ts {
		if trait == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("nil trait"), "index %d", i)
		}

		if err := trait.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if _, seen := seenIDs[trait.Id]; seen {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("duplicate id"), "index %d", i)
		}
		seenIDs[trait.Id] = struct{}{}
	}

	return nil
}

func ValidateTokenID(id string) error {
	if _, err := sdk.AccAddressFromBech32(id); err != nil {
		return ErrInvalidTokenID.Wrap(id)
	}

	return nil
}

func (t Token) ValidateCompatibility() error {
	if t.ClassId == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil class id")
	}

	if t.Id == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil id")
	}

	return nil
}

func (t Token) ValidateBasic() error {
	if err := t.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(t.ClassId); err != nil {
		return err
	}

	if err := ValidateTokenID(t.Id); err != nil {
		return err
	}

	return nil
}

func (p Property) ValidateCompatibility() error {
	if p.TraitId == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil trait id")
	}

	if p.Fact == "" {
		return sdkerrors.ErrNotSupported.Wrap("nil fact")
	}

	return nil
}

func (p Property) ValidateBasic() error {
	if err := p.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateTraitID(p.TraitId); err != nil {
		return err
	}

	return nil
}

type Properties []*Property

func (ps Properties) ValidateBasic() error {
	seenTraitIDs := map[string]struct{}{}
	for i, property := range ps {
		if property == nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("nil property"), "index %d", i)
		}

		if err := property.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(err, "index %d", i)
		}

		if _, seen := seenTraitIDs[property.TraitId]; seen {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest.Wrap("duplicate trait id"), "index %d", i)
		}
		seenTraitIDs[property.TraitId] = struct{}{}
	}

	return nil
}

func ValidateOperator(operator, classID string) error {
	if operator != classID {
		return sdkerrors.ErrUnauthorized.Wrapf("%s over class %s", operator, classID)
	}

	return nil
}
