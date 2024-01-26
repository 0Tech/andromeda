package module

import (
	"context"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic defines the basic application module used by the module.
type AppModuleBasic struct{}

// ____________________________________________________________________________

var _ module.AppModuleBasic = (*AppModuleBasic)(nil)

// Name returns the name of the module.
func (AppModuleBasic) Name() string {
	return testv1alpha1.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
	// Amino deprecated.
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	testv1alpha1.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := testv1alpha1.RegisterQueryHandlerClient(context.Background(), mux, testv1alpha1.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}
