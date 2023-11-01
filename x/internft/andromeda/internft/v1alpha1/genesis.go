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

		seenNFTIDs := map[string]struct{}{}
		for nftIndex, genNFT := range genClass.Nfts {
			errHint := fmt.Sprintf("%s.nfts[%d]", errHint, nftIndex)

			id := genNFT.Id
			if err := ValidateNFTID(id); err != nil {
				return errorsmod.Wrap(err, errHint)
			}

			if _, seen := seenNFTIDs[id]; seen {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest.Wrap("unsorted nfts"), errHint)
			}
			seenNFTIDs[id] = struct{}{}

			if err := Properties(genNFT.Properties).ValidateBasic(); err != nil {
				return errorsmod.Wrap(err, errHint)
			}

			for _, property := range genNFT.Properties {
				if _, hasTrait := traits[property.Id]; !hasTrait {
					return errorsmod.Wrap(ErrTraitNotFound.Wrap(property.Id), errHint)
				}
			}

			if err := ValidateAddress(genNFT.Owner); err != nil {
				return errorsmod.Wrap(errorsmod.Wrap(err, "owner"), errHint)
			}
		}
	}

	return nil
}
