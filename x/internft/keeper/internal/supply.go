package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) NewClass(ctx context.Context, class internftv1alpha1.Class, traits []internftv1alpha1.Trait) error {
	if err := k.hasClass(ctx, class.Id); err == nil {
		return internftv1alpha1.ErrClassAlreadyExists.Wrap(class.Id)
	}
	k.setClass(ctx, class)

	for _, trait := range traits {
		k.setTrait(ctx, class.Id, trait)
	}

	return nil
}

func (k Keeper) UpdateClass(ctx context.Context, class internftv1alpha1.Class) error {
	if err := k.hasClass(ctx, class.Id); err != nil {
		return err
	}
	k.setClass(ctx, class)

	return nil
}

func (k Keeper) hasClass(ctx context.Context, classID string) error {
	_, err := k.GetClass(ctx, classID)
	return err
}

func (k Keeper) GetClass(ctx context.Context, classID string) (*internftv1alpha1.Class, error) {
	class, err := k.classes.Get(ctx, classID)
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			err = internftv1alpha1.ErrClassNotFound.Wrap(classID)
		}

		return nil, err
	}

	return &class, nil
}

func (k Keeper) setClass(ctx context.Context, class internftv1alpha1.Class) {
	if err := k.classes.Set(ctx, class.Id, class); err != nil {
		panic(err)
	}
}

func (k Keeper) hasTrait(ctx context.Context, classID string, traitID string) error {
	_, err := k.GetTrait(ctx, classID, traitID)
	return err
}

func (k Keeper) GetTrait(ctx context.Context, classID string, traitID string) (*internftv1alpha1.Trait, error) {
	trait, err := k.traits.Get(ctx, collections.Join(classID, traitID))
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			err = internftv1alpha1.ErrTraitNotFound.Wrapf("%s in %s", traitID, classID)
		}

		return nil, err
	}

	return &trait, nil
}

func (k Keeper) setTrait(ctx context.Context, classID string, trait internftv1alpha1.Trait) {
	if err := k.traits.Set(ctx, collections.Join(classID, trait.Id), trait); err != nil {
		panic(err)
	}
}

func (k Keeper) iterateTraitsOfClass(ctx context.Context, classID string, fn func(trait internftv1alpha1.Trait)) {
	rng := collections.NewPrefixedPairRange[string, string](classID)
	iter, err := k.traits.Iterate(ctx, rng)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		trait, err := iter.Value()
		if err != nil {
			panic(err)
		}

		fn(trait)
	}
}

func (k Keeper) iterateClasses(ctx context.Context, fn func(class internftv1alpha1.Class)) {
	iter, err := k.classes.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		class, err := iter.Value()
		if err != nil {
			panic(err)
		}

		fn(class)
	}
}

func (k Keeper) MintNFT(ctx context.Context, owner sdk.AccAddress, nft internftv1alpha1.NFT, properties []internftv1alpha1.Property) error {
	if err := k.hasClass(ctx, nft.ClassId); err != nil {
		return err
	}

	if err := k.hasNFT(ctx, nft); err == nil {
		// TODO(@0Tech): define the error
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest.Wrap("nft already exists"), nft.String())
	}
	k.setNFT(ctx, nft)

	for _, property := range properties {
		if err := k.hasTrait(ctx, nft.ClassId, property.Id); err != nil {
			return errorsmod.Wrap(err, property.Id)
		}

		k.setProperty(ctx, nft, property)
	}

	k.setOwner(ctx, nft, owner)

	return nil
}

// func (k Keeper) hasProperty(ctx context.Context, nft internftv1alpha1.NFT, propertyID string) error {
// 	_, err := k.GetProperty(ctx, nft, propertyID)
// 	return err
// }

func (k Keeper) GetProperty(ctx context.Context, nft internftv1alpha1.NFT, propertyID string) (*internftv1alpha1.Property, error) {
	property, err := k.properties.Get(ctx, collections.Join3(nft.ClassId, nft.Id, propertyID))
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			err = internftv1alpha1.ErrTraitNotFound.Wrapf("%s, %s", nft.ClassId, propertyID)
		}

		return nil, err
	}

	return &property, nil
}

func (k Keeper) setProperty(ctx context.Context, nft internftv1alpha1.NFT, property internftv1alpha1.Property) {
	if err := k.properties.Set(ctx, collections.Join3(nft.ClassId, nft.Id, property.Id), property); err != nil {
		panic(err)
	}
}

func (k Keeper) iteratePropertiesOfClass(ctx context.Context, nft internftv1alpha1.NFT, fn func(property internftv1alpha1.Property)) {
	rng := collections.NewSuperPrefixedTripleRange[string, string, string](nft.ClassId, nft.Id)
	iter, err := k.properties.Iterate(ctx, rng)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		property, err := iter.Value()
		if err != nil {
			panic(err)
		}

		fn(property)
	}
}

func (k Keeper) BurnNFT(ctx context.Context, owner sdk.AccAddress, nft internftv1alpha1.NFT) error {
	if err := k.validateOwner(ctx, nft, owner); err != nil {
		return err
	}
	k.deleteOwner(ctx, nft)

	if err := k.hasNFT(ctx, nft); err != nil {
		panic(err)
	}
	k.deleteNFT(ctx, nft)

	return nil
}

func (k Keeper) UpdateNFT(ctx context.Context, nft internftv1alpha1.NFT, properties []internftv1alpha1.Property) error {
	if err := k.hasNFT(ctx, nft); err != nil {
		return err
	}

	for _, property := range properties {
		trait, err := k.GetTrait(ctx, nft.ClassId, property.Id)
		if err != nil {
			return err
		}

		if !trait.Variable {
			return internftv1alpha1.ErrTraitImmutable.Wrap(property.Id)
		}

		k.setProperty(ctx, nft, property)
	}

	return nil
}

func (k Keeper) hasNFT(ctx context.Context, nft internftv1alpha1.NFT) error {
	_, err := k.GetNFT(ctx, nft)
	return err
}

func (k Keeper) GetNFT(ctx context.Context, nft internftv1alpha1.NFT) (*internftv1alpha1.NFT, error) {
	nft, err := k.nfts.Get(ctx, collections.Join(nft.ClassId, nft.Id))
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			err = internftv1alpha1.ErrNFTNotFound.Wrap(nft.String())
		}

		return nil, err
	}

	return &nft, nil
}

func (k Keeper) setNFT(ctx context.Context, nft internftv1alpha1.NFT) {
	if err := k.nfts.Set(ctx, collections.Join(nft.ClassId, nft.Id), nft); err != nil {
		panic(err)
	}
}

func (k Keeper) deleteNFT(ctx context.Context, nft internftv1alpha1.NFT) {
	if err := k.nfts.Remove(ctx, collections.Join(nft.ClassId, nft.Id)); err != nil {
		panic(err)
	}
}

func (k Keeper) iterateNFTsOfClass(ctx context.Context, classID string, fn func(nft internftv1alpha1.NFT)) {
	rng := collections.NewPrefixedPairRange[string, string](classID)
	iter, err := k.nfts.Iterate(ctx, rng)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		nft, err := iter.Value()
		if err != nil {
			panic(err)
		}

		fn(nft)
	}
}
