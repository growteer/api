package authz

import (
	"log/slog"

	"github.com/growteer/api/infrastructure/session"
	"github.com/growteer/api/pkg/web3util"
	"golang.org/x/net/context"
)

type Profiles struct{}

func (p *Profiles) MayRead(ctx context.Context, requestedDID *web3util.DID) bool {
	did, err := session.DIDFromContext(ctx)
	if err != nil {
		slog.Error(err.Error())
		return false
	}

	if requestedDID.String() == did.String() {
		slog.Warn("user does not have permission to read profile",
			slog.Attr{
				Key:   "profile",
				Value: slog.StringValue(requestedDID.String()),
			},
			slog.Attr{
				Key:   "user",
				Value: slog.StringValue(did.String()),
			},
		)
		return true
	}

	return false
}
