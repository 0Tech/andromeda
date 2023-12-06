package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestMsgSendToken() {
	tester := func(subject internftv1alpha1.MsgSendToken) error {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.msgServer.SendToken(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&internftv1alpha1.EventSendToken{
			Sender: subject.Sender,
			Recipient: subject.Recipient,
			Token: subject.Token,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
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
					subject.Sender = s.customer.String()
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
					subject.Recipient = s.stranger.String()
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
					subject.Token.ClassId = s.classID
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
					subject.Token.Id = s.tokenIDs[s.customer.String()]
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

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestMsgCreateClass() {
	tester := func(subject internftv1alpha1.MsgCreateClass) error {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.msgServer.CreateClass(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&internftv1alpha1.EventCreateClass{
			Operator: subject.Operator,
			Class: subject.Class,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
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
					subject.Operator = s.customer.String()
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
					subject.Class.Id = internftv1alpha1.GetClassID(internftv1alpha1.Address(s.customer))
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

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestMsgUpdateTrait() {
	tester := func(subject internftv1alpha1.MsgUpdateTrait) error {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.msgServer.UpdateTrait(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&internftv1alpha1.EventUpdateTrait{
			Operator: subject.Operator,
			Class: subject.Class,
			Trait: subject.Trait,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
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
					subject.Operator = s.vendor.String()
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
					subject.Class.Id = s.classID
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
					subject.Trait.Id = s.mutableTraitID
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

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestMsgMintToken() {
	tester := func(subject internftv1alpha1.MsgMintToken) error {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.msgServer.MintToken(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&internftv1alpha1.EventMintToken{
			Operator: subject.Operator,
			Token: subject.Token,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
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
					subject.Operator = s.vendor.String()
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
					subject.Token.ClassId = s.classID
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
					subject.Token.Id = createIDs(3, "token")[2]
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

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestMsgBurnToken() {
	tester := func(subject internftv1alpha1.MsgBurnToken) error {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.msgServer.BurnToken(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&internftv1alpha1.EventBurnToken{
			Owner: subject.Owner,
			Token: subject.Token,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
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
					subject.Owner = s.customer.String()
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
					subject.Token.ClassId = s.classID
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
					subject.Token.Id = s.tokenIDs[s.customer.String()]
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

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestMsgUpdateProperty() {
	tester := func(subject internftv1alpha1.MsgUpdateProperty) error {
		ctx, _ := s.ctx.CacheContext()
		res, err := s.msgServer.UpdateProperty(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&internftv1alpha1.EventUpdateProperty{
			Operator: subject.Operator,
			Token: subject.Token,
			Property: subject.Property,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
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
					subject.Operator = s.vendor.String()
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
					subject.Token.ClassId = s.classID
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
					subject.Token.Id = s.tokenIDs[s.customer.String()]
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
					subject.Property.TraitId = s.mutableTraitID
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

	doTest(s, tester, cases)
}
