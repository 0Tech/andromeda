package internal_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/keeper/expected"
	keeper "github.com/0tech/andromeda/x/escrow/keeper/internal"
	"github.com/0tech/andromeda/x/escrow/module"
	escrowtestutil "github.com/0tech/andromeda/x/escrow/testutil"
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

	keeper keeper.Keeper

	queryServer escrowv1alpha1.QueryServer
	msgServer   escrowv1alpha1.MsgServer

	seller   sdk.AccAddress
	buyer    sdk.AccAddress
	stranger sdk.AccAddress

	agentDedicated sdk.AccAddress
	agentAny       sdk.AccAddress
	agentIdle      sdk.AccAddress
	agentLast      sdk.AccAddress

	proposalDedicated uint64
	proposalAny       uint64
	proposalLast      uint64
}

func (s *KeeperTestSuite) addressBytesToString(address sdk.AccAddress) string {
	addressStr, err := s.addressCodec.BytesToString(address)
	s.NoError(err)

	return addressStr
}

func (s *KeeperTestSuite) SetupTest() {
	var authKeeper expected.AuthKeeper
	s.keeper, authKeeper, s.addressCodec, s.ctx = setupEscrowKeeper(s.T())

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	err := s.keeper.InitGenesis(s.ctx, s.keeper.DefaultGenesis())
	s.NoError(err)

	// create accounts
	addresses := []*sdk.AccAddress{
		&s.seller,
		&s.buyer,
		&s.stranger,
	}
	for i, address := range simtestutil.CreateRandomAccounts(len(addresses)) {
		*addresses[i] = address

		account := &authtypes.BaseAccount{
			Address: s.addressBytesToString(address),
		}
		authKeeper.SetAccount(s.ctx, authKeeper.NewAccount(s.ctx, account))
	}

	// seller creates agents
	for _, agent := range []*sdk.AccAddress{
		&s.agentDedicated,
		&s.agentAny,
		&s.agentIdle,
	} {
		var err error
		*agent, err = s.keeper.CreateAgent(s.ctx, s.seller)
		s.NoError(err)
	}
	s.agentLast = s.agentIdle

	// submit proposal for the buyer
	s.proposalDedicated, err = s.keeper.SubmitProposal(s.ctx, s.seller, s.agentDedicated, nil, nil)
	s.NoError(err)

	// submit proposal for anyone
	s.proposalAny, err = s.keeper.SubmitProposal(s.ctx, s.seller, s.agentAny, nil, nil)
	s.NoError(err)

	// store the last proposal id
	s.proposalLast = s.proposalAny
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func setupEscrowKeeper(t *testing.T) (
	keeper.Keeper,
	expected.AuthKeeper,
	address.Codec,
	context.Context,
) {
	key := storetypes.NewKVStoreKey(escrowv1alpha1.ModuleName)
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

	escrowv1alpha1.RegisterInterfaces(ir)

	bapp := baseapp.NewBaseApp(
		"escrow",
		log.NewNopLogger(),
		testCtx.DB,
		encCfg.TxConfig.TxDecoder(),
	)
	bapp.SetInterfaceRegistry(ir)

	ctrl := gomock.NewController(t)

	// mock auth keeper
	authKeeper := escrowtestutil.NewMockAuthKeeper(ctrl)
	authPrefix := []byte{0xff, 0x00}

	accountNumberPrefix := append(append([]byte{}, authPrefix...), 0x00)
	getAccountNumber := func(ctx context.Context) uint64 {
		store := runtime.NewKVStoreService(key).OpenKVStore(ctx)

		bz, err := store.Get(accountNumberPrefix)
		assert.NoError(t, err)

		prev := math.ZeroUint()
		if bz != nil {
			err = prev.Unmarshal(bz)
			assert.NoError(t, err)
		}

		curr := prev.Incr()

		bz, err = curr.Marshal()
		assert.NoError(t, err)

		err = store.Set(accountNumberPrefix, bz)
		assert.NoError(t, err)

		return curr.Uint64()
	}

	accountPrefix := append(append([]byte{}, authPrefix...), 0x01)
	hasAccount := func(ctx context.Context, address sdk.AccAddress) bool {
		store := runtime.NewKVStoreService(key).OpenKVStore(ctx)

		key := append(append([]byte{}, accountPrefix...), address...)
		has, err := store.Has(key)
		assert.NoError(t, err)

		return has
	}
	setAccount := func(ctx context.Context, account sdk.AccountI) {
		store := runtime.NewKVStoreService(key).OpenKVStore(ctx)

		bz, err := encCfg.Codec.Marshal(account)
		assert.NoError(t, err)

		key := append(append([]byte{}, accountPrefix...), account.GetAddress()...)
		err = store.Set(key, bz)
		assert.NoError(t, err)
	}

	authKeeper.EXPECT().HasAccount(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, address sdk.AccAddress) bool {
		return hasAccount(ctx, address)
	}).AnyTimes()
	authKeeper.EXPECT().NewAccount(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, account sdk.AccountI) sdk.AccountI {
		accNum := getAccountNumber(ctx)
		err := account.SetAccountNumber(accNum)
		assert.NoError(t, err)
		return account
	}).AnyTimes()
	authKeeper.EXPECT().SetAccount(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, account sdk.AccountI) {
		setAccount(ctx, account)
	}).AnyTimes()

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	k, err := keeper.NewKeeper(encCfg.Codec, runtime.NewKVStoreService(key), authority, bapp.MsgServiceRouter(), authKeeper)
	assert.NoError(t, err)

	msgServer := keeper.NewMsgServer(*k)
	queryServer := keeper.NewQueryServer(*k)

	escrowv1alpha1.RegisterMsgServer(bapp.MsgServiceRouter(), msgServer)
	escrowv1alpha1.RegisterQueryServer(bapp.GRPCQueryRouter(), queryServer)

	return *k, authKeeper, encCfg.InterfaceRegistry.SigningContext().AddressCodec(), testCtx.Ctx
}
