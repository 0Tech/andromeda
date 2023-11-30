package internftv1alpha1_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func TestMsgSendToken(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgSendToken) error {
		var parsed internftv1alpha1.MsgSendTokenInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.MsgSendToken]{
		{
			"nil sender": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid sender": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					subject.Sender = createAddresses(2, "addr")[0].String()
				},
			},
			"invalid sender": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					subject.Sender = "not-in-bech32"
				},
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
		{
			"nil recipient": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid recipient": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					subject.Recipient = createAddresses(2, "addr")[1].String()
				},
			},
			"invalid recipient": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					subject.Recipient = "not-in-bech32"
				},
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
	}

	addedToken := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgSendToken]{
		{
			"[nil token": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					addedToken = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil token": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					addedToken = true
					subject.Token = &internftv1alpha1.Token{}
				},
			},
		},
		{
			"nil class id": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = createIDs(1, "class")[0]
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id]": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id]": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					if !addedToken {
						return
					}
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"invalid token id]": {
				malleate: func(subject *internftv1alpha1.MsgSendToken) {
					if !addedToken {
						return
					}
					subject.Token.Id = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}...)

	doTest(t, tester, cases)
}

func TestMsgCreateClass(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgCreateClass) error {
		var parsed internftv1alpha1.MsgCreateClassInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.MsgCreateClass]{
		{
			"nil operator": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgCreateClass) {
					subject.Operator = createAddresses(1, "addr")[0].String()
				},
			},
			"invalid operator": {
				malleate: func(subject *internftv1alpha1.MsgCreateClass) {
					subject.Operator = "not-in-bech32"
				},
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
	}

	addedClass := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgCreateClass]{
		{
			"[nil class": {
				malleate: func(subject *internftv1alpha1.MsgCreateClass) {
					addedClass = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil class": {
				malleate: func(subject *internftv1alpha1.MsgCreateClass) {
					addedClass = true
					subject.Class = &internftv1alpha1.Class{}
				},
			},
		},
		{
			"nil class id]": {
				err: func() error {
					if !addedClass {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id]": {
				malleate: func(subject *internftv1alpha1.MsgCreateClass) {
					if !addedClass {
						return
					}
					subject.Class.Id = createIDs(1, "class")[0]
				},
			},
			"invalid class id]": {
				malleate: func(subject *internftv1alpha1.MsgCreateClass) {
					if !addedClass {
						return
					}
					subject.Class.Id = "not/in/caip19"
				},
				err: func() error {
					if !addedClass {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}...)

	doTest(t, tester, cases)
}

func TestMsgUpdateTrait(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgUpdateTrait) error {
		var parsed internftv1alpha1.MsgUpdateTraitInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.MsgUpdateTrait]{
		{
			"nil operator": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					subject.Operator = createAddresses(1, "addr")[0].String()
				},
			},
			"invalid operator": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					subject.Operator = "not-in-bech32"
				},
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
	}

	addedClass := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateTrait]{
		{
			"[nil class": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					addedClass = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil class": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					addedClass = true
					subject.Class = &internftv1alpha1.Class{}
				},
			},
		},
		{
			"nil class id]": {
				err: func() error {
					if !addedClass {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id]": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					if !addedClass {
						return
					}
					subject.Class.Id = createIDs(1, "class")[0]
				},
			},
			"invalid class id]": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					if !addedClass {
						return
					}
					subject.Class.Id = "not/in/caip19"
				},
				err: func() error {
					if !addedClass {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}...)

	addedTrait := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateTrait]{
		{
			"[nil trait": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					addedTrait = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil trait": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					addedTrait = true
					subject.Trait = &internftv1alpha1.Trait{}
				},
			},
		},
		{
			"nil trait id": {
				err: func() error {
					if !addedTrait {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid trait id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					if !addedTrait {
						return
					}
					subject.Trait.Id = createIDs(1, "trait")[0]
				},
			},
			"invalid trait id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					if !addedTrait {
						return
					}
					subject.Trait.Id = "not/in/caip19"
				},
				err: func() error {
					if !addedTrait {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil mutability]": {
				err: func() error {
					if !addedTrait {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"immutable]": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					if !addedTrait {
						return
					}
					subject.Trait.Mutability = internftv1alpha1.Trait_MUTABILITY_IMMUTABLE
				},
			},
			"mutable]": {
				malleate: func(subject *internftv1alpha1.MsgUpdateTrait) {
					if !addedTrait {
						return
					}
					subject.Trait.Mutability = internftv1alpha1.Trait_MUTABILITY_MUTABLE
				},
			},
		},
	}...)

	doTest(t, tester, cases)
}

