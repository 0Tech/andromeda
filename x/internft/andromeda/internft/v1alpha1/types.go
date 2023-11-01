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

func (class Class) ValidateBasic() error {
	if err := ValidateClassID(class.Id); err != nil {
		return err
	}

	return nil
}

func ValidateTraitID(id string) error {
	if len(id) == 0 {
		return ErrInvalidTraitID.Wrap("empty")
	}

	return nil
}

func (t Trait) ValidateBasic() error {
	if err := ValidateTraitID(t.Id); err != nil {
		return err
	}

	return nil
}

type Traits []Trait

func (t Traits) ValidateBasic() error {
	seenIDs := map[string]struct{}{}
	for _, trait := range t {
		if err := trait.ValidateBasic(); err != nil {
			return errorsmod.Wrap(err, trait.Id)
		}

		if _, seen := seenIDs[trait.Id]; seen {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest.Wrap("duplicate id"), trait.Id)
		}
		seenIDs[trait.Id] = struct{}{}
	}

	return nil
}

func ValidateNFTID(id string) error {
	if _, err := sdk.AccAddressFromBech32(id); err != nil {
		return ErrInvalidNFTID.Wrap(id)
	}

	return nil
}

func (nft NFT) ValidateBasic() error {
	if err := ValidateClassID(nft.ClassId); err != nil {
		return err
	}

	if err := ValidateNFTID(nft.Id); err != nil {
		return err
	}

	return nil
}

func (p Property) ValidateBasic() error {
	if len(p.Id) == 0 {
		return ErrInvalidTraitID.Wrap("empty")
	}

	return nil
}

type Properties []Property

func (p Properties) ValidateBasic() error {
	seenIDs := map[string]struct{}{}
	for _, property := range p {
		if err := property.ValidateBasic(); err != nil {
			return errorsmod.Wrap(err, property.Id)
		}

		if _, seen := seenIDs[property.Id]; seen {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest.Wrap("duplicate id"), property.Id)
		}
		seenIDs[property.Id] = struct{}{}
	}

	return nil
}

func ValidateOperator(operator, classID string) error {
	if operator != classID {
		return sdkerrors.ErrUnauthorized.Wrapf("%s over class %s", operator, classID)
	}

	return nil
}
