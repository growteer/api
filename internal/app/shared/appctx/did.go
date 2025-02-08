package appctx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/growteer/api/pkg/web3util"
)

type contextKey struct {
	Name string
}

var ctxKeyUserDID = &contextKey{"userDID"}

func ContextWithDID(ctx context.Context, did string) context.Context {
	return context.WithValue(ctx, ctxKeyUserDID, did)
}

func DIDFromContext(ctx context.Context) (*web3util.DID, error) {
	rawDid, ok := ctx.Value(ctxKeyUserDID).(string)
	if !ok {
		err := fmt.Errorf("no did found in context")
		return nil, err
	}

	did, err := web3util.DIDFromString(rawDid)
	if err != nil {
		slog.Error(err.Error(), slog.Attr{
			Key:   "did",
			Value: slog.StringValue(rawDid),
		})

		return nil, err
	}

	return did, nil
}
