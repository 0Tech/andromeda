package module

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
	modulev1alpha1 "github.com/0tech/andromeda/x/test/api/andromeda/test/module/v1alpha1"
	"github.com/0tech/andromeda/x/test/keeper"
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

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// testv1alpha1.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	testv1alpha1.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := testv1alpha1.RegisterQueryHandlerClient(context.Background(), mux, testv1alpha1.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
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

var _ module.HasServices = (*AppModule)(nil)

func (am AppModule) RegisterServices(cfg module.Configurator) {
	testv1alpha1.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.keeper))
	testv1alpha1.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))

	// m := keeper.NewMigrator(am.keeper)
	// migrations := map[uint64]func(sdk.Context) error{}
	// for ver, handler := range migrations {
	// 	if err := cfg.RegisterMigration(testv1alpha1.ModuleName, ver, handler); err != nil {
	// 		panic(fmt.Sprintf("failed to migrate x/%s from version %d to %d: %v", testv1alpha1.ModuleName, ver, ver+1, err))
	// 	}
	// }
}

// ____________________________________________________________________________

var _ module.HasConsensusVersion = (*AppModule)(nil)

func (AppModule) ConsensusVersion() uint64 { return 1 }

// ____________________________________________________________________________

var _ appmodule.AppModule = (*AppModule)(nil)

func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule()        {}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

func init() {
	appmodule.Register(&modulev1alpha1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type TestInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *modulev1alpha1.Module
}

type TestOutputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in TestInputs) TestOutputs {
	k, err := keeper.NewKeeper(in.Cdc, in.StoreService, nil)
	if err != nil {
		panic(err)
	}

	m := NewAppModule(in.Cdc, *k)

	return TestOutputs{Keeper: *k, Module: m}
}
