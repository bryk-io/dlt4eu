package service

import (
	"encoding/base64"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/bryk-io/dlt4eu/service/server"
	"github.com/google/uuid"
	"go.bryk.io/x/jwx"
)

// Config define the settings available when configuring a service
// handler instance.
type Config struct {
	Issuer string `json:"issuer" mapstructure:"issuer"`
	Key    string `json:"key" mapstructure:"key"`
}

// Handler provides the main service operator.
type Handler struct {
	conf     *Config
	schema   gqlgen.ExecutableSchema
	resolver *Resolver
}

// New service handler instance.
func New(conf *Config) (*Handler, error) {
	// Initialize token generator
	tg, err := tokenGenerator(conf)
	if err != nil {
		return nil, err
	}

	// Get resolver
	r := &Resolver{tg: tg}
	if err := r.init(); err != nil {
		return nil, err
	}

	// Load executable schema
	cs := server.Config{
		Resolvers: r,
	}

	// Return handler instance
	return &Handler{
		conf:     conf,
		resolver: r,
		schema:   server.NewExecutableSchema(cs),
	}, nil
}

// Schema required when exposing the service via GraphQL.
func (h *Handler) Schema() gqlgen.ExecutableSchema {
	return h.schema
}

// Shutdown the instance and free resources.
func (h *Handler) Shutdown() {
	h.resolver.shutdown()
}

// AdminToken generates the required credentials to access the API as
// an administrator.
func (h *Handler) AdminToken(subject string) (string, error) {
	params := &jwx.TokenParameters{
		Subject:          subject,
		Audience:         []string{h.conf.Issuer},
		ContentType:      "dlt4eu.vc/0.1.0",
		Expiration:       "720h",
		UniqueIdentifier: uuid.New().String(),
		Method:           jwx.ES512,
		CustomClaims:     &cc{Role: "admin"},
	}
	token, err := h.resolver.tg.NewToken("master", params)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

// Initialize token generator
func tokenGenerator(conf *Config) (*jwx.Generator, error) {
	k, err := base64.RawURLEncoding.DecodeString(conf.Key)
	if err != nil {
		return nil, err
	}
	key, err := jwx.ECFromPEM(k)
	if err != nil {
		return nil, err
	}
	gen := jwx.NewGenerator(conf.Issuer)
	gen.SupportNone(false)
	if err := gen.AddKey("master", key); err != nil {
		return nil, err
	}
	return gen, nil
}

type cc struct {
	Role string `json:"role"`
}
