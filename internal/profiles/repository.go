package profiles

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const db_collection_profiles = "profiles"

type Repository struct {
	profiles *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	repo := &Repository{
		profiles: db.Collection(db_collection_profiles),
	}

	return repo
}
