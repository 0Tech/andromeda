package internal_test

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
	keeper "github.com/0tech/andromeda/x/internft/keeper/internal"
	"github.com/0tech/andromeda/x/internft/module"
	internfttestutil "github.com/0tech/andromeda/x/internft/testutil"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	keeper keeper.Keeper

	queryServer internftv1alpha1.QueryServer
	msgServer   internftv1alpha1.MsgServer

	vendor   sdk.AccAddress
	customer sdk.AccAddress
	stranger sdk.AccAddress

	mutableTraitID   string
	immutableTraitID string

	nftIDs map[string]string
}

func createAddresses(size int, prefix string) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, size)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(fmt.Sprintf("%s%d", prefix, i))
	}

	return addrs
}

func createIDs(size int, prefix string) []string {
	addrs := createAddresses(size, prefix)
	ids := make([]string, len(addrs))
	for i, addr := range addrs {
		ids[i] = addr.String()
	}

	return ids
}

//nolint:gosec
func randomString(size int) string {
	res := make([]rune, size)

	letters := []rune("0123456789abcdef")
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}

	return string(res)
}

type Case[T any] struct {
	malleate func(*T)
	err func() error
}

func doTest[T any](
	s *KeeperTestSuite,
	tester func(T) error,
	cases []map[string]Case[T]) {
	for iter := internfttestutil.NewCaseIterator(cases); iter.Valid(); iter.Next() {
		names := iter.Key()

		var subject T
		var errs []error
		for i, name := range names {
			c := cases[i][name]

			if malleate := c.malleate; malleate != nil {
				malleate(&subject)
			}
			if errGen := c.err; errGen != nil {
				if err := errGen(); err != nil {
					errs = append(errs, err)
				}
			}
		}

		testName := func(names []string) string {
			display := make([]string, 0, len(names))
			for _, name := range names {
				if len(name) != 0 {
					display = append(display, name)
				}
			}
			return strings.Join(display, ",")
		}
		s.Run(testName(names), func() {
			err := tester(subject)
			if len(errs) != 0 {
				s.Error(err)

				for _, candidate := range errs {
					if errors.Is(err, candidate) {
						return
					}
				}
				s.FailNow("unexpected error", err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(internftv1alpha1.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	s.keeper = keeper.NewKeeper(encCfg.Codec, storeService, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	// create accounts
	addresses := []*sdk.AccAddress{
		&s.vendor,
		&s.customer,
		&s.stranger,
	}
	for i, address := range createAddresses(len(addresses), "addr") {
		*addresses[i] = address
	}

	s.keeper.SetParams(s.ctx, internftv1alpha1.Params{})

	// vendor creates a class
	class := internftv1alpha1.Class{
		Id: s.vendor.String(),
	}
	err := class.ValidateBasic()
	s.Assert().NoError(err)

	s.mutableTraitID = "level"
	s.immutableTraitID = "color"

	traits := []internftv1alpha1.Trait{
		{
			Id:      s.mutableTraitID,
			Variable: true,
		},
		{
			Id: s.immutableTraitID,
		},
	}

	err = s.keeper.NewClass(s.ctx, class, traits)
	s.Assert().NoError(err)

	// vendor mints nfts to all accounts
	nftIDs := createIDs(len(addresses), "nft")
	s.nftIDs = map[string]string{}
	for i := range addresses {
		recipient := *addresses[i]
		nftID := nftIDs[i]
		s.nftIDs[recipient.String()] = nftID

		// except for stranger
		if recipient.Equals(s.stranger) {
			continue
		}

		nft := internftv1alpha1.NFT{
			ClassId: class.Id,
			Id: nftID,
		}
		properties := []internftv1alpha1.Property{
			{
				Id: s.mutableTraitID,
				Fact: "42",
			},
			{
				Id: s.immutableTraitID,
				Fact: "black",
			},
		}

		err := s.keeper.MintNFT(s.ctx, recipient, nft, properties)
		s.Assert().NoError(err)
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
