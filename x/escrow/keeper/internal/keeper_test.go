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
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/gogoproto/proto"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/keeper/expected"
	keeper "github.com/0tech/andromeda/x/escrow/keeper/internal"
	"github.com/0tech/andromeda/x/escrow/module"
	escrowtestutil "github.com/0tech/andromeda/x/escrow/testutil"
	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
	testkeeper "github.com/0tech/andromeda/x/test/keeper"
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

func createRandomAccounts(size int) []sdk.AccAddress {
	return simtestutil.CreateRandomAccounts(size)
}

type KeeperTestSuite struct {
	suite.Suite

	ctx context.Context

	addressCodec address.Codec

	keeper *keeper.Keeper

	queryServer escrowv1alpha1.QueryServer
	msgServer   escrowv1alpha1.MsgServer

	seller   sdk.AccAddress
	buyer    sdk.AccAddress
	stranger sdk.AccAddress

	agentDedicated sdk.AccAddress
	agentAny       sdk.AccAddress
	agentIdle      sdk.AccAddress
}

func (s *KeeperTestSuite) addressBytesToString(address sdk.AccAddress) string {
	addressStr, err := s.addressCodec.BytesToString(address)
	s.NoError(err)

	return addressStr
}

func (s *KeeperTestSuite) encodeMsgs(msgs []sdk.Msg) []*codectypes.Any {
	encoded := make([]*codectypes.Any, len(msgs))
	for i, msg := range msgs {
		s.NotNil(msg)

		bz, err := proto.Marshal(msg)
		s.NoError(err)

		encoded[i] = &codectypes.Any{
			TypeUrl: sdk.MsgTypeURL(msg),
			Value:   bz,
		}
	}

	return encoded
}

