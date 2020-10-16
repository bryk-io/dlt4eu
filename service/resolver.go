package service

import (
	"context"
	"errors"
	"sync"

	"github.com/bryk-io/dlt4eu/service/model"
	"go.bryk.io/x/jwx"
	"go.bryk.io/x/net/middleware"
)

// Resolver provides the application's main entry point. The
// implementations required are located on "schema.resolvers.go".
type Resolver struct {
	tg    *jwx.Generator
	store map[string]*model.Identifier
	mu    sync.Mutex
}

func (r *Resolver) init() error {
	r.store = make(map[string]*model.Identifier)
	return nil
}

func (r *Resolver) authenticate(ctx context.Context) error {
	md, ok := middleware.MetadataFromContext(ctx)
	if !ok {
		return errors.New("unauthenticated")
	}
	creds := md.Get("authorization")
	if len(creds) == 0 {
		return errors.New("unauthenticated")
	}
	return r.tg.Validate(creds[0], isAdmin)
}

func (r *Resolver) addIdentifier(id *model.Identifier) {
	r.mu.Lock()
	r.store[id.ID()] = id
	r.mu.Unlock()
}

func (r *Resolver) getIdentifier(id string) *model.Identifier {
	r.mu.Lock()
	m, ok := r.store[id]
	r.mu.Unlock()
	if !ok {
		return nil
	}
	return m
}

func (r *Resolver) shutdown() {}

func isAdmin(token *jwx.Token) error {
	claims := &cc{}
	if err := token.Decode(claims); err != nil {
		return err
	}
	if claims.Role != "admin" {
		return errors.New("not admin")
	}
	return nil
}
