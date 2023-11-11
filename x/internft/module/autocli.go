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
					RpcMethod: "Token",
					Use: "token",
					Short: "Query a token",
				},
				{
					RpcMethod: "Tokens",
					Use: "tokens",
					Short: "Query all the tokens",
				},
				{
					RpcMethod: "Property",
					Use: "property",
					Short: "Query a property of a token",
				},
				{
					RpcMethod: "Properties",
					Use: "properties",
					Short: "Query all the properties of a token",
				},
				{
					RpcMethod: "Owner",
					Use: "owner",
					Short: "Query the owner of a token",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: internftv1alpha1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// {
				// 	RpcMethod: "UpdateParams",
				// 	Use: "update-params",
				// 	Short: "Update the module parameters",
				// },
				{
					RpcMethod: "Send",
					Use: "send",
					Short: "Send a token",
				},
				{
					RpcMethod: "NewClass",
					Use: "new-class",
					Short: "Create a new class",
				},
				{
					RpcMethod: "UpdateClass",
					Use: "update-class",
					Short: "Update a class",
				},
				{
					RpcMethod: "NewToken",
					Use: "new-token",
					Short: "Create a new token",
				},
				{
					RpcMethod: "BurnToken",
					Use: "burn-token",
					Short: "Burn a token",
				},
				{
					RpcMethod: "UpdateToken",
					Use: "update-token",
					Short: "Update properties of a token",
				},
			},
		},
	}
}
