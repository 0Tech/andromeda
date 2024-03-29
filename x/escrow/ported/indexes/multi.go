package indexes

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
)

// Multi defines the most common index. It can be used to create a reference between
// a field of value and its primary key. Multiple primary keys can be mapped to the same
// reference key as the index does not enforce uniqueness constraints.
type Multi[ReferenceKey, PrimaryKey, Value any] struct {
	getRefKey func(pk PrimaryKey, value Value) (ReferenceKey, error)
	refKeys   collections.KeySet[collections.Pair[ReferenceKey, PrimaryKey]]
}

// NewMulti instantiates a new Multi instance given a schema,
// a Prefix, the humanized name for the index, the reference key key codec
// and the primary key key codec. The getRefKeyFunc is a function that
// given the primary key and value returns the referencing key.
func NewMulti[ReferenceKey, PrimaryKey, Value any](
	schema *collections.SchemaBuilder,
	prefix collections.Prefix,
	name string,
	refCodec codec.KeyCodec[ReferenceKey],
	pkCodec codec.KeyCodec[PrimaryKey],
	getRefKeyFunc func(pk PrimaryKey, value Value) (ReferenceKey, error),
) *Multi[ReferenceKey, PrimaryKey, Value] {
	return &Multi[ReferenceKey, PrimaryKey, Value]{
		getRefKey: getRefKeyFunc,
		refKeys:   collections.NewKeySet(schema, prefix, name, collections.PairKeyCodec(refCodec, pkCodec)),
	}
}

func (i *Multi[ReferenceKey, PrimaryKey, Value]) Reference(ctx context.Context, pk PrimaryKey, newValue Value, lazyOldValue func() (Value, error)) error {
	oldValue, err := lazyOldValue()
	switch {
	// if no error it means the value existed, and we need to remove the old indexes
	case err == nil:
		err = i.unreference(ctx, pk, oldValue)
		if err != nil {
			return err
		}
	// if error is ErrNotFound, it means that the object does not exist, so we're creating indexes for the first time.
	// we do nothing.
	case errors.Is(err, collections.ErrNotFound):
	// default case means that there was some other error
	default:
		return err
	}
	// create new indexes
	refKey, err := i.getRefKey(pk, newValue)
	if err != nil {
		return err
	}
	return i.refKeys.Set(ctx, collections.Join(refKey, pk))
}

func (i *Multi[ReferenceKey, PrimaryKey, Value]) Unreference(ctx context.Context, pk PrimaryKey, getValue func() (Value, error)) error {
	value, err := getValue()
	if err != nil {
		return err
	}
	return i.unreference(ctx, pk, value)
}

func (i *Multi[ReferenceKey, PrimaryKey, Value]) unreference(ctx context.Context, pk PrimaryKey, value Value) error {
	refKey, err := i.getRefKey(pk, value)
	if err != nil {
		return err
	}
	return i.refKeys.Remove(ctx, collections.Join(refKey, pk))
}

func (i *Multi[ReferenceKey, PrimaryKey, Value]) Iterate(ctx context.Context, ranger collections.Ranger[collections.Pair[ReferenceKey, PrimaryKey]]) (MultiIterator[ReferenceKey, PrimaryKey], error) {
	iter, err := i.refKeys.Iterate(ctx, ranger)
	return (MultiIterator[ReferenceKey, PrimaryKey])(iter), err
}

func (i *Multi[ReferenceKey, PrimaryKey, Value]) Walk(
	ctx context.Context,
	ranger collections.Ranger[collections.Pair[ReferenceKey, PrimaryKey]],
	walkFunc func(indexingKey ReferenceKey, indexedKey PrimaryKey) (stop bool, err error),
) error {
	return i.refKeys.Walk(ctx, ranger, func(key collections.Pair[ReferenceKey, PrimaryKey]) (bool, error) {
		return walkFunc(key.K1(), key.K2())
	})
}

func (i *Multi[ReferenceKey, PrimaryKey, Value]) IterateRaw(ctx context.Context, start, end []byte, order collections.Order) (_ collections.Iterator[collections.Pair[ReferenceKey, PrimaryKey], collections.NoValue], err error) {
	iter, err := i.refKeys.IterateRaw(ctx, start, end, order)
	if err != nil {
		return
	}
	return iter, nil
}

// MatchExact returns a MultiIterator containing all the primary keys referenced by the provided reference key.
func (i *Multi[ReferenceKey, PrimaryKey, Value]) MatchExact(ctx context.Context, refKey ReferenceKey) (MultiIterator[ReferenceKey, PrimaryKey], error) {
	return i.Iterate(ctx, collections.NewPrefixedPairRange[ReferenceKey, PrimaryKey](refKey))
}

func (i *Multi[K1, K2, Value]) KeyCodec() codec.KeyCodec[collections.Pair[K1, K2]] {
	return i.refKeys.KeyCodec()
}

// MultiIterator is just a KeySetIterator with key as Pair[ReferenceKey, PrimaryKey].
type MultiIterator[ReferenceKey, PrimaryKey any] collections.KeySetIterator[collections.Pair[ReferenceKey, PrimaryKey]]

// PrimaryKey returns the iterator's current primary key.
func (i MultiIterator[ReferenceKey, PrimaryKey]) PrimaryKey() (PrimaryKey, error) {
	fullKey, err := i.FullKey()
	return fullKey.K2(), err
}

// PrimaryKeys fully consumes the iterator and returns the list of primary keys.
func (i MultiIterator[ReferenceKey, PrimaryKey]) PrimaryKeys() ([]PrimaryKey, error) {
	fullKeys, err := i.FullKeys()
	if err != nil {
		return nil, err
	}
	pks := make([]PrimaryKey, len(fullKeys))
	for i, fullKey := range fullKeys {
		pks[i] = fullKey.K2()
	}
	return pks, nil
}

// FullKey returns the current full reference key as Pair[ReferenceKey, PrimaryKey].
func (i MultiIterator[ReferenceKey, PrimaryKey]) FullKey() (collections.Pair[ReferenceKey, PrimaryKey], error) {
	return (collections.KeySetIterator[collections.Pair[ReferenceKey, PrimaryKey]])(i).Key()
}

// FullKeys fully consumes the iterator and returns all the list of full reference keys.
func (i MultiIterator[ReferenceKey, PrimaryKey]) FullKeys() ([]collections.Pair[ReferenceKey, PrimaryKey], error) {
	return (collections.KeySetIterator[collections.Pair[ReferenceKey, PrimaryKey]])(i).Keys()
}

// Next advances the iterator.
func (i MultiIterator[ReferenceKey, PrimaryKey]) Next() {
	(collections.KeySetIterator[collections.Pair[ReferenceKey, PrimaryKey]])(i).Next()
}

// Valid asserts if the iterator is still valid or not.
func (i MultiIterator[ReferenceKey, PrimaryKey]) Valid() bool {
	return (collections.KeySetIterator[collections.Pair[ReferenceKey, PrimaryKey]])(i).Valid()
}

// Close closes the iterator.
func (i MultiIterator[ReferenceKey, PrimaryKey]) Close() error {
	return (collections.KeySetIterator[collections.Pair[ReferenceKey, PrimaryKey]])(i).Close()
}
