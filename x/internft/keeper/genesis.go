package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internft "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *internft.GenesisState) error {
	k.SetParams(ctx, gs.Params)

	for _, genClass := range gs.Classes {
		class := internft.Class{
			Id: genClass.Id,
		}
		k.setClass(ctx, class)

		for _, trait := range genClass.Traits {
			k.setTrait(ctx, class.Id, trait)
		}

		k.setPreviousID(ctx, class.Id, genClass.LastMintedNftId)

		for _, genNFT := range genClass.Nfts {
			nft := internft.NFT{
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

func (k Keeper) ExportGenesis(ctx sdk.Context) *internft.GenesisState {
	classes := k.getClasses(ctx)

	var genClasses []internft.GenesisClass
	if len(classes) != 0 {
		genClasses = make([]internft.GenesisClass, len(classes))
	}

	for classIndex, class := range classes {
		genClasses[classIndex].Id = class.Id
		genClasses[classIndex].LastMintedNftId = k.GetPreviousID(ctx, class.Id)

		genClasses[classIndex].Traits = k.getTraitsOfClass(ctx, class.Id)

		nfts := k.getNFTsOfClass(ctx, class.Id)

		var genNFTs []internft.GenesisNFT
		if len(nfts) != 0 {
			genNFTs = make([]internft.GenesisNFT, len(nfts))
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

	return &internft.GenesisState{
		Params:  k.GetParams(ctx),
		Classes: genClasses,
	}
}

func (k Keeper) getClasses(ctx sdk.Context) (classes []internft.Class) {
	k.iterateClasses(ctx, func(class internft.Class) {
		classes = append(classes, class)
	})

	return
}

func (k Keeper) getTraitsOfClass(ctx sdk.Context, classID string) (traits []internft.Trait) {
	k.iterateTraitsOfClass(ctx, classID, func(trait internft.Trait) {
		traits = append(traits, trait)
	})

	return
}

func (k Keeper) getNFTsOfClass(ctx sdk.Context, classID string) (nfts []internft.NFT) {
	k.iterateNFTsOfClass(ctx, classID, func(nft internft.NFT) {
		nfts = append(nfts, nft)
	})

	return
}

func (k Keeper) getPropertiesOfNFT(ctx sdk.Context, nft internft.NFT) (properties []internft.Property) {
	k.iteratePropertiesOfClass(ctx, nft, func(property internft.Property) {
		properties = append(properties, property)
	})

	return
}
