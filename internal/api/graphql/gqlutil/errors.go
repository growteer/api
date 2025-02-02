package gqlutil

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/growteer/api/internal/app/apperrors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ErrType = string

const (
	errTypeBadInput        ErrType = "BAD_INPUT"
	errTypeInternal        ErrType = "INTERNAL_SERVER_ERROR"
	errTypeNotFound        ErrType = "NOT_FOUND"
	errTypeUnauthenticated ErrType = "UNAUTHENTICATED"
	errTypeUnauthorized    ErrType = "UNAUTHORIZED"
)

func Recover(ctx context.Context, err interface{}) error {
	asErr, ok := err.(error)
	if ok {
		return apperrors.Internal{
			Code:    apperrors.ErrCodeInternalError,
			Message: "panic while executing graphql request",
			Wrapped: asErr,
		}
	}

	return apperrors.Internal{
		Code:    apperrors.ErrCodeInternalError,
		Message: fmt.Sprintf("panic while executing graphql request: %v", err),
	}
}

func PresentError(ctx context.Context, err error) *gqlerror.Error {
	slog.Error(err.Error())
	graphql.AddError(ctx, toGQLError(ctx, err))

	return nil
}

func toGQLError(ctx context.Context, err error) *gqlerror.Error {
	var badInput apperrors.BadInput
	if errors.As(err, &badInput) {
		return badInputToGQLError(ctx, &badInput)
	}

	var internal apperrors.Internal
	if errors.As(err, &internal) {
		return internalToGQLError(ctx, &internal)
	}

	var notFound apperrors.NotFound
	if errors.As(err, &notFound) {
		return notFoundToGQLError(ctx, &notFound)
	}

	var unauthenticated apperrors.Unauthenticated
	if errors.As(err, &unauthenticated) {
		return unauthenticatedToGQLError(ctx, &unauthenticated)
	}

	var unauthorized apperrors.Unauthorized
	if errors.As(err, &unauthorized) {
		return unauthorizedToGQLError(ctx, &unauthorized)
	}

	return &gqlerror.Error{
		Message: apperrors.ErrCodeInternalError,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeInternal,
		},
	}
}

func badInputToGQLError(ctx context.Context, err *apperrors.BadInput) *gqlerror.Error {
	return &gqlerror.Error{
		Message: err.Code,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeBadInput,
		},
	}
}

func internalToGQLError(ctx context.Context, err *apperrors.Internal) *gqlerror.Error {
	return &gqlerror.Error{
		Message: err.Code,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeInternal,
		},
	}
}

func notFoundToGQLError(ctx context.Context, err *apperrors.NotFound) *gqlerror.Error {
	return &gqlerror.Error{
		Message: err.Code,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeNotFound,
		},
	}
}

func unauthenticatedToGQLError(ctx context.Context, err *apperrors.Unauthenticated) *gqlerror.Error {
	return &gqlerror.Error{
		Message: err.Code,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeUnauthenticated,
		},
	}
}

func unauthorizedToGQLError(ctx context.Context, err *apperrors.Unauthorized) *gqlerror.Error {
	return &gqlerror.Error{
		Message: err.Code,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeUnauthorized,
		},
	}
}
