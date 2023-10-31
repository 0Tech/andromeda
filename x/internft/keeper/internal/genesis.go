package internal

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) InitGenesis(ctx context.Context, gs *internftv1alpha1.GenesisState) error {
	k.SetParams(ctx, gs.Params)

	for _, genClass := range gs.Classes {
		class := internftv1alpha1.Class{
			Id: genClass.Id,
		}
		k.setClass(ctx, class)

		for _, trait := range genClass.Traits {
			k.setTrait(ctx, class.Id, trait)
		}

		for _, genNFT := range genClass.Nfts {
			nft := internftv1alpha1.NFT{
				ClassId: class.Id,
				Id:      genNFT.Id,
			}
			k.setNFT(ctx, nft)

			for _, property := range genNFT.Properties {
				k.setProperty(ctx, nft, property)
			}

			owner := genNFT.Owner
			k.setOwner(ctx, nft, sdk.MustAccAddressFromBech32(owner))
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) *internftv1alpha1.GenesisState {
	classes := k.getClasses(ctx)

	var genClasses []internftv1alpha1.GenesisClass
	if len(classes) != 0 {
		genClasses = make([]internftv1alpha1.GenesisClass, len(classes))
	}

	for classIndex, class := range classes {
		genClasses[classIndex].Id = class.Id

		genClasses[classIndex].Traits = k.getTraitsOfClass(ctx, class.Id)

		nfts := k.getNFTsOfClass(ctx, class.Id)

		var genNFTs []internftv1alpha1.GenesisNFT
		if len(nfts) != 0 {
			genNFTs = make([]internftv1alpha1.GenesisNFT, len(nfts))
		}

		for nftIndex, nft := range nfts {
			genNFTs[nftIndex].Id = nft.Id

			genNFTs[nftIndex].Properties = k.getPropertiesOfNFT(ctx, nft)

			owner, err := k.getOwner(ctx, nft)
			if err != nil {
				panic(err)
			}
			genNFTs[nftIndex].Owner = owner.String()
		}

		genClasses[classIndex].Nfts = genNFTs
	}

	return &internftv1alpha1.GenesisState{
		Params:  k.GetParams(ctx),
		Classes: genClasses,
	}
}

func (k Keeper) getClasses(ctx context.Context) (classes []internftv1alpha1.Class) {
	k.iterateClasses(ctx, func(class internftv1alpha1.Class) {
		classes = append(classes, class)
	})

	return
}

func (k Keeper) getTraitsOfClass(ctx context.Context, classID string) (traits []internftv1alpha1.Trait) {
	k.iterateTraitsOfClass(ctx, classID, func(trait internftv1alpha1.Trait) {
		traits = append(traits, trait)
	})

	return
}

func (k Keeper) getNFTsOfClass(ctx context.Context, classID string) (nfts []internftv1alpha1.NFT) {
	k.iterateNFTsOfClass(ctx, classID, func(nft internftv1alpha1.NFT) {
		nfts = append(nfts, nft)
	})

	return
}

func (k Keeper) getPropertiesOfNFT(ctx context.Context, nft internftv1alpha1.NFT) (properties []internftv1alpha1.Property) {
	k.iteratePropertiesOfClass(ctx, nft, func(property internftv1alpha1.Property) {
		properties = append(properties, property)
	})

	return
}
