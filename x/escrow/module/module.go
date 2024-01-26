package module

import (
	"context"
	"encoding/json"
	"fmt"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/keeper"
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
	return escrowv1alpha1.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
	// Amino deprecated.
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	escrowv1alpha1.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := escrowv1alpha1.RegisterQueryHandlerClient(context.Background(), mux, escrowv1alpha1.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements an application module for the module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

// ____________________________________________________________________________

var _ module.HasGenesisBasics = (*AppModule)(nil)

// DefaultGenesis returns default genesis state as raw bytes for the module.
func (am AppModule) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(am.keeper.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the module.
func (am AppModule) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var gs escrowv1alpha1.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", escrowv1alpha1.ModuleName, err)
	}

	return am.keeper.ValidateGenesis(&gs)
}

// ____________________________________________________________________________

var _ module.HasInvariants = (*AppModule)(nil)

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
	// No invariants.
}

// ____________________________________________________________________________

var _ module.HasServices = (*AppModule)(nil)

func (am AppModule) RegisterServices(cfg module.Configurator) {
	escrowv1alpha1.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
	escrowv1alpha1.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))

	// m := keeper.NewMigrator(am.keeper)
	// migrations := map[uint64]func(sdk.Context) error{}
	// for ver, handler := range migrations {
	// 	if err := cfg.RegisterMigration(escrowv1alpha1.ModuleName, ver, handler); err != nil {
	// 		panic(fmt.Sprintf("failed to migrate x/%s from version %d to %d: %v", escrowv1alpha1.ModuleName, ver, ver+1, err))
	// 	}
	// }
}

// ____________________________________________________________________________

var _ module.HasGenesis = (*AppModule)(nil)

// InitGenesis performs genesis initialization for the module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var gs escrowv1alpha1.GenesisState
	cdc.MustUnmarshalJSON(data, &gs)

	if err := am.keeper.InitGenesis(sdk.UnwrapSDKContext(ctx), &gs); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the exported genesis state as raw bytes for the module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs, err := am.keeper.ExportGenesis(sdk.UnwrapSDKContext(ctx))
	if err != nil {
		panic(err)
	}

	return cdc.MustMarshalJSON(gs)
}

// ____________________________________________________________________________

var _ module.HasConsensusVersion = (*AppModule)(nil)

func (AppModule) ConsensusVersion() uint64 { return 1 }
