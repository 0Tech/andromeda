package internal

import (
	"cosmossdk.io/collections"
)

var assetsKeyPrefix = collections.NewPrefix(append([]byte{0xff}, []byte("andromeda/test/asset")...))
