package internal

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) InitGenesis(ctx context.Context, gs *internftv1alpha1.GenesisState) error {
	k.SetParams(ctx, *gs.Params)

	for _, genClass := range gs.Classes {
		class := &internftv1alpha1.Class{
			Id: genClass.Id,
		}
		k.setClass(ctx, class)

		for _, trait := range genClass.Traits {
			k.setTrait(ctx, class.Id, trait)
		}

		for _, genToken := range genClass.Tokens {
			token := &internftv1alpha1.Token{
				ClassId: class.Id,
				Id:      genToken.Id,
			}
			k.setToken(ctx, token)

			for _, property := range genToken.Properties {
				k.setProperty(ctx, token, property)
			}

			owner := genToken.Owner
			k.setOwner(ctx, token, sdk.MustAccAddressFromBech32(owner))
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) *internftv1alpha1.GenesisState {
	genesis := &internftv1alpha1.GenesisState{}
	genesis.Params = k.GetParams(ctx)
	
	classes := k.getClasses(ctx)
	genClasses := make([]*internftv1alpha1.GenesisClass, len(classes))
	for classIndex, class := range classes {
		genClasses[classIndex] = &internftv1alpha1.GenesisClass{}
		genClasses[classIndex].Id = class.Id
		genClasses[classIndex].Traits = k.getTraitsOfClass(ctx, class.Id)

		tokens := k.getTokensOfClass(ctx, class.Id)
		genTokens := make([]*internftv1alpha1.GenesisToken, len(tokens))
		for tokenIndex, token := range tokens {
			genTokens[tokenIndex] = &internftv1alpha1.GenesisToken{}
			genTokens[tokenIndex].Id = token.Id
			genTokens[tokenIndex].Properties = k.getPropertiesOfToken(ctx, token)

			owner, err := k.getOwner(ctx, token)
			if err != nil {
				panic(err)
			}
			genTokens[tokenIndex].Owner = owner.String()
		}
		genClasses[classIndex].Tokens = genTokens
	}
	genesis.Classes = genClasses

	return genesis
}

func (k Keeper) getClasses(ctx context.Context) (classes []internftv1alpha1.Class) {
	k.iterateClasses(ctx, func(class internftv1alpha1.Class) {
		classes = append(classes, class)
	})

	return
}

func (k Keeper) getTraitsOfClass(ctx context.Context, classID string) (traits []*internftv1alpha1.Trait) {
	traits = []*internftv1alpha1.Trait{}
	k.iterateTraitsOfClass(ctx, classID, func(trait internftv1alpha1.Trait) {
		traits = append(traits, &trait)
	})

	return
}

func (k Keeper) getTokensOfClass(ctx context.Context, classID string) (tokens []*internftv1alpha1.Token) {
	tokens = []*internftv1alpha1.Token{}
	k.iterateTokensOfClass(ctx, classID, func(token internftv1alpha1.Token) {
		tokens = append(tokens, &token)
	})

	return
}

func (k Keeper) getPropertiesOfToken(ctx context.Context, token *internftv1alpha1.Token) (properties []*internftv1alpha1.Property) {
	properties = []*internftv1alpha1.Property{}
	k.iteratePropertiesOfToken(ctx, token, func(property internftv1alpha1.Property) {
		properties = append(properties, &property)
	})

	return
}
