package internftv1alpha1_test

import (
	"fmt"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func TestMsgSend(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgSend) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgSend]{
		{
			"valid sender": {
				malleate: func(subject *internftv1alpha1.MsgSend) {
					subject.Sender = createAddresses(2, "addr")[0].String()
				},
			},
			"empty sender": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid recipient": {
				malleate: func(subject *internftv1alpha1.MsgSend) {
					subject.Recipient = createAddresses(2, "addr")[1].String()
				},
			},
			"empty recipient": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgSend) {
					subject.Nft.ClassId = createIDs(1, "class")[0]
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"valid nft id": {
				malleate: func(subject *internftv1alpha1.MsgSend) {
					subject.Nft.Id = createIDs(1, "nft")[0]
				},
			},
			"empty nft id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidNFTID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestMsgNewClass(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgNewClass) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgNewClass]{
		{
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgNewClass) {
					subject.Operator = createAddresses(2, "addr")[0].String()
				},
			},
			"empty operator": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgNewClass) {
					subject.Class.Id = createAddresses(2, "addr")[0].String()
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
			"unauthorized class id": {
				malleate: func(subject *internftv1alpha1.MsgNewClass) {
					subject.Class.Id = createAddresses(2, "addr")[1].String()
				},
				err: func() error {
					return sdkerrors.ErrUnauthorized
				},
			},
		},
	}

	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)

		added := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgNewClass]{
			{
				"no trait": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						added = false
					},
				},
				"add trait": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						added = true
						subject.Traits = append(subject.Traits, internftv1alpha1.Trait{})
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						if added {
							subject.Traits[len(subject.Traits) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if added {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"immutable": {
				},
				"mutable": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						if added {
							subject.Traits[len(subject.Traits) - 1].Variable = true
						}
					},
				},
			},
		}...)

		addedDup := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgNewClass]{
			{
				"no duplicate trait": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						addedDup = false
					},
				},
				"add duplicate trait": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						addedDup = true
						subject.Traits = append(subject.Traits, internftv1alpha1.Trait{})
					},
					err: func() error {
						if added && addedDup {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						if addedDup {
							subject.Traits[len(subject.Traits) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if addedDup {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"immutable": {
				},
				"mutable": {
					malleate: func(subject *internftv1alpha1.MsgNewClass) {
						if addedDup {
							subject.Traits[len(subject.Traits) - 1].Variable = true
						}
					},
				},
			},
		}...)
	}

	doTest(t, tester, cases)
}

func TestMsgUpdateClass(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgUpdateClass) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgUpdateClass]{
		{
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgUpdateClass) {
					subject.Operator = createAddresses(2, "addr")[0].String()
				},
			},
			"empty operator": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateClass) {
					subject.Class.Id = createAddresses(2, "addr")[0].String()
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
			"unauthorized class id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateClass) {
					subject.Class.Id = createAddresses(2, "addr")[1].String()
				},
				err: func() error {
					return sdkerrors.ErrUnauthorized
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestMsgMintNFT(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgMintNFT) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgMintNFT]{
		{
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgMintNFT) {
					subject.Operator = createAddresses(2, "addr")[0].String()
				},
			},
			"empty operator": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid recipient": {
				malleate: func(subject *internftv1alpha1.MsgMintNFT)  {
					subject.Recipient = createAddresses(2, "addr")[1].String()
				},
			},
			"empty recipient": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgMintNFT) {
					subject.Nft.ClassId = createAddresses(2, "addr")[0].String()
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
			"unauthorized class id": {
				malleate: func(subject *internftv1alpha1.MsgMintNFT) {
					subject.Nft.ClassId = createAddresses(2, "addr")[1].String()
				},
				err: func() error {
					return sdkerrors.ErrUnauthorized
				},
			},
		},
		{
			"valid nft id": {
				malleate: func(subject *internftv1alpha1.MsgMintNFT) {
					subject.Nft.Id = createIDs(1, "nft")[0]
				},
			},
			"empty nft id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidNFTID
				},
			},
		},
	}

	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)
		fact := fmt.Sprintf("fact%02d", i)

		added := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgMintNFT]{
			{
				"no property": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						added = false
					},
				},
				"add property": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						added = true
						subject.Properties = append(subject.Properties, internftv1alpha1.Property{})
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						if added {
							subject.Properties[len(subject.Properties) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if added {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						if added {
							subject.Properties[len(subject.Properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)

		addedDup := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgMintNFT]{
			{
				"no duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						addedDup = false
					},
				},
				"add duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						addedDup = true
						subject.Properties = append(subject.Properties, internftv1alpha1.Property{})
					},
					err: func() error {
						if added && addedDup {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						if addedDup {
							subject.Properties[len(subject.Properties) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if addedDup {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *internftv1alpha1.MsgMintNFT) {
						if addedDup {
							subject.Properties[len(subject.Properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)
	}

	doTest(t, tester, cases)
}

func TestMsgBurnNFT(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgBurnNFT) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgBurnNFT]{
		{
			"valid owner": {
				malleate: func(subject *internftv1alpha1.MsgBurnNFT) {
					subject.Owner = createAddresses(1, "addr")[0].String()
				},
			},
			"empty owner": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgBurnNFT) {
					subject.Nft.ClassId = createIDs(1, "class")[0]
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"valid nft id": {
				malleate: func(subject *internftv1alpha1.MsgBurnNFT) {
					subject.Nft.Id = createIDs(1, "nft")[0]
				},
			},
			"empty nft id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidNFTID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestMsgUpdateNFT(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgUpdateNFT) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgUpdateNFT]{
		{
			"valid owner": {
				malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
					subject.Owner = createAddresses(1, "addr")[0].String()
				},
			},
			"empty owner": {
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
					subject.Nft.ClassId = createIDs(1, "class")[0]
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"valid nft id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
					subject.Nft.Id = createIDs(1, "nft")[0]
				},
			},
			"empty nft id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidNFTID
				},
			},
		},
	}

	addedEver := false
	cases = append(cases, map[string]Case[internftv1alpha1.MsgUpdateNFT]{
		"": {
			malleate: func(_ *internftv1alpha1.MsgUpdateNFT) {
				addedEver = false
			},
		},
	})
	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)
		fact := fmt.Sprintf("fact%02d", i)

		added := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateNFT]{
			{
				"no property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						added = false
					},
				},
				"add property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						addedEver = true
						added = true
						subject.Properties = append(subject.Properties, internftv1alpha1.Property{})
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						if added {
							subject.Properties[len(subject.Properties) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if added {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						if added {
							subject.Properties[len(subject.Properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)

		addedDup := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateNFT]{
			{
				"no duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						addedDup = false
					},
				},
				"add duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						addedEver = true
						addedDup = true
						subject.Properties = append(subject.Properties, internftv1alpha1.Property{})
					},
					err: func() error {
						if added && addedDup {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						if addedDup {
							subject.Properties[len(subject.Properties) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if addedDup {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *internftv1alpha1.MsgUpdateNFT) {
						if addedDup {
							subject.Properties[len(subject.Properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)
	}
	cases = append(cases, map[string]Case[internftv1alpha1.MsgUpdateNFT]{
		"": {
			err: func() error {
				if !addedEver {
					return sdkerrors.ErrInvalidRequest
				}
				return nil
			},
		},
	})

	doTest(t, tester, cases)
}
