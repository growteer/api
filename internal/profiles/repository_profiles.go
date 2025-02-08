package profiles

import (
	"context"
	"log/slog"
	"time"

	"github.com/growteer/api/internal/app/apperrors"
	"github.com/growteer/api/pkg/web3util"
	"go.mongodb.org/mongo-driver/bson"
)

type Profile struct {
	DID          string    `bson:"_id"`
	FirstName    string    `bson:"firstName"`
	LastName     string    `bson:"lastName"`
	DateOfBirth  time.Time `bson:"dateOfBirth"`
	PrimaryEmail string    `bson:"primaryEmail"`
	Location     Location  `bson:"location,omitempty"`
	Website      string    `bson:"website,omitempty"`
	PersonalGoal string    `bson:"personalGoal,omitempty"`
	About        string    `bson:"about,omitempty"`
	CreatedAt    time.Time `bson:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt"`
}

type Location struct {
	Country    string `bson:"country"`
	PostalCode string `bson:"postalCode"`
	City       string `bson:"city"`
}

func (r *repository) Exists(ctx context.Context, did *web3util.DID) bool {
	var result Profile
	if err := r.profiles.FindOne(ctx, bson.M{"_id": did.String()}).Decode(&result); err != nil {
		slog.Debug("profile not found", slog.Attr{
			Key:   "did",
			Value: slog.StringValue(did.String()),
		})

		return false
	}

	return true
}

func (r *repository) Create(ctx context.Context, profile Profile) (*Profile, error) {
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = profile.CreatedAt

	_, err := r.profiles.InsertOne(ctx, profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *repository) GetByDID(ctx context.Context, did *web3util.DID) (*Profile, error) {
	var result Profile
	err := r.profiles.FindOne(ctx, bson.M{"_id": did.String()}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Update(ctx context.Context, profile Profile) (*Profile, error) {
	profile.UpdatedAt = time.Now()
	result, err := r.profiles.ReplaceOne(ctx, bson.M{"_id": profile.DID}, profile)
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

	return &profile, nil
}
