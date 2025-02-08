package profiles

import (
	"time"

	"github.com/growteer/api/internal/entities"
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

func (p *Profile) ToEntity() *entities.Profile {
	location := entities.Location{
		Country:    p.Location.Country,
		PostalCode: p.Location.PostalCode,
		City:       p.Location.City,
	}

	return &entities.Profile{
		About:        p.About,
		DateOfBirth:  p.DateOfBirth,
		DID:          p.DID,
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Location:     location,
		PersonalGoal: p.PersonalGoal,
		PrimaryEmail: p.PrimaryEmail,
		Website:      p.Website,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func DAOFromEntity(entity *entities.Profile) *Profile {
	location := Location{
		Country:    entity.Location.Country,
		PostalCode: entity.Location.PostalCode,
		City:       entity.Location.City,
	}

	return &Profile{
		About:        entity.About,
		DateOfBirth:  entity.DateOfBirth,
		DID:          entity.DID,
		FirstName:    entity.FirstName,
		LastName:     entity.LastName,
		Location:     location,
		PersonalGoal: entity.PersonalGoal,
		PrimaryEmail: entity.PrimaryEmail,
		Website:      entity.Website,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}
