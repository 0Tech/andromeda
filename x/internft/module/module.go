package module

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
	modulev1alpha1 "github.com/0tech/andromeda/x/internft/api/andromeda/internft/module/v1alpha1"
	"github.com/0tech/andromeda/x/internft/keeper"
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
	return internftv1alpha1.ModuleName
}

// RegisterLegacyAminoCodec registers the types on the LegacyAmino codec
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	internftv1alpha1.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	internftv1alpha1.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := internftv1alpha1.RegisterQueryHandlerClient(context.Background(), mux, internftv1alpha1.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// ____________________________________________________________________________

var _ module.HasGenesisBasics = (*AppModuleBasic)(nil)

// DefaultGenesis returns default genesis state as raw bytes for the module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(internftv1alpha1.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var gs internftv1alpha1.GenesisState
	if err := cdc.UnmarshalJSON(bz, &gs); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", internftv1alpha1.ModuleName, err)
	}

	return gs.ValidateBasic()
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

var _ module.AppModule = (*AppModule)(nil)

// AppModule implements an application module for the module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

// ____________________________________________________________________________

var _ module.HasInvariants = (*AppModule)(nil)

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO(@0Tech): add invariants
}

// ____________________________________________________________________________

var _ module.HasServices = (*AppModule)(nil)

func (am AppModule) RegisterServices(cfg module.Configurator) {
	internftv1alpha1.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
	internftv1alpha1.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))

	// m := keeper.NewMigrator(am.keeper)
	// migrations := map[uint64]func(sdk.Context) error{}
	// for ver, handler := range migrations {
	// 	if err := cfg.RegisterMigration(internftv1alpha1.ModuleName, ver, handler); err != nil {
	// 		panic(fmt.Sprintf("failed to migrate x/%s from version %d to %d: %v", internftv1alpha1.ModuleName, ver, ver+1, err))
	// 	}
	// }
}

// ____________________________________________________________________________

var _ module.HasGenesis = (*AppModule)(nil)

// InitGenesis performs genesis initialization for the module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var gs internftv1alpha1.GenesisState
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

// ____________________________________________________________________________

var _ appmodule.AppModule = (*AppModule)(nil)

func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule() {}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

func init() {
	appmodule.Register(&modulev1alpha1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type InterNFTInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc    codec.Codec
	Config *modulev1alpha1.Module
}

type InterNFTOutputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in InterNFTInputs) InterNFTOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k, err := keeper.NewKeeper(in.Cdc, in.StoreService, authority.String())
	if err != nil {
		panic(err)
	}

	m := NewAppModule(in.Cdc, *k)

	return InterNFTOutputs{Keeper: *k, Module: m}
}
