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
	errorCodeInvalidID
)

var (
	ErrUnimplemented = errors.RegisterWithGRPCCode(errorCodespace, errorCodeUnimplemented, codes.Unimplemented, "unimplemented request")
	ErrInvalidID     = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvalidID, codes.InvalidArgument, "invalid id")
)

// stateful errors
const (
	_ uint32 = iota + 1 << (32 - 1) - 1
	errorCodeClassNotFound
	errorCodeClassAlreadyExists
	errorCodeTraitNotFound
	errorCodeTraitImmutable
	errorCodeTokenNotFound
	errorCodeTokenAlreadyExists
	errorCodePermissionDenied
)

var (
	ErrClassNotFound      = errors.RegisterWithGRPCCode(errorCodespace, errorCodeClassNotFound, codes.NotFound, "class not found")
	ErrClassAlreadyExists = errors.RegisterWithGRPCCode(errorCodespace, errorCodeClassAlreadyExists, codes.AlreadyExists, "class already exists")
	ErrTraitNotFound      = errors.RegisterWithGRPCCode(errorCodespace, errorCodeTraitNotFound, codes.NotFound, "trait not found")
	ErrTraitImmutable     = errors.RegisterWithGRPCCode(errorCodespace, errorCodeTraitImmutable, codes.FailedPrecondition, "trait immutable")
	ErrTokenNotFound        = errors.RegisterWithGRPCCode(errorCodespace, errorCodeTokenNotFound, codes.NotFound, "token not found")
	ErrTokenAlreadyExists = errors.RegisterWithGRPCCode(errorCodespace, errorCodeTokenAlreadyExists, codes.AlreadyExists, "token already exists")
	ErrPermissionDenied = errors.RegisterWithGRPCCode(errorCodespace, errorCodePermissionDenied, codes.PermissionDenied, "permission denied")
)
