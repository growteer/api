package authn

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/growteer/api/internal/app/shared/appctx"
	"github.com/growteer/api/internal/app/shared/apperrors"
	"github.com/growteer/api/pkg/web3util"
)

func UserSessionMiddleware(provider TokenProvider) func(http.Handler) http.Handler {
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

				ctx := appctx.ContextWithDID(r.Context(), claims.Subject)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetAuthenticatedDID(ctx context.Context) (*web3util.DID, error) {
	did, err := appctx.DIDFromContext(ctx)
	if err != nil {
		return nil, apperrors.Unauthenticated{
			Message: "no did found in context",
			Wrapped: err,
		}
	}

	if err := web3util.VerifySolanaPublicKey(did.Address); err != nil {
		return nil, apperrors.Unauthenticated{
			Message: "invalid solana public key",
			Wrapped: err,
		}
	}

	return did, nil
}
