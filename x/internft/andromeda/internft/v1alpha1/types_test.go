package internftv1alpha1_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internft "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

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

func TestClass(t *testing.T) {
	id := createIDs(1, "class")[0]

	testCases := map[string]struct {
		id  string
		err error
	}{
		"valid class": {
			id: id,
		},
		"invalid id": {
			err: internft.ErrInvalidClassID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			class := internft.Class{
				Id: tc.id,
			}

			err := class.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestTraits(t *testing.T) {
	testCases := map[string]struct {
		ids []string
		err error
	}{
		"valid traits": {},
		"invalid id": {
			ids: []string{
				"",
			},
			err: internft.ErrInvalidTraitID,
		},
		"duplicate id": {
			ids: []string{
				"uri",
				"uri",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			traits := make([]internft.Trait, len(tc.ids))
			for i, id := range tc.ids {
				traits[i] = internft.Trait{
					Id: id,
				}
			}

			err := internft.Traits(traits).ValidateBasic()
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestNFT(t *testing.T) {
	classIDs := createIDs(1, "class")
	nftIDs := createIDs(1, "nft")

	// ValidateBasic()
	testCases := map[string]struct {
		classID string
		id      string
		err     error
	}{
		"valid nft": {
			classID: classIDs[0],
			id:      nftIDs[0],
		},
		"invalid class id": {
			id:  nftIDs[0],
			err: internft.ErrInvalidClassID,
		},
		"invalid id": {
			classID: classIDs[0],
			err:     internft.ErrInvalidNFTID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			nft := internft.NFT{
				ClassId: tc.classID,
				Id:      tc.id,
			}

			err := nft.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestProperties(t *testing.T) {
	testCases := map[string]struct {
		ids []string
		err error
	}{
		"valid properties": {},
		"invalid id": {
			ids: []string{
				"",
			},
			err: internft.ErrInvalidTraitID,
		},
		"duplicate id": {
			ids: []string{
				"uri",
				"uri",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			properties := make([]internft.Property, len(tc.ids))
			for i, id := range tc.ids {
				properties[i] = internft.Property{
					Id: id,
				}
			}

			err := internft.Properties(properties).ValidateBasic()
			require.ErrorIs(t, err, tc.err)
		})
	}
}
