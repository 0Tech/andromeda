package module

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/api/andromeda/internft/v1alpha1"
)

func (AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: internftv1alpha1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use: "params",
					Short: "Query the module parameters",
				},
				{
					RpcMethod: "Class",
					Use: "class",
					Short: "Query a class",
				},
				{
					RpcMethod: "Classes",
					Use: "classes",
					Short: "Query all the classes",
				},
				{
					RpcMethod: "Trait",
					Use: "trait",
					Short: "Query a trait of a class",
				},
				{
					RpcMethod: "Traits",
					Use: "traits",
					Short: "Query all the traits of a class",
				},
				{
					RpcMethod: "NFT",
					Use: "nft",
					Short: "Query an NFT",
				},
				{
					RpcMethod: "NFTs",
					Use: "nfts",
					Short: "Query all the NFTs",
				},
				{
					RpcMethod: "Property",
					Use: "property",
					Short: "Query a property of an NFT",
				},
				{
					RpcMethod: "Properties",
					Use: "properties",
					Short: "Query all the properties of an NFT",
				},
				{
					RpcMethod: "Owner",
					Use: "owner",
					Short: "Query the owner of an NFT",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: internftv1alpha1.Msg_ServiceDesc.ServiceName,
		},
	}
}
