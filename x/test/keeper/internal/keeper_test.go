package internal_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
	keeper "github.com/0tech/andromeda/x/test/keeper/internal"
	"github.com/0tech/andromeda/x/test/module"
)

const notInBech32 = "addresslverygoodtestingaddress" // does not include the separator

//nolint:gosec
func randomString(size int) string {
	res := make([]rune, size)

	letters := []rune("0123456789abcdef")
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}

	return string(res)
}

type KeeperTestSuite struct {
	suite.Suite

	ctx context.Context

	addressCodec address.Codec

	keeper *keeper.Keeper

	queryServer testv1alpha1.QueryServer
	msgServer   testv1alpha1.MsgServer

	notPetPerson sdk.AccAddress
	catPerson    sdk.AccAddress
	dogPerson    sdk.AccAddress
	petPerson    sdk.AccAddress

	cat string
	dog string
}

func (s *KeeperTestSuite) addressBytesToString(address sdk.AccAddress) string {
	addressStr, err := s.addressCodec.BytesToString(address)
	s.NoError(err)

	return addressStr
}

func (s *KeeperTestSuite) SetupTest() {
	s.addressCodec, s.ctx, s.keeper = setupTestKeeper(s.T())

	s.queryServer = keeper.NewQueryServer(*s.keeper)
	s.msgServer = keeper.NewMsgServer(*s.keeper)

	s.cat = "cat"
	s.dog = "dog"

	accounts := []*sdk.AccAddress{
		&s.notPetPerson,
		&s.catPerson,
		&s.dogPerson,
		&s.petPerson,
	}
	for i, account := range simtestutil.CreateRandomAccounts(len(accounts)) {
		*accounts[i] = account
	}

	for _, allocation := range []struct {
		account sdk.AccAddress
		assets  []string
	}{
		{
			account: s.notPetPerson,
		},
		{
			account: s.catPerson,
			assets:  []string{s.cat},
		},
		{
			account: s.dogPerson,
			assets:  []string{s.dog},
		},
		{
			account: s.petPerson,
			assets:  []string{s.cat, s.dog},
		},
	} {
		for _, asset := range allocation.assets {
			err := s.keeper.Create(s.ctx, allocation.account, asset)
			s.NoError(err)
		}
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func setupTestKeeper(t *testing.T) (
	address.Codec,
	context.Context,
	*keeper.Keeper,
) {
	key := storetypes.NewKVStoreKey(testv1alpha1.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContextWithDB(t, key, tkey)
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})

	mainPrefix := randomString(8)
	ir := codectestutil.CodecOptions{
		AccAddressPrefix: mainPrefix,
		ValAddressPrefix: mainPrefix + "valoper",
	}.NewInterfaceRegistry()
	encCfg.InterfaceRegistry = ir
	encCfg.Codec = codec.NewProtoCodec(ir)

	bapp := baseapp.NewBaseApp(
		"test",
		log.NewNopLogger(),
		testCtx.DB,
		encCfg.TxConfig.TxDecoder(),
	)
	bapp.SetInterfaceRegistry(ir)

	testKeeper, err := keeper.NewKeeper(encCfg.Codec, runtime.NewKVStoreService(key), nil)
	assert.NoError(t, err)

	msgServer := keeper.NewMsgServer(*testKeeper)
	queryServer := keeper.NewQueryServer(*testKeeper)

	testv1alpha1.RegisterInterfaces(ir)
	testv1alpha1.RegisterMsgServer(bapp.MsgServiceRouter(), msgServer)
	testv1alpha1.RegisterQueryServer(bapp.GRPCQueryRouter(), queryServer)

	return encCfg.InterfaceRegistry.SigningContext().AddressCodec(), testCtx.Ctx, testKeeper
}
