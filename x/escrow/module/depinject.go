package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	modulev1alpha1 "github.com/0tech/andromeda/x/escrow/api/andromeda/escrow/module/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/keeper"
)

var _ appmodule.AppModule = (*AppModule)(nil)

func (AppModule) IsOnePerModuleType() {}
func (AppModule) IsAppModule()        {}

func init() {
	appmodule.Register(&modulev1alpha1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type Inputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *modulev1alpha1.Module

	Router     keeper.MessageRouter
	AuthKeeper keeper.AuthKeeper
}

type Outputs struct {
	depinject.Out

	Keeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in Inputs) Outputs {
	addressCodec := in.Cdc.InterfaceRegistry().SigningContext().AddressCodec()

	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		bz, err := addressCodec.StringToBytes(in.Config.Authority)
		if err != nil {
			authority = authtypes.NewModuleAddress(in.Config.Authority)
		} else {
			authority = bz
		}
	}

	k, err := keeper.NewKeeper(in.Cdc, in.StoreService, authority, in.Router, in.AuthKeeper)
	if err != nil {
		panic(err)
	}

	m := NewAppModule(*k)

	return Outputs{Keeper: *k, Module: m}
}
