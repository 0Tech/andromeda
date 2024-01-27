package escrowv1alpha1

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
	errorCodeInvalidMessage
)

var (
	ErrUnimplemented  = errors.RegisterWithGRPCCode(errorCodespace, errorCodeUnimplemented, codes.Unimplemented, "unimplemented request")
	ErrInvalidAddress = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvalidAddress, codes.InvalidArgument, "invalid address")
	ErrInvalidMessage = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvalidMessage, codes.InvalidArgument, "invalid message")
)

// stateful errors
const (
	_ uint32 = iota + 1<<(32-1) - 1

	errorCodeInvariantBroken
	errorCodeAgentNotFound
	errorCodeProposalNotFound
	errorCodePermissionDenied
	errorCodeLargeMetadata
)

var (
	ErrInvariantBroken  = errors.RegisterWithGRPCCode(errorCodespace, errorCodeInvariantBroken, codes.Internal, "invariant broken")
	ErrAgentNotFound    = errors.RegisterWithGRPCCode(errorCodespace, errorCodeAgentNotFound, codes.NotFound, "agent not found")
	ErrProposalNotFound = errors.RegisterWithGRPCCode(errorCodespace, errorCodeProposalNotFound, codes.NotFound, "proposal not found")
	ErrPermissionDenied = errors.RegisterWithGRPCCode(errorCodespace, errorCodePermissionDenied, codes.PermissionDenied, "permission denied")
	ErrLargeMetadata    = errors.RegisterWithGRPCCode(errorCodespace, errorCodeLargeMetadata, codes.ResourceExhausted, "large metadata")
)
