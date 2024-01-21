package testv1alpha1

import (
	"google.golang.org/grpc/codes"

	"cosmossdk.io/errors"
)

const errorCodespace = ModuleName

// stateless errors
const (
	_ uint32 = iota + 1<<1 - 1

	errorCodeUnimplemented
	errorCodeInvalidAddress
)

var (
	ErrUnimplemented  = errors.RegisterWithGRPCCode(errorCodespace, errorCodeUnimplemented, codes.Unimplemented, "unimplemented request")
	ErrInvalidAddress = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvalidAddress, codes.InvalidArgument, "invalid address")
)

// stateful errors
const (
	_ uint32 = iota + 1<<(32-1) - 1

	errorCodeInvariantBroken
	errorCodeAssetNotFound
	errorCodeAssetAlreadyExists
)

var (
	ErrInvariantBroken    = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvariantBroken, codes.Internal, "invariant broken")
	ErrAssetNotFound      = errors.RegisterWithGRPCCode(errorCodespace, errorCodeAssetNotFound, codes.NotFound, "asset not found")
	ErrAssetAlreadyExists = errors.RegisterWithGRPCCode(errorCodespace, errorCodeAssetAlreadyExists, codes.AlreadyExists, "asset already exists")
)