func (s *KeeperTestSuite) SetupTest() {
	var authKeeper expected.AuthKeeper
	var testKeeper *testkeeper.Keeper
	var cdc codec.Codec
	cdc, s.ctx, s.keeper, authKeeper, testKeeper = setupKeepers(s.T())
	s.addressCodec = cdc.InterfaceRegistry().SigningContext().AddressCodec()

	s.queryServer = keeper.NewQueryServer(*s.keeper)
	s.msgServer = keeper.NewMsgServer(*s.keeper)

	err := s.keeper.InitGenesis(s.ctx, s.keeper.DefaultGenesis())
	s.NoError(err)

	// create accounts
	accounts := []*sdk.AccAddress{
		&s.seller,
		&s.buyer,
		&s.stranger,
	}
	for i, account := range createRandomAccounts(len(accounts)) {
		*accounts[i] = account

		account := &authtypes.BaseAccount{
			Address: s.addressBytesToString(account),
		}
		authKeeper.SetAccount(s.ctx, authKeeper.NewAccount(s.ctx, account))
	}

	// allocate assets
	testMsgServer := testkeeper.NewMsgServer(*testKeeper)
	for _, allocation := range []struct {
		account sdk.AccAddress
		assets  []string
	}{
		{
			account: s.seller,
			assets:  []string{"cat", "dog", "snake"},
		},
		{
			account: s.buyer,
			assets:  []string{"voucher"},
		},
		{
			account: s.stranger,
			assets:  []string{"voucher"},
		},
	} {
		accountStr := s.addressBytesToString(allocation.account)
		for _, asset := range allocation.assets {
			_, err := testMsgServer.Create(s.ctx, &testv1alpha1.MsgCreate{
				Creator: accountStr,
				Asset:   asset,
			})
			s.NoError(err)
		}
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

	// submit proposal for the buyer
	err = s.keeper.SubmitProposal(s.ctx, s.seller, s.agentDedicated,
		s.encodeMsgs([]sdk.Msg{
			&testv1alpha1.MsgSend{
				Sender:    s.addressBytesToString(s.seller),
				Recipient: s.addressBytesToString(s.agentDedicated),
				Asset:     "cat",
			},
		}),
		s.encodeMsgs([]sdk.Msg{
			&testv1alpha1.MsgSend{
				Sender:    s.addressBytesToString(s.agentDedicated),
				Recipient: s.addressBytesToString(s.seller),
				Asset:     "voucher",
			},
			&testv1alpha1.MsgSend{
				Sender:    s.addressBytesToString(s.agentDedicated),
				Recipient: s.addressBytesToString(s.buyer),
				Asset:     "cat",
			},
		}),
		"sell a cat to buyer for a voucher",
	)
	s.NoError(err)

	// submit proposal for anyone
	err = s.keeper.SubmitProposal(s.ctx, s.seller, s.agentAny,
		s.encodeMsgs([]sdk.Msg{
			&testv1alpha1.MsgSend{
				Sender:    s.addressBytesToString(s.seller),
				Recipient: s.addressBytesToString(s.agentAny),
				Asset:     "dog",
			},
		}),
		s.encodeMsgs([]sdk.Msg{
			&testv1alpha1.MsgSend{
				Sender:    s.addressBytesToString(s.agentAny),
				Recipient: s.addressBytesToString(s.seller),
				Asset:     "voucher",
			},
		}),
		"sell a dog for a voucher",
	)
	s.NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func setupEscrowKeeper(t *testing.T) (
	codec.Codec,
	context.Context,
	*keeper.Keeper,
) {
	cdc, ctx, keeper, _, _ := setupKeepers(t)
	return cdc, ctx, keeper
}

func setupKeepers(t *testing.T) (
	codec.Codec,
	context.Context,
	*keeper.Keeper,
	expected.AuthKeeper,
	*testkeeper.Keeper,
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

	bapp := baseapp.NewBaseApp(
		"escrow",
		log.NewNopLogger(),
		testCtx.DB,
		encCfg.TxConfig.TxDecoder(),
	)
	bapp.SetInterfaceRegistry(ir)

	ctrl := gomock.NewController(t)
	authKeeper := newAuthKeeper(t, ctrl, encCfg.Codec, key)

	escrowKeeper := newEscrowKeeper(t, bapp, encCfg.Codec, key, authKeeper)
	testKeeper := newTestKeeper(t, bapp, encCfg.Codec, key) // register test keeper

	return encCfg.Codec, testCtx.Ctx, escrowKeeper, authKeeper, testKeeper
}

// mock auth keeper
func newAuthKeeper(t *testing.T, ctrl *gomock.Controller, cdc codec.Codec, key *storetypes.KVStoreKey) expected.AuthKeeper {
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

		bz, err := cdc.Marshal(account)
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

	return authKeeper
}

func newEscrowKeeper(t *testing.T, bapp *baseapp.BaseApp, cdc codec.Codec, key *storetypes.KVStoreKey, authKeeper expected.AuthKeeper) *keeper.Keeper {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	escrowKeeper, err := keeper.NewKeeper(cdc, runtime.NewKVStoreService(key), authority, bapp.MsgServiceRouter(), authKeeper)
	assert.NoError(t, err)

	msgServer := keeper.NewMsgServer(*escrowKeeper)
	queryServer := keeper.NewQueryServer(*escrowKeeper)

	escrowv1alpha1.RegisterInterfaces(cdc.InterfaceRegistry())
	escrowv1alpha1.RegisterMsgServer(bapp.MsgServiceRouter(), msgServer)
	escrowv1alpha1.RegisterQueryServer(bapp.GRPCQueryRouter(), queryServer)

	return escrowKeeper
}

func newTestKeeper(t *testing.T, bapp *baseapp.BaseApp, cdc codec.Codec, key *storetypes.KVStoreKey) *testkeeper.Keeper {
	testKeeper, err := testkeeper.NewKeeper(cdc, runtime.NewKVStoreService(key), nil)
	assert.NoError(t, err)

	msgServer := testkeeper.NewMsgServer(*testKeeper)
	queryServer := testkeeper.NewQueryServer(*testKeeper)

	testv1alpha1.RegisterInterfaces(cdc.InterfaceRegistry())
	testv1alpha1.RegisterMsgServer(bapp.MsgServiceRouter(), msgServer)
	testv1alpha1.RegisterQueryServer(bapp.GRPCQueryRouter(), queryServer)

	return testKeeper
}
