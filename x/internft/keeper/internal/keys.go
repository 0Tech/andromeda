package internal

import (
	"cosmossdk.io/collections"
)

var (
	paramsKey = collections.NewPrefix(0x00)

	classKeyPrefix      = collections.NewPrefix(0x10)
	traitKeyPrefix      = collections.NewPrefix(0x11)

	nftKeyPrefix      = collections.NewPrefix(0x20)
	propertyKeyPrefix = collections.NewPrefix(0x21)
	ownerKeyPrefix = collections.NewPrefix(0x22)
)
