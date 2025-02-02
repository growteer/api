package profiles

import (
	"context"

	"github.com/growteer/api/internal/app/apperrors"
	"github.com/growteer/api/internal/authz"
	"github.com/growteer/api/pkg/web3util"
)

type Repository interface {
	Create(ctx context.Context, profile Profile) (*Profile, error)
	GetByDID(ctx context.Context, did *web3util.DID) (*Profile, error)
	Update(ctx context.Context, profile Profile) (*Profile, error)
}

type Service struct {
	authz *authz.Profiles
	repo  Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateProfile(ctx context.Context, profile Profile) (*Profile, error) {
	savedProfile, err := s.repo.Create(ctx, profile)
	if err != nil {
		return nil, apperrors.Internal{
			Code:    apperrors.ErrCodeCouldNotSaveProfile,
			Message: "could not save new profile",
			Wrapped: err,
		}
	}

	return savedProfile, nil
}

func (s *Service) GetProfile(ctx context.Context, did *web3util.DID) (*Profile, error) {
	if !s.authz.MayRead(ctx, did) {
		return nil, apperrors.NotFound{
			Message: "profile not found",
		}
	}

	return s.repo.GetByDID(ctx, did)
}

func (s *Service) UpdateProfile(ctx context.Context, did *web3util.DID, profile *Profile) (*Profile, error) {
	if !s.authz.MayUpdate(ctx, did) {
		return nil, apperrors.NotFound{
			Message: "profile not found",
		}
	}

	return s.repo.Update(ctx, *profile)
}
