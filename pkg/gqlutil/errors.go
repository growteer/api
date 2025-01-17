package gqlutil

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ErrType = string
type ErrCode = string

const (
	errTypeBadRequest ErrType = "BAD_REQUEST"
	errTypeInternal ErrType = "INTERNAL_SERVER_ERROR"
)

func BadInputError(ctx context.Context, message string, code ErrCode, err error) error {
	gqlErr := &gqlerror.Error{
		Err: err,
		Message: message,
		Path: graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeBadRequest,
			"code": code,
		},
	}

	graphql.AddError(ctx, gqlErr)
	return gqlErr
}

func InternalError(ctx context.Context, message string, err error) error {
	gqlErr := &gqlerror.Error{
		Err: err,
		Message: message,
		Path: graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeInternal,
			"code": ErrCodeInternalError,
		},
	}

	graphql.AddError(ctx, gqlErr)
	return gqlErr
}