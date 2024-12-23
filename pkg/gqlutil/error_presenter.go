package gqlutil

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func PresentError(ctx context.Context, err error) *gqlerror.Error {
	slog.Error(err.Error())

	return graphql.DefaultErrorPresenter(ctx, err)
}

func Recover(ctx context.Context, err interface{}) error {
	asErr, ok := err.(error)
	if !ok {
		asErr = fmt.Errorf("internal server error: %v", err)
	}

	slog.Error(asErr.Error())

	return gqlerror.Wrap(asErr)
}