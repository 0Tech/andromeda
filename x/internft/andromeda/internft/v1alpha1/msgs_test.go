package internftv1alpha1_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func TestMsgSend(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.MsgSend) error{
		{
			"valid_sender": func(subject *internftv1alpha1.MsgSend) error {
				subject.Sender = createAddresses(2, "addr")[0].String()
				return nil
			},
			"empty_sender": func(subject *internftv1alpha1.MsgSend) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_recipient": func(subject *internftv1alpha1.MsgSend) error {
				subject.Recipient = createAddresses(2, "addr")[1].String()
				return nil
			},
			"empty_recipient": func(subject *internftv1alpha1.MsgSend) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_class_id": func(subject *internftv1alpha1.MsgSend) error {
				subject.Nft.ClassId = createIDs(1, "class")[0]
				return nil
			},
			"empty_class_id": func(subject *internftv1alpha1.MsgSend) error {
				return internftv1alpha1.ErrInvalidClassID
			},
		},
		{
			"valid_nft_id": func(subject *internftv1alpha1.MsgSend) error {
				subject.Nft.Id = createIDs(1, "nft")[0]
				return nil
			},
			"empty_nft_id": func(subject *internftv1alpha1.MsgSend) error {
				return internftv1alpha1.ErrInvalidNFTID
			},
		},
	}
	tester := func(subject internftv1alpha1.MsgSend) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestMsgNewClass(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.MsgNewClass) error{
		{
			"valid_operator": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Operator = createAddresses(2, "addr")[0].String()
				return nil
			},
			"empty_operator": func(subject *internftv1alpha1.MsgNewClass) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_class_id": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Class.Id = createAddresses(2, "addr")[0].String()
				return nil
			},
			"empty_class_id": func(subject *internftv1alpha1.MsgNewClass) error {
				return internftv1alpha1.ErrInvalidClassID
			},
			"unauthorized_class_id": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Class.Id = createAddresses(2, "addr")[1].String()
				return sdkerrors.ErrUnauthorized
			},
		},
		{
			"no_trait": func(subject *internftv1alpha1.MsgNewClass) error {
				return nil
			},
			"duplicate_trait": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Traits = []internftv1alpha1.Trait{
					{Id: "color"},
					{Id: "color", Variable: true},
				}
				return sdkerrors.ErrInvalidRequest
			},
			"immutable_trait": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Traits = []internftv1alpha1.Trait{
					{Id: "color"},
				}
				return nil
			},
			"mutable_trait": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Traits = []internftv1alpha1.Trait{
					{Id: "color", Variable: true},
				}
				return nil
			},
			"empty_trait_id_immutable_trait": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Traits = []internftv1alpha1.Trait{{}}
				return internftv1alpha1.ErrInvalidTraitID
			},
			"empty_trait_id_mutable_trait": func(subject *internftv1alpha1.MsgNewClass) error {
				subject.Traits = []internftv1alpha1.Trait{
					{Variable: true},
				}
				return internftv1alpha1.ErrInvalidTraitID
			},
		},
	}
	tester := func(subject internftv1alpha1.MsgNewClass) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestMsgUpdateClass(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.MsgUpdateClass) error{
		{
			"valid_operator": func(subject *internftv1alpha1.MsgUpdateClass) error {
				subject.Operator = createAddresses(2, "addr")[0].String()
				return nil
			},
			"empty_operator": func(subject *internftv1alpha1.MsgUpdateClass) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_class_id": func(subject *internftv1alpha1.MsgUpdateClass) error {
				subject.Class.Id = createAddresses(2, "addr")[0].String()
				return nil
			},
			"empty_class_id": func(subject *internftv1alpha1.MsgUpdateClass) error {
				return internftv1alpha1.ErrInvalidClassID
			},
			"unauthorized_class_id": func(subject *internftv1alpha1.MsgUpdateClass) error {
				subject.Class.Id = createAddresses(2, "addr")[1].String()
				return sdkerrors.ErrUnauthorized
			},
		},
	}
	tester := func(subject internftv1alpha1.MsgUpdateClass) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestMsgMintNFT(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.MsgMintNFT) error{
		{
			"valid_operator": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Operator = createAddresses(2, "addr")[0].String()
				return nil
			},
			"empty_operator": func(subject *internftv1alpha1.MsgMintNFT) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_recipient": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Recipient = createAddresses(2, "addr")[1].String()
				return nil
			},
			"empty_recipient": func(subject *internftv1alpha1.MsgMintNFT) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_class_id": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Nft.ClassId = createAddresses(2, "addr")[0].String()
				return nil
			},
			"empty_class_id": func(subject *internftv1alpha1.MsgMintNFT) error {
				return internftv1alpha1.ErrInvalidClassID
			},
			"unauthorized_class_id": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Nft.ClassId = createAddresses(2, "addr")[1].String()
				return sdkerrors.ErrUnauthorized
			},
		},
		{
			"valid_nft_id": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Nft.Id = createIDs(1, "nft")[0]
				return nil
			},
			"empty_nft_id": func(subject *internftv1alpha1.MsgMintNFT) error {
				return internftv1alpha1.ErrInvalidNFTID
			},
		},
		{
			"no_property": func(subject *internftv1alpha1.MsgMintNFT) error {
				return nil
			},
			"duplicate_property": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Properties = []internftv1alpha1.Property{
					{Id: "color"},
					{Id: "color", Fact: "black"},
				}
				return sdkerrors.ErrInvalidRequest
			},
			"empty_trait_id_empty_fact": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Properties = []internftv1alpha1.Property{{}}
				return internftv1alpha1.ErrInvalidTraitID
			},
			"empty_trait_id_nonempty_fact": func(subject *internftv1alpha1.MsgMintNFT) error {
				subject.Properties = []internftv1alpha1.Property{
					{Fact: "black"},
				}
				return internftv1alpha1.ErrInvalidTraitID
			},
		},
	}
	tester := func(subject internftv1alpha1.MsgMintNFT) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestMsgBurnNFT(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.MsgBurnNFT) error{
		{
			"valid_owner": func(subject *internftv1alpha1.MsgBurnNFT) error {
				subject.Owner = createAddresses(1, "addr")[0].String()
				return nil
			},
			"empty_owner": func(subject *internftv1alpha1.MsgBurnNFT) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_class_id": func(subject *internftv1alpha1.MsgBurnNFT) error {
				subject.Nft.ClassId = createIDs(1, "class")[0]
				return nil
			},
			"empty_class_id": func(subject *internftv1alpha1.MsgBurnNFT) error {
				return internftv1alpha1.ErrInvalidClassID
			},
		},
		{
			"valid_nft_id": func(subject *internftv1alpha1.MsgBurnNFT) error {
				subject.Nft.Id = createIDs(1, "nft")[0]
				return nil
			},
			"empty_nft_id": func(subject *internftv1alpha1.MsgBurnNFT) error {
				return internftv1alpha1.ErrInvalidNFTID
			},
		},
	}
	tester := func(subject internftv1alpha1.MsgBurnNFT) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestMsgUpdateNFT(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.MsgUpdateNFT) error{
		{
			"valid_owner": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				subject.Owner = createAddresses(1, "addr")[0].String()
				return nil
			},
			"empty_owner": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				return sdkerrors.ErrInvalidAddress
			},
		},
		{
			"valid_class_id": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				subject.Nft.ClassId = createIDs(1, "class")[0]
				return nil
			},
			"empty_class_id": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				return internftv1alpha1.ErrInvalidClassID
			},
		},
		{
			"valid_nft_id": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				subject.Nft.Id = createIDs(1, "nft")[0]
				return nil
			},
			"empty_nft_id": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				return internftv1alpha1.ErrInvalidNFTID
			},
		},
		{
			"no_property": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				return sdkerrors.ErrInvalidRequest
			},
			"duplicate_property": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				subject.Properties = []internftv1alpha1.Property{
					{Id: "color"},
					{Id: "color", Fact: "black"},
				}
				return sdkerrors.ErrInvalidRequest
			},
			"empty_trait_id_empty_fact": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				subject.Properties = []internftv1alpha1.Property{{}}
				return internftv1alpha1.ErrInvalidTraitID
			},
			"empty_trait_id_nonempty_fact": func(subject *internftv1alpha1.MsgUpdateNFT) error {
				subject.Properties = []internftv1alpha1.Property{
					{Fact: "black"},
				}
				return internftv1alpha1.ErrInvalidTraitID
			},
		},
	}
	tester := func(subject internftv1alpha1.MsgUpdateNFT) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}
