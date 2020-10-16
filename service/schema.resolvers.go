package service

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/bryk-io/dlt4eu/service/model"
	"github.com/bryk-io/dlt4eu/service/server"
	"github.com/google/uuid"
	"go.bryk.io/x/jwx"
)

func (r *mutationResolver) NewIdentifier(ctx context.Context) (*model.Identifier, error) {
	if err := r.authenticate(ctx); err != nil {
		return nil, err
	}
	id, err := model.NewIdentifier()
	if err == nil {
		r.addIdentifier(id)
	}
	return id, err
}

func (r *mutationResolver) NewProof(ctx context.Context, req *model.ProofRequest) (*model.Proof, error) {
	if err := r.authenticate(ctx); err != nil {
		return nil, err
	}
	id := r.getIdentifier(req.ID)
	if id == nil {
		return nil, errors.New("unknown identifier")
	}
	return id.ProduceProof(req)
}

func (r *mutationResolver) NewCredential(ctx context.Context, req *model.CredentialRequest) (*model.Credential, error) {
	if err := r.authenticate(ctx); err != nil {
		return nil, err
	}
	params := &jwx.TokenParameters{
		Subject:          req.Subject,
		Audience:         req.Audience,
		ContentType:      "dlt4eu.vc/0.1.0",
		Expiration:       req.Expiration,
		NotBefore:        req.NotBefore,
		UniqueIdentifier: uuid.New().String(),
		Method:           jwx.ES512,
	}
	if strings.TrimSpace(req.Payload) != "" {
		claims := make(map[string]interface{})
		if err := json.Unmarshal([]byte(req.Payload), &claims); err == nil {
			params.CustomClaims = claims
		}
	}
	vc, err := r.tg.NewToken("master", params)
	if err != nil {
		return nil, err
	}
	return &model.Credential{
		Token: vc.String(),
	}, nil
}

func (r *queryResolver) Resolve(ctx context.Context, id string) (*model.Identifier, error) {
	if err := r.authenticate(ctx); err != nil {
		return nil, err
	}
	el := r.getIdentifier(id)
	if el == nil {
		return nil, errors.New("unknown identifier")
	}
	return el, nil
}

func (r *queryResolver) IsCredentialValid(ctx context.Context, token string) (bool, error) {
	if err := r.authenticate(ctx); err != nil {
		return false, err
	}
	params := &jwx.TokenParameters{
		ContentType: "dlt4eu.vc/0.1.0",
		Method:      jwx.ES512,
	}
	if err := r.tg.Validate(token, params.GetValidators()...); err != nil {
		return false, err
	}
	return true, nil
}

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }

// Query returns server.QueryResolver implementation.
func (r *Resolver) Query() server.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
