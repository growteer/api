package profiles

import (
	"context"
	"time"

	"github.com/growteer/api/pkg/web3util"
	"go.mongodb.org/mongo-driver/bson"
)

type Profile struct {
	DID string `bson:"_id"`
	FirstName string `bson:"firstName"`
	LastName string `bson:"lastName"`
	DateOfBirth time.Time `bson:"dateOfBirth"`
	PrimaryEmail string `bson:"primaryEmail"`
	Location struct {
		Country string `bson:"country"`
		PostalCode string `bson:"postalCode"`
		City string `bson:"city"`
	} `bson:"location,omitempty"`
	Website string `bson:"website,omitempty"`
	PersonalGoal string `bson:"personalGoal,omitempty"`
	About string `bson:"about,omitempty"`
	CreatedAt time.Time `bson:"createdAt"`
}

func (r *Repository) Create(ctx context.Context, profile Profile) (*Profile, error) {
	profile.CreatedAt = time.Now()
	_, err := r.profiles.InsertOne(ctx, profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *Repository) GetByDID(ctx context.Context, did *web3util.DID) (*Profile, error) {
	var result *Profile
	err := r.profiles.FindOne(ctx, bson.M{"_id": did.String()}).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}