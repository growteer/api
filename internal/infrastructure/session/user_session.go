package session

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/growteer/api/internal/infrastructure/tokens"
	"github.com/growteer/api/pkg/gqlutil"
	"github.com/growteer/api/pkg/web3util"
)

type contextKey struct {
	Name string
}

var ctxKeyUserDID = &contextKey{"userDID"}

func UserSessionMiddleware(provider *tokens.Provider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			headerParts := strings.SplitN(authHeader, " ", 2)

			if len(headerParts) == 2 && strings.ToLower(headerParts[0]) == "bearer" {
				sessionToken := headerParts[1]
				claims, err := provider.ParseSessionToken(sessionToken)
				if err != nil {
					slog.Error("unable to parse session token")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), ctxKeyUserDID, claims.Subject)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetAuthenticatedDID(ctx context.Context) (*web3util.DID, error) {
	did, err := DIDFromContext(ctx)
	if err != nil {
		return nil, gqlutil.AuthenticationError(ctx, err.Error(), err)
	}

	if err := web3util.VerifySolanaPublicKey(did.Address); err != nil {
		slog.Warn("did in context did not include a valid solana public key", slog.Attr{
			Key:   "did",
			Value: slog.StringValue(did.String()),
		})

		return nil, gqlutil.AuthenticationError(ctx, err.Error(), err)
	}

	return did, nil
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
