package internal

import (
	"cosmossdk.io/collections"
	collcodec "cosmossdk.io/collections/codec"
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeService store.KVStoreService

	authority string

	schema collections.Schema
	params collections.Item[internftv1alpha1.Params]
	classes collections.Map[string, internftv1alpha1.Class]
	traits collections.Map[collections.Pair[string, string], internftv1alpha1.Trait]
	tokens collections.Map[collections.Pair[string, string], internftv1alpha1.Token]
	properties collections.Map[collections.Triple[string, string, string], internftv1alpha1.Property]
	owners collections.Map[collections.Pair[string, string], sdk.AccAddress]
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
) Keeper {
	// TODO(@0Tech): add authority check

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:       cdc,
		storeService: storeService,
		authority: authority,
		params: collections.NewItem(sb, paramsKey, "params",
			codec.CollValue[internftv1alpha1.Params](cdc)),
		classes: collections.NewMap(sb, classKeyPrefix, "classes",
			collections.StringKey,
			codec.CollValue[internftv1alpha1.Class](cdc)),
		traits: collections.NewMap(sb, traitKeyPrefix, "traits",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			codec.CollValue[internftv1alpha1.Trait](cdc)),
		tokens: collections.NewMap(sb, tokenKeyPrefix, "tokens",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			codec.CollValue[internftv1alpha1.Token](cdc)),
		properties: collections.NewMap(sb, propertyKeyPrefix, "properties",
			collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.StringKey),
			codec.CollValue[internftv1alpha1.Property](cdc)),
		owners: collections.NewMap(sb, ownerKeyPrefix, "owners",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			collcodec.KeyToValueCodec(sdk.AccAddressKey)),
	}

	if schema, err := sb.Build(); err != nil {
		panic(err)
	} else {
		k.schema = schema
	}

	return k
}
