package internftv1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internft "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func TestMsgSend(t *testing.T) {
	addrs := createAddresses(2, "addr")
	classID := createIDs(1, "class")[0]
	nftID := createIDs(1, "nft")[0]

	testCases := map[string]struct {
		sender    sdk.AccAddress
		recipient sdk.AccAddress
		classID   string
		nftID     string
		err       error
	}{
		"valid msg": {
			sender:    addrs[0],
			recipient: addrs[1],
			classID:   classID,
			nftID:   nftID,
		},
		"invalid sender": {
			recipient: addrs[1],
			classID:   classID,
			nftID:   nftID,
			err:       sdkerrors.ErrInvalidAddress,
		},
		"invalid recipient": {
			sender:  addrs[0],
			classID: classID,
			nftID:   nftID,
			err:     sdkerrors.ErrInvalidAddress,
		},
		"invalid class id": {
			sender:    addrs[0],
			recipient: addrs[1],
			nftID:   nftID,
			err:       internft.ErrInvalidClassID,
		},
		"invalid nft id": {
			sender:    addrs[0],
			recipient: addrs[1],
			classID: classID,
			err:       internft.ErrInvalidNFTID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := internft.MsgSend{
				Sender:    tc.sender.String(),
				Recipient: tc.recipient.String(),
				Nft: internft.NFT{
					ClassId: tc.classID,
					Id:      tc.nftID,
				},
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func TestMsgNewClass(t *testing.T) {
	operator := createAddresses(1, "addr")[0]
	allowedClassID := operator.String()
	otherClassID := createIDs(1, "class")[0]
	const traitID = "uri"

	testCases := map[string]struct {
		operator   sdk.AccAddress
		classID string
		traitID string
		err     error
	}{
		"valid msg": {
			operator:   operator,
			classID: allowedClassID,
			traitID: traitID,
		},
		"invalid operator": {
			traitID: traitID,
			classID: allowedClassID,
			err:     sdkerrors.ErrInvalidAddress,
		},
		"invalid class id": {
			operator: operator,
			traitID: traitID,
			err:   internft.ErrInvalidClassID,
		},
		"invalid trait id": {
			operator: operator,
			classID: allowedClassID,
			err:   internft.ErrInvalidTraitID,
		},
		"unauthorized": {
			operator:   operator,
			classID: otherClassID,
			traitID: traitID,
			err: sdkerrors.ErrUnauthorized,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := internft.MsgNewClass{
				Operator: tc.operator.String(),
				Class: internft.Class{
					Id: tc.classID,
				},
				Traits: []internft.Trait{
					{
						Id: tc.traitID,
					},
				},
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func TestMsgUpdateClass(t *testing.T) {
	operator := createAddresses(1, "addr")[0]
	classID := operator.String()

	testCases := map[string]struct {
		operator sdk.AccAddress
		classID string
		err     error
	}{
		"valid msg": {
			operator: operator,
			classID: classID,
		},
		"invalid opreator": {
			classID: classID,
			err: sdkerrors.ErrInvalidAddress,
		},
		"invalid class id": {
			operator: operator,
			err: internft.ErrInvalidClassID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := internft.MsgUpdateClass{
				Operator: tc.operator.String(),
				Class: internft.Class{
					Id: tc.classID,
				},
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func TestMsgMintNFT(t *testing.T) {
	addrs := createAddresses(2, "addr")
	operator := addrs[0]
	recipient := addrs[1]
	classID := operator.String()
	nftID := createIDs(1, "nft")[0]
	const traitID = "uri"

	testCases := map[string]struct {
		operator sdk.AccAddress
		recipient sdk.AccAddress
		classID   string
		nftID string
		traitID   string
		err       error
	}{
		"valid msg": {
			operator: operator,
			recipient: recipient,
			classID:   classID,
			nftID: nftID,
			traitID:   traitID,
		},
		"invalid class id": {
			operator: operator,
			recipient: recipient,
			nftID: nftID,
			traitID:   traitID,
			err:       internft.ErrInvalidClassID,
		},
		"invalid trait id": {
			operator: operator,
			classID: classID,
			nftID: nftID,
			err:     internft.ErrInvalidTraitID,
		},
		"invalid recipient": {
			operator: operator,
			classID: classID,
			nftID: nftID,
			traitID: traitID,
			err:     sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := internft.MsgMintNFT{
				Operator: tc.operator.String(),
				Recipient: tc.recipient.String(),
				Nft: internft.NFT{
					ClassId: tc.classID,
					Id: tc.nftID,
				},
				Properties: []internft.Property{
					{
						Id: tc.traitID,
					},
				},
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func TestMsgBurnNFT(t *testing.T) {
	addr := createAddresses(1, "addr")[0]
	classID := createIDs(1, "class")[0]
	nftID := createIDs(1, "nft")[0]

	testCases := map[string]struct {
		owner   sdk.AccAddress
		classID string
		nftID string
		err     error
	}{
		"valid msg": {
			owner:   addr,
			classID: classID,
			nftID: nftID,
		},
		"invalid owner": {
			classID: classID,
			nftID: nftID,
			err:     sdkerrors.ErrInvalidAddress,
		},
		"invalid class id": {
			owner: addr,
			nftID: nftID,
			err:   internft.ErrInvalidClassID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := internft.MsgBurnNFT{
				Owner: tc.owner.String(),
				Nft: internft.NFT{
					ClassId: tc.classID,
					Id:      tc.nftID,
				},
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}

func TestMsgUpdateNFT(t *testing.T) {
	classID := createIDs(1, "class")[0]
	nftID := createIDs(1, "nft")[0]
	traitID := "uri"

	testCases := map[string]struct {
		classID  string
		nftID string
		traitIDs []string
		err      error
	}{
		"valid msg": {
			classID: classID,
			nftID: nftID,
			traitIDs: []string{
				traitID,
			},
		},
		"invalid class id": {
			nftID: nftID,
			traitIDs: []string{
				traitID,
			},
			err: internft.ErrInvalidClassID,
		},
		"empty properties": {
			classID: classID,
			nftID: nftID,
			err:     sdkerrors.ErrInvalidRequest,
		},
		"invalid trait id": {
			classID: classID,
			nftID: nftID,
			traitIDs: []string{
				"",
			},
			err: internft.ErrInvalidTraitID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			properties := make([]internft.Property, len(tc.traitIDs))
			for i, id := range tc.traitIDs {
				properties[i] = internft.Property{
					Id: id,
				}
			}

			msg := internft.MsgUpdateNFT{
				Nft: internft.NFT{
					ClassId: tc.classID,
					Id:      tc.nftID,
				},
				Properties: properties,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}
		})
	}
}
