package internal

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) CreateClass(ctx context.Context, class *internftv1alpha1.Class) error {
	if err := k.hasClass(ctx, class.Id); err == nil {
		return internftv1alpha1.ErrClassAlreadyExists.Wrap(class.Id)
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

func (k Keeper) setClass(ctx context.Context, class *internftv1alpha1.Class) {
	if err := k.classes.Set(ctx, class.Id, *class); err != nil {
		panic(err)
	}
}

func (k Keeper) UpdateTrait(ctx context.Context, class *internftv1alpha1.Class, trait *internftv1alpha1.Trait) error {
	if err := k.hasClass(ctx, class.Id); err != nil {
		return err
	}

	// mutating existing trait
	if prevTrait, err := k.GetTrait(ctx, class.Id, trait.Id); err == nil {
		switch prevTrait.Mutability {
		case internftv1alpha1.Trait_MUTABILITY_IMMUTABLE:
			return internftv1alpha1.ErrTraitImmutable.Wrap(trait.Id)
		}
	}

	k.setTrait(ctx, class.Id, trait)

	return nil
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

func (k Keeper) setTrait(ctx context.Context, classID string, trait *internftv1alpha1.Trait) {
	if err := k.traits.Set(ctx, collections.Join(classID, trait.Id), *trait); err != nil {
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

func (k Keeper) MintToken(ctx context.Context, owner sdk.AccAddress, token *internftv1alpha1.Token) error {
	if err := k.hasClass(ctx, token.ClassId); err != nil {
		return err
	}

	if err := k.hasToken(ctx, token); err == nil {
		return internftv1alpha1.ErrTokenAlreadyExists.Wrap(token.String())
	}
	k.setToken(ctx, token)

	k.setOwner(ctx, token, owner)

	return nil
}

// func (k Keeper) hasProperty(ctx context.Context, token internftv1alpha1.Token, propertyID string) error {
// 	_, err := k.GetProperty(ctx, token, propertyID)
// 	return err
// }

func (k Keeper) GetProperty(ctx context.Context, token *internftv1alpha1.Token, propertyID string) (*internftv1alpha1.Property, error) {
	property, err := k.properties.Get(ctx, collections.Join3(token.ClassId, token.Id, propertyID))
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			err = internftv1alpha1.ErrTraitNotFound.Wrapf("%s, %s", token.ClassId, propertyID)
		}

		return nil, err
	}

	return &property, nil
}

func (k Keeper) setProperty(ctx context.Context, token *internftv1alpha1.Token, property *internftv1alpha1.Property) {
	if err := k.properties.Set(ctx, collections.Join3(token.ClassId, token.Id, property.TraitId), *property); err != nil {
		panic(err)
	}
}

func (k Keeper) iteratePropertiesOfToken(ctx context.Context, token *internftv1alpha1.Token, fn func(property internftv1alpha1.Property)) {
	rng := collections.NewSuperPrefixedTripleRange[string, string, string](token.ClassId, token.Id)
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

func (k Keeper) BurnToken(ctx context.Context, owner sdk.AccAddress, token *internftv1alpha1.Token) error {
	if err := k.validateOwner(ctx, token, owner); err != nil {
		return err
	}
	k.deleteOwner(ctx, token)

	if err := k.hasToken(ctx, token); err != nil {
		panic(err)
	}
	k.deleteToken(ctx, token)

	k.properties.Clear(ctx, collections.NewSuperPrefixedTripleRange[string, string, string](token.ClassId, token.Id))

	return nil
}

func (k Keeper) UpdateProperty(ctx context.Context, token *internftv1alpha1.Token, property *internftv1alpha1.Property) error {
	if err := k.hasToken(ctx, token); err != nil {
		return err
	}

	trait, err := k.GetTrait(ctx, token.ClassId, property.TraitId)
	if err != nil {
		return err
	}

	// mutating existing property
	if _, err := k.GetProperty(ctx, token, property.TraitId); err == nil {
		switch trait.Mutability {
		case internftv1alpha1.Trait_MUTABILITY_IMMUTABLE:
			return internftv1alpha1.ErrTraitImmutable.Wrap(property.TraitId)
		}
	}
		
	k.setProperty(ctx, token, property)

	return nil
}

func (k Keeper) hasToken(ctx context.Context, token *internftv1alpha1.Token) error {
	_, err := k.GetToken(ctx, token)
	return err
}

func (k Keeper) GetToken(ctx context.Context, token *internftv1alpha1.Token) (*internftv1alpha1.Token, error) {
	_, err := k.tokens.Get(ctx, collections.Join(token.ClassId, token.Id))
	if err != nil {
		if errorsmod.IsOf(err, collections.ErrNotFound) {
			err = internftv1alpha1.ErrTokenNotFound.Wrap(token.String())
		}

		return nil, err
	}

	return token, nil
}

func (k Keeper) setToken(ctx context.Context, token *internftv1alpha1.Token) {
	if err := k.tokens.Set(ctx, collections.Join(token.ClassId, token.Id), *token); err != nil {
		panic(err)
	}
}

func (k Keeper) deleteToken(ctx context.Context, token *internftv1alpha1.Token) {
	if err := k.tokens.Remove(ctx, collections.Join(token.ClassId, token.Id)); err != nil {
		panic(err)
	}
}

func (k Keeper) iterateTokensOfClass(ctx context.Context, classID string, fn func(token internftv1alpha1.Token)) {
	rng := collections.NewPrefixedPairRange[string, string](classID)
	iter, err := k.tokens.Iterate(ctx, rng)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		token, err := iter.Value()
		if err != nil {
			panic(err)
		}

		fn(token)
	}
}
