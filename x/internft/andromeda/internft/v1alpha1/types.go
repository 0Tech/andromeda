package internftv1alpha1

import (
	"encoding/hex"
	"regexp"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Address sdk.AccAddress

func (a *Address) Parse(bech32 string) error {
	addr, err := sdk.AccAddressFromBech32(bech32)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrap(bech32)
	}
	*a = Address(addr)

	return nil
}

// TODO(@0Tech): implement caip
func (a Address) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(a)
}

type Reference string

var (
	referenceExpr = "[-.%[:alnum:]]{1,128}"
	referenceRegexp = regexp.MustCompilePOSIX("^" + referenceExpr + "$")
)

func (r *Reference) Parse(caip19 string) error {
	if !referenceRegexp.MatchString(caip19) {
		return errorsmod.Wrap(ErrInvalidID.Wrapf("reference must in form of %s", referenceExpr), caip19)
	}
	*r = Reference(caip19)

	return nil
}

type TokenID string

var (
	tokenIDExpr = "[-.%[:alnum:]]{1,78}"
	tokenIDRegexp = regexp.MustCompilePOSIX("^" + tokenIDExpr + "$")
)

func (t *TokenID) Parse(caip19 string) error {
	if !tokenIDRegexp.MatchString(caip19) {
		return errorsmod.Wrap(ErrInvalidID.Wrapf("token id must in form of %s", tokenIDExpr), caip19)
	}
	*t = TokenID(caip19)

	return nil
}

func (p Params) ValidateCompatibility() error {
	return nil
}

type ParamsInternal struct {
}

func (pi *ParamsInternal) Parse(p Params) error {
	if err := p.ValidateCompatibility(); err != nil {
		return err
	}

	return nil
}

func (c Class) ValidateCompatibility() error {
	if c.Id == "" {
		return ErrUnimplemented.Wrap("nil id")
	}

	return nil
}

type ClassInternal struct {
	ID Reference
}

func (ci *ClassInternal) Parse(c Class) error {
	if err := c.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ci.ID.Parse(c.Id); err != nil {
		return err
	}

	return nil
}

func (t Trait) ValidateCompatibility() error {
	if t.Id == "" {
		return ErrUnimplemented.Wrap("nil id")
	}

	if t.Mutability == Trait_MUTABILITY_UNSPECIFIED {
		return ErrUnimplemented.Wrap("nil mutability")
	}

	return nil
}

type TraitInternal struct {
	ID Reference
	Mutable bool
}

func (ti *TraitInternal) Parse(t Trait) error {
	if err := t.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ti.ID.Parse(t.Id); err != nil {
		return err
	}

	ti.Mutable = (t.Mutability == Trait_MUTABILITY_MUTABLE)

	return nil
}

func (t Token) ValidateCompatibility() error {
	if t.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	if t.Id == "" {
		return ErrUnimplemented.Wrap("nil id")
	}

	return nil
}

type TokenInternal struct {
	ClassID Reference
	ID TokenID
}

func (ti *TokenInternal) Parse(t Token) error {
	if err := t.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ti.ClassID.Parse(t.ClassId); err != nil {
		return err
	}

	if err := ti.ID.Parse(t.Id); err != nil {
		return err
	}

	return nil
}

func (p Property) ValidateCompatibility() error {
	if p.TraitId == "" {
		return ErrUnimplemented.Wrap("nil trait id")
	}

	if p.Fact == "" {
		return ErrUnimplemented.Wrap("nil fact")
	}

	return nil
}

type PropertyInternal struct {
	TraitID Reference
	Fact string
}

func (pi *PropertyInternal) Parse(p Property) error {
	if err := p.ValidateCompatibility(); err != nil {
		return err
	}

	if err := pi.TraitID.Parse(p.TraitId); err != nil {
		return err
	}

	pi.Fact = p.Fact

	return nil
}

// TODO(@0Tech): move to msg server
func ValidateOperator(operator, classID string) error {
	var addr Address
	if err := addr.Parse(operator); err != nil {
		return err
	}

	if classID != GetClassID(addr) {
		return errorsmod.Wrap(ErrPermissionDenied.Wrapf("not operator of class %s", classID), operator)
	}

	return nil
}

func GetClassID(operator Address) string {
	return hex.EncodeToString(operator)
}
