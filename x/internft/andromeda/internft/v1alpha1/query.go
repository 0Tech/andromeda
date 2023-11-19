package internftv1alpha1

func (q QueryParamsRequest) ValidateCompatibility() error {
	return nil
}

func (q QueryParamsRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	return nil
}

func (q QueryClassRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	return nil
}

func (q QueryClassRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	return nil
}

func (q QueryClassesRequest) ValidateCompatibility() error {
	return nil
}

func (q QueryClassesRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	return nil
}

func (q QueryTraitRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	if q.TraitId == "" {
		return ErrUnimplemented.Wrap("nil trait id")
	}

	return nil
}

func (q QueryTraitRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	return nil
}

func (q QueryTraitsRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	return nil
}

func (q QueryTraitsRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	return nil
}

func (q QueryTokenRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	if q.TokenId == "" {
		return ErrUnimplemented.Wrap("nil token id")
	}

	return nil
}

func (q QueryTokenRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	if err := ValidateTokenID(q.TokenId); err != nil {
		return err
	}

	return nil
}

func (q QueryTokensRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	return nil
}

func (q QueryTokensRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	return nil
}

func (q QueryPropertyRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	if q.TokenId == "" {
		return ErrUnimplemented.Wrap("nil token id")
	}

	if q.TraitId == "" {
		return ErrUnimplemented.Wrap("nil trait id")
	}

	return nil
}

func (q QueryPropertyRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	if err := ValidateTokenID(q.TokenId); err != nil {
		return err
	}

	return nil
}

func (q QueryPropertiesRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	if q.TokenId == "" {
		return ErrUnimplemented.Wrap("nil token id")
	}

	return nil
}

func (q QueryPropertiesRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	if err := ValidateTokenID(q.TokenId); err != nil {
		return err
	}

	return nil
}

func (q QueryOwnerRequest) ValidateCompatibility() error {
	if q.ClassId == "" {
		return ErrUnimplemented.Wrap("nil class id")
	}

	if q.TokenId == "" {
		return ErrUnimplemented.Wrap("nil token id")
	}

	return nil
}

func (q QueryOwnerRequest) ValidateBasic() error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := ValidateClassID(q.ClassId); err != nil {
		return err
	}

	if err := ValidateTokenID(q.TokenId); err != nil {
		return err
	}

	return nil
}
