package internal

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) InitGenesis(ctx context.Context, gs *internftv1alpha1.GenesisState) error {
	if err := k.setParams(ctx, gs.Params); err != nil {
		return err
	}

	for _, genClass := range gs.Classes {
		class := &internftv1alpha1.Class{
			Id: genClass.Id,
		}
		if err := k.setClass(ctx, class); err != nil {
			return err
		}

		for _, trait := range genClass.Traits {
			if err := k.setTrait(ctx, class.Id, trait); err != nil {
				return err
			}
		}

		for _, genToken := range genClass.Tokens {
			token := &internftv1alpha1.Token{
				ClassId: class.Id,
				Id:      genToken.Id,
			}
			if err := k.setToken(ctx, token); err != nil {
				return err
			}

			for _, property := range genToken.Properties {
				if err := k.setProperty(ctx, token, property); err != nil {
					return err
				}
			}

			owner := genToken.Owner
			if err := k.setOwner(ctx, token, sdk.MustAccAddressFromBech32(owner)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) (*internftv1alpha1.GenesisState, error) {
	genesis := &internftv1alpha1.GenesisState{}

	var err error
	genesis.Params, err = k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	
	classes, err := k.getClasses(ctx)
	if err != nil {
		return nil, err
	}
	genClasses := make([]*internftv1alpha1.GenesisClass, len(classes))
	for classIndex, class := range classes {
		genClasses[classIndex] = &internftv1alpha1.GenesisClass{}
		genClasses[classIndex].Id = class.Id
		genClasses[classIndex].Traits, err = k.getTraitsOfClass(ctx, class.Id)
		if err != nil {
			return nil, err
		}

		tokens, err := k.getTokensOfClass(ctx, class.Id)
		if err != nil {
			return nil, err
		}
		genTokens := make([]*internftv1alpha1.GenesisToken, len(tokens))
		for tokenIndex, token := range tokens {
			genTokens[tokenIndex] = &internftv1alpha1.GenesisToken{}
			genTokens[tokenIndex].Id = token.Id
			genTokens[tokenIndex].Properties, err = k.getPropertiesOfToken(ctx, token)
			if err != nil {
				return nil, err
			}

			owner, err := k.getOwner(ctx, token)
			if err != nil {
				return nil, err
			}
			genTokens[tokenIndex].Owner = owner.String()
		}
		genClasses[classIndex].Tokens = genTokens
	}
	genesis.Classes = genClasses

	return genesis, nil
}

func (k Keeper) getClasses(ctx context.Context) ([]*internftv1alpha1.Class, error) {
	classes := []*internftv1alpha1.Class{}
	if err := k.iterateClasses(ctx, func(class internftv1alpha1.Class) {
		classes = append(classes, &class)
	}); err != nil {
		return nil, err
	}

	return classes, nil
}

func (k Keeper) getTraitsOfClass(ctx context.Context, classID string) ([]*internftv1alpha1.Trait, error) {
	traits := []*internftv1alpha1.Trait{}
	if err := k.iterateTraitsOfClass(ctx, classID, func(trait internftv1alpha1.Trait) {
		traits = append(traits, &trait)
	}); err != nil {
		return nil, err
	}

	return traits, nil
}

func (k Keeper) getTokensOfClass(ctx context.Context, classID string) ([]*internftv1alpha1.Token, error) {
	tokens := []*internftv1alpha1.Token{}
	if err := k.iterateTokensOfClass(ctx, classID, func(token internftv1alpha1.Token) {
		tokens = append(tokens, &token)
	}); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (k Keeper) getPropertiesOfToken(ctx context.Context, token *internftv1alpha1.Token) ([]*internftv1alpha1.Property, error) {
	properties := []*internftv1alpha1.Property{}
	if err := k.iteratePropertiesOfToken(ctx, token, func(property internftv1alpha1.Property) {
		properties = append(properties, &property)
	}); err != nil {
		return nil, err
	}

	return properties, nil
}
