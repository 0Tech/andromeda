package internftv1alpha1

import (
	"google.golang.org/grpc/codes"

	"cosmossdk.io/errors"
)

const errorCodespace = ModuleName

var (
	ErrInvalidClassID     = errors.RegisterWithGRPCCode(errorCodespace, 2, codes.InvalidArgument, "invalid class id")
	ErrInvalidTraitID     = errors.RegisterWithGRPCCode(errorCodespace, 3, codes.InvalidArgument, "invalid trait id")
	ErrInvalidTokenID       = errors.RegisterWithGRPCCode(errorCodespace, 4, codes.InvalidArgument, "invalid token id")
	ErrClassNotFound      = errors.RegisterWithGRPCCode(errorCodespace, 5, codes.NotFound, "class not found")
	ErrClassAlreadyExists = errors.RegisterWithGRPCCode(errorCodespace, 6, codes.AlreadyExists, "class already exists")
	ErrTraitNotFound      = errors.RegisterWithGRPCCode(errorCodespace, 7, codes.NotFound, "trait not found")
	ErrTraitImmutable     = errors.RegisterWithGRPCCode(errorCodespace, 8, codes.FailedPrecondition, "trait immutable")
	ErrTokenNotFound        = errors.RegisterWithGRPCCode(errorCodespace, 9, codes.NotFound, "token not found")
	ErrInsufficientToken    = errors.RegisterWithGRPCCode(errorCodespace, 10, codes.FailedPrecondition, "insufficient token")
)
