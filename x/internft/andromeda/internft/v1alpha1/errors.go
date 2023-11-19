package internftv1alpha1

import (
	"google.golang.org/grpc/codes"

	"cosmossdk.io/errors"
)

const errorCodespace = ModuleName

// stateless errors
const (
	_ uint32 = iota + 1 << 1 - 1
	errorCodeUnimplemented
	errorCodeInvalidClassID
	errorCodeInvalidTokenID
)

var (
	ErrUnimplemented = errors.RegisterWithGRPCCode(errorCodespace, errorCodeUnimplemented, codes.Unimplemented, "unimplemented request")
	ErrInvalidClassID     = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvalidClassID, codes.InvalidArgument, "invalid class id")
	ErrInvalidTokenID       = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvalidTokenID, codes.InvalidArgument, "invalid token id")
)

// stateful errors
const (
	_ uint32 = iota + 1 << (32 - 1) - 1
	errorCodeClassNotFound
	errorCodeClassAlreadyExists
	errorCodeTraitNotFound
	errorCodeTraitImmutable
	errorCodeTokenNotFound
	errorCodePermissionDenied
)

var (
	ErrClassNotFound      = errors.RegisterWithGRPCCode(errorCodespace, errorCodeClassNotFound, codes.NotFound, "class not found")
	ErrClassAlreadyExists = errors.RegisterWithGRPCCode(errorCodespace, errorCodeClassAlreadyExists, codes.AlreadyExists, "class already exists")
	ErrTraitNotFound      = errors.RegisterWithGRPCCode(errorCodespace, errorCodeTraitNotFound, codes.NotFound, "trait not found")
	ErrTraitImmutable     = errors.RegisterWithGRPCCode(errorCodespace, errorCodeTraitImmutable, codes.FailedPrecondition, "trait immutable")
	ErrTokenNotFound        = errors.RegisterWithGRPCCode(errorCodespace, errorCodeTokenNotFound, codes.NotFound, "token not found")
	ErrPermissionDenied = errors.RegisterWithGRPCCode(errorCodespace, errorCodePermissionDenied, codes.PermissionDenied, "permission denied")
)
