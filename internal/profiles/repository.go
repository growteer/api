package profiles

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const db_collection_profiles = "profiles"

type repository struct {
	profiles *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	repo := &repository{
		profiles: db.Collection(db_collection_profiles),
	}

	return repo
}
