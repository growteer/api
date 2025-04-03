package profiles

import (
	"context"
	"log/slog"
	"time"

	"github.com/growteer/api/internal/app/shared/apperrors"
	"github.com/growteer/api/internal/entities"
	"github.com/growteer/api/pkg/web3util"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) Exists(ctx context.Context, did *web3util.DID) bool {
	var result Profile
	if err := r.profiles.FindOne(ctx, bson.M{"did": did.String()}).Decode(&result); err != nil {
		slog.Debug("profile not found", slog.Attr{
			Key:   "did",
			Value: slog.StringValue(did.String()),
		})

		return false
	}

	return true
}

func (r *repository) Create(ctx context.Context, profile *entities.Profile) (*entities.Profile, error) {
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = profile.CreatedAt

	dao := DAOFromEntity(profile)
	_, err := r.profiles.InsertOne(ctx, dao)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *repository) GetByDID(ctx context.Context, did *web3util.DID) (*entities.Profile, error) {
	var dao Profile
	err := r.profiles.FindOne(ctx, bson.M{"did": did.String()}).Decode(&dao)
	if err != nil {
		return nil, err
	}

	return dao.ToEntity(), nil
}

func (r *repository) Update(ctx context.Context, profile *entities.Profile) (*entities.Profile, error) {
	profile.UpdatedAt = time.Now()

	dao := DAOFromEntity(profile)
	result, err := r.profiles.ReplaceOne(ctx, bson.M{"did": profile.DID}, dao)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		err := apperrors.NotFound{
			Message: "no profile found for updating",
		}

		slog.Warn(err.Error(), slog.Attr{
			Key:   "did",
			Value: slog.StringValue(profile.DID),
		})

		return nil, err
	}

	return profile, nil
}
