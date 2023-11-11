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
					subject.Token.ClassId = createIDs(1, "class")[0]
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"valid token id": {
				malleate: func(subject *internftv1alpha1.MsgSend) {
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"empty token id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidTokenID
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

func TestMsgNewToken(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgNewToken) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgNewToken]{
		{
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgNewToken) {
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
				malleate: func(subject *internftv1alpha1.MsgNewToken)  {
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
				malleate: func(subject *internftv1alpha1.MsgNewToken) {
					subject.Token.ClassId = createAddresses(2, "addr")[0].String()
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
			"unauthorized class id": {
				malleate: func(subject *internftv1alpha1.MsgNewToken) {
					subject.Token.ClassId = createAddresses(2, "addr")[1].String()
				},
				err: func() error {
					return sdkerrors.ErrUnauthorized
				},
			},
		},
		{
			"valid token id": {
				malleate: func(subject *internftv1alpha1.MsgNewToken) {
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"empty token id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidTokenID
				},
			},
		},
	}

	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)
		fact := fmt.Sprintf("fact%02d", i)

		added := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgNewToken]{
			{
				"no property": {
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
						added = false
					},
				},
				"add property": {
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
						added = true
						subject.Properties = append(subject.Properties, internftv1alpha1.Property{})
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
						if added {
							subject.Properties[len(subject.Properties) - 1].TraitId = traitID
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
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
						if added {
							subject.Properties[len(subject.Properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)

		addedDup := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgNewToken]{
			{
				"no duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
						addedDup = false
					},
				},
				"add duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
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
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
						if addedDup {
							subject.Properties[len(subject.Properties) - 1].TraitId = traitID
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
					malleate: func(subject *internftv1alpha1.MsgNewToken) {
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

func TestMsgBurnToken(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgBurnToken) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgBurnToken]{
		{
			"valid owner": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
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
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					subject.Token.ClassId = createIDs(1, "class")[0]
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"valid token id": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"empty token id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidTokenID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestMsgUpdateToken(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgUpdateToken) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.MsgUpdateToken]{
		{
			"valid owner": {
				malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
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
				malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
					subject.Token.ClassId = createIDs(1, "class")[0]
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"valid token id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"empty token id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidTokenID
				},
			},
		},
	}

	addedEver := false
	cases = append(cases, map[string]Case[internftv1alpha1.MsgUpdateToken]{
		"": {
			malleate: func(_ *internftv1alpha1.MsgUpdateToken) {
				addedEver = false
			},
		},
	})
	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)
		fact := fmt.Sprintf("fact%02d", i)

		added := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateToken]{
			{
				"no property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
						added = false
					},
				},
				"add property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
						addedEver = true
						added = true
						subject.Properties = append(subject.Properties, internftv1alpha1.Property{})
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
						if added {
							subject.Properties[len(subject.Properties) - 1].TraitId = traitID
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
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
						if added {
							subject.Properties[len(subject.Properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)

		addedDup := false
		cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateToken]{
			{
				"no duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
						addedDup = false
					},
				},
				"add duplicate property": {
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
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
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
						if addedDup {
							subject.Properties[len(subject.Properties) - 1].TraitId = traitID
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
					malleate: func(subject *internftv1alpha1.MsgUpdateToken) {
						if addedDup {
							subject.Properties[len(subject.Properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)
	}
	cases = append(cases, map[string]Case[internftv1alpha1.MsgUpdateToken]{
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
