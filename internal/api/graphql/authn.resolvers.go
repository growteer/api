package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.60

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/growteer/api/internal/api/graphql/gqlutil"
	"github.com/growteer/api/internal/api/graphql/model"
	"github.com/growteer/api/pkg/web3util"
)

// GenerateNonce is the resolver for the generateNonce field.
func (r *mutationResolver) GenerateNonce(ctx context.Context, address string) (*model.NonceResult, error) {
	if err := web3util.VerifySolanaPublicKey(address); err != nil {
		return nil, gqlutil.BadInputError(ctx, "invalid solana address", gqlutil.ErrCodeInvalidCredentials, err)
	}

	did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, address)

	nonce, err := r.authnService.GenerateNonce(ctx, did)
	if err != nil {
		return nil, err
	}

	return &model.NonceResult{
		Nonce: nonce,
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginDetails) (*model.AuthResult, error) {
	if err := web3util.VerifySolanaPublicKey(input.Address); err != nil {
		return nil, gqlutil.BadInputError(ctx, "invalid solana address", gqlutil.ErrCodeInvalidCredentials, err)
	}

	did := web3util.NewDID(web3util.DIDMethodPKH, web3util.NamespaceSolana, input.Address)

	sessionToken, refreshToken, err := r.authnService.Login(ctx, did, input.Message, input.Signature)
	if err != nil {
		graphql.AddError(ctx, err)
	}

	return &model.AuthResult{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshSession is the resolver for the refreshSession field.
func (r *mutationResolver) RefreshSession(ctx context.Context, input model.RefreshInput) (*model.AuthResult, error) {
	sessionToken, refreshToken, err := r.authnService.RefreshSession(ctx, input.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &model.AuthResult{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
	}, nil
}
