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

func BadInputError(ctx context.Context, message string, err error) error {
	return &gqlerror.Error{
		Err: err,
		Message: message,
		Path: graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeBadRequest,
		},
	}
}

func InternalError(ctx context.Context, message string, err error) error {
	return &gqlerror.Error{
		Err: err,
		Message: message,
		Path: graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"type": errTypeInternal,
		},
	}
}