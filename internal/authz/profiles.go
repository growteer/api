package authz

import (
	"log/slog"

	"github.com/growteer/api/internal/infrastructure/session"
	"github.com/growteer/api/pkg/web3util"
	"golang.org/x/net/context"
)

type Profiles struct{}

func (p *Profiles) MayRead(ctx context.Context, profileToRead *web3util.DID) bool {
	did, err := session.DIDFromContext(ctx)
	if err != nil {
		slog.Error(err.Error())
		return false
	}

	if profileToRead.String() == did.String() {
		return true
	}

	slog.Warn("user does not have permission to read profile",
		slog.Attr{
			Key:   "profile",
			Value: slog.StringValue(profileToRead.String()),
		},
		slog.Attr{
			Key:   "user",
			Value: slog.StringValue(did.String()),
		},
	)

	return false
}

func (p *Profiles) MayUpdate(ctx context.Context, profileToUpdate *web3util.DID) bool {
	did, err := session.DIDFromContext(ctx)
	if err != nil {
		slog.Error(err.Error())
		return false
	}

	if profileToUpdate.String() == did.String() {
		return true
	}

	slog.Warn("user does not have permission to update profile",
		slog.Attr{
			Key:   "profile",
			Value: slog.StringValue(profileToUpdate.String()),
		},
		slog.Attr{
			Key:   "user",
			Value: slog.StringValue(did.String()),
		},
	)

	return false
}
