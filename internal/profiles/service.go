package profiles

import (
	"context"

	"github.com/growteer/api/pkg/gqlutil"
	"github.com/growteer/api/pkg/web3util"
)

type Repository interface {
	Create(ctx context.Context, profile Profile) (*Profile, error)
	GetByDID(ctx context.Context, did *web3util.DID) (*Profile, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateProfile(ctx context.Context, profile Profile) (*Profile, error) {
	savedProfile, err := s.repo.Create(ctx, profile)
	if err != nil {
		return nil, gqlutil.InternalError(ctx, "signup.profile_not_saved", err)
	}

	return savedProfile, nil
}