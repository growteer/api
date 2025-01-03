package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.60

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/growteer/api/graph/model"
)

// GenerateNonce is the resolver for the generateNonce field.
func (r *mutationResolver) GenerateNonce(ctx context.Context, input model.NonceInput) (*model.NonceResult, error) {
	if input.Address == "" {
		return nil, fmt.Errorf("bad request")
	}

	nonce, err := r.authnService.GenerateNonce(ctx, input.Address)
	if err != nil {
		return nil, err
	}

	return &model.NonceResult{
		Nonce: nonce,
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthResult, error) {
	signatureBytes, err := base64.StdEncoding.DecodeString(input.Signature)
	if err != nil {
		return nil, err
	}

	sessionToken, refreshToken, err := r.authnService.Login(ctx, input.Address, input.Message, signatureBytes)
	if err != nil {
		return nil, err
	}

	return &model.AuthResult{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	}, nil
}

// Refresh is the resolver for the refresh field.
func (r *mutationResolver) Refresh(ctx context.Context, input *model.RefreshInput) (*model.AuthResult, error) {
	sessionToken, refreshToken, err := r.authnService.RefreshSession(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &model.AuthResult{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	}, nil
}

// Nonce is the resolver for the nonce field.
func (r *queryResolver) Nonce(ctx context.Context, address string) (*model.NonceResult, error) {
	panic(fmt.Errorf("not implemented: Nonce - nonce"))
}

// Nonces is the resolver for the nonces field.
func (r *queryResolver) Nonces(ctx context.Context) ([]*model.NonceResult, error) {
	panic(fmt.Errorf("not implemented: Nonces - nonces"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
