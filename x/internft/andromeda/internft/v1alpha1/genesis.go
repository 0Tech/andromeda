package internftv1alpha1

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

func DefaultParams() Params {
	return Params{}
}

// ValidateBasic check the given genesis state has no integrity issues
func (s GenesisState) ValidateBasic() error {
	seenClassIDs := map[string]struct{}{}
	for classIndex, genClass := range s.Classes {
		errHint := fmt.Sprintf("classes[%d]", classIndex)

		id := genClass.Id
		if err := ValidateClassID(id); err != nil {
			return errorsmod.Wrap(err, errHint)
		}

		if _, seen := seenClassIDs[id]; seen {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest.Wrapf("duplicate class id %s", genClass.Id), errHint)
		}
		seenClassIDs[id] = struct{}{}

		if err := Traits(genClass.Traits).ValidateBasic(); err != nil {
			return errorsmod.Wrap(err, errHint)
		}

		traits := map[string]struct{}{}
		for _, trait := range genClass.Traits {
			traits[trait.Id] = struct{}{}
		}

		seenTokenIDs := map[string]struct{}{}
		for tokenIndex, genToken := range genClass.Tokens {
			errHint := fmt.Sprintf("%s.tokens[%d]", errHint, tokenIndex)

			id := genToken.Id
			if err := ValidateTokenID(id); err != nil {
				return errorsmod.Wrap(err, errHint)
			}

			if _, seen := seenTokenIDs[id]; seen {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest.Wrap("unsorted tokens"), errHint)
			}
			seenTokenIDs[id] = struct{}{}

			if err := Properties(genToken.Properties).ValidateBasic(); err != nil {
				return errorsmod.Wrap(err, errHint)
			}

			for _, property := range genToken.Properties {
				if _, hasTrait := traits[property.TraitId]; !hasTrait {
					return errorsmod.Wrap(ErrTraitNotFound.Wrap(property.TraitId), errHint)
				}
			}

			if err := ValidateAddress(genToken.Owner); err != nil {
				return errorsmod.Wrap(errorsmod.Wrap(err, "owner"), errHint)
			}
		}
	}

	return nil
}