func TestMsgMintToken(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgMintToken) error {
		var parsed internftv1alpha1.MsgMintTokenInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.MsgMintToken]{
		{
			"nil operator": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					subject.Operator = createAddresses(1, "addr")[0].String()
				},
			},
			"invalid operator": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					subject.Operator = "not-in-bech32"
				},
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
	}

	addedToken := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgMintToken]{
		{
			"[nil token": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					addedToken = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil token": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					addedToken = true
					subject.Token = &internftv1alpha1.Token{}
				},
			},
		},
		{
			"nil class id": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = createIDs(1, "class")[0]
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id]": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id]": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					if !addedToken {
						return
					}
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"invalid token id]": {
				malleate: func(subject *internftv1alpha1.MsgMintToken) {
					if !addedToken {
						return
					}
					subject.Token.Id = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}...)

	doTest(t, tester, cases)
}

func TestMsgBurnToken(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgBurnToken) error {
		var parsed internftv1alpha1.MsgBurnTokenInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.MsgBurnToken]{
		{
			"nil owner": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid owner": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					subject.Owner = createAddresses(1, "addr")[0].String()
				},
			},
			"invalid owner": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					subject.Owner = "not-in-bech32"
				},
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
	}

	addedToken := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgBurnToken]{
		{
			"[nil token": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					addedToken = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil token": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					addedToken = true
					subject.Token = &internftv1alpha1.Token{}
				},
			},
		},
		{
			"nil class id": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = createIDs(1, "class")[0]
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id]": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id]": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					if !addedToken {
						return
					}
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"invalid token id]": {
				malleate: func(subject *internftv1alpha1.MsgBurnToken) {
					if !addedToken {
						return
					}
					subject.Token.Id = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}...)

	doTest(t, tester, cases)
}

func TestMsgUpdateProperty(t *testing.T) {
	tester := func(subject internftv1alpha1.MsgUpdateProperty) error {
		var parsed internftv1alpha1.MsgUpdatePropertyInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.MsgUpdateProperty]{
		{
			"nil operator": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid operator": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					subject.Operator = createAddresses(1, "addr")[0].String()
				},
			},
			"invalid operator": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					subject.Operator = "not-in-bech32"
				},
				err: func() error {
					return sdkerrors.ErrInvalidAddress
				},
			},
		},
	}

	addedToken := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateProperty]{
		{
			"[nil token": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					addedToken = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil token": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					addedToken = true
					subject.Token = &internftv1alpha1.Token{}
				},
			},
		},
		{
			"nil class id": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = createIDs(1, "class")[0]
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					if !addedToken {
						return
					}
					subject.Token.ClassId = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id]": {
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id]": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					if !addedToken {
						return
					}
					subject.Token.Id = createIDs(1, "token")[0]
				},
			},
			"invalid token id]": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					if !addedToken {
						return
					}
					subject.Token.Id = "not/in/caip19"
				},
				err: func() error {
					if !addedToken {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}...)

	addedProperty := false
	cases = append(cases, []map[string]Case[internftv1alpha1.MsgUpdateProperty]{
		{
			"[nil property": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					addedProperty = false
				},
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"[non-nil property": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					addedProperty = true
					subject.Property = &internftv1alpha1.Property{}
				},
			},
		},
		{
			"nil trait id": {
				err: func() error {
					if !addedProperty {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid trait id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					if !addedProperty {
						return
					}
					subject.Property.TraitId = createIDs(1, "trait")[0]
				},
			},
			"invalid trait id": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					if !addedProperty {
						return
					}
					subject.Property.TraitId = "not/in/caip19"
				},
				err: func() error {
					if !addedProperty {
						return nil
					}
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil fact": {
				err: func() error {
					if !addedProperty {
						return nil
					}
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid fact": {
				malleate: func(subject *internftv1alpha1.MsgUpdateProperty) {
					if addedProperty {
						subject.Property.Fact = "valid fact"
					}
				},
			},
		},
	}...)

	doTest(t, tester, cases)
}
