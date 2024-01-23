package internal

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

type Keeper struct {
	cdc codec.Codec

	schema collections.Schema

	assets collections.Map[collections.Pair[sdk.AccAddress, string], testv1alpha1.Asset]
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	prefix []byte,
) (*Keeper, error) {
	sb := collections.NewSchemaBuilder(storeService)

	if len(prefix) == 0 {
		prefix = assetsKeyPrefix
	}

	k := &Keeper{
		cdc: cdc,
		assets: collections.NewMap(sb, prefix, "assets",
			collections.PairKeyCodec(sdk.AccAddressKey, collections.StringKey),
			codec.CollValue[testv1alpha1.Asset](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		return nil, err
	}
	k.schema = schema

	return k, nil
}

func (k Keeper) addressBytesToString(addr []byte) (string, error) {
	addressCodec := k.cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addrStr, err := addressCodec.BytesToString(addr)
	if err != nil {
		return "", testv1alpha1.ErrInvalidAddress.Wrap(err.Error())
	}

	return addrStr, nil
}

func (k Keeper) addressStringToBytes(addr string) (sdk.AccAddress, error) {
	addressCodec := k.cdc.InterfaceRegistry().SigningContext().AddressCodec()
	addrBytes, err := addressCodec.StringToBytes(addr)
	if err != nil {
		return nil, testv1alpha1.ErrInvalidAddress.Wrap(err.Error())
	}

	return addrBytes, nil
}
