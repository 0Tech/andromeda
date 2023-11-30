package internftv1alpha1

func (q QueryParamsRequest) ValidateCompatibility() error {
	return nil
}

type QueryParamsInternal struct {
}

func (qi *QueryParamsInternal) Parse(q QueryParamsRequest) error {
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

type QueryClassInternal struct {
	ClassID Reference
}

func (qi *QueryClassInternal) Parse(q QueryClassRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
		return err
	}

	return nil
}

func (q QueryClassesRequest) ValidateCompatibility() error {
	return nil
}

type QueryClassesInternal struct {
}

func (qi *QueryClassesInternal) Parse(q QueryClassesRequest) error {
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

type QueryTraitInternal struct {
	ClassID Reference
	TraitID Reference
}

func (qi *QueryTraitInternal) Parse(q QueryTraitRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
		return err
	}

	if err := qi.TraitID.Parse(q.TraitId); err != nil {
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

type QueryTraitsInternal struct {
	ClassID Reference
}

func (qi *QueryTraitsInternal) Parse(q QueryTraitsRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
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

type QueryTokenInternal struct {
	ClassID Reference
	TokenID TokenID
}

func (qi *QueryTokenInternal) Parse(q QueryTokenRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
		return err
	}

	if err := qi.TokenID.Parse(q.TokenId); err != nil {
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

type QueryTokensInternal struct {
	ClassID Reference
}

func (qi *QueryTokensInternal) Parse(q QueryTokensRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
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

type QueryPropertyInternal struct {
	ClassID Reference
	TokenID TokenID
	TraitID Reference
}

func (qi *QueryPropertyInternal) Parse(q QueryPropertyRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
		return err
	}

	if err := qi.TokenID.Parse(q.TokenId); err != nil {
		return err
	}

	if err := qi.TraitID.Parse(q.TraitId); err != nil {
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

type QueryPropertiesInternal struct {
	ClassID Reference
	TokenID TokenID
}

func (qi *QueryPropertiesInternal) Parse(q QueryPropertiesRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
		return err
	}

	if err := qi.TokenID.Parse(q.TokenId); err != nil {
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

type QueryOwnerInternal struct {
	ClassID Reference
	TokenID TokenID
}

func (qi *QueryOwnerInternal) Parse(q QueryOwnerRequest) error {
	if err := q.ValidateCompatibility(); err != nil {
		return err
	}

	if err := qi.ClassID.Parse(q.ClassId); err != nil {
		return err
	}

	if err := qi.TokenID.Parse(q.TokenId); err != nil {
		return err
	}

	return nil
}
