package model

import (
	"errors"
	"time"

	"go.bryk.io/x/ccg/did"
)

// Identifier provides a DID instance based on the specification v1.0
type Identifier struct {
	id      *did.Identifier
	created int64
	updated int64
}

// NewIdentifier properly initialize a new DID instance with a single "master" key.
func NewIdentifier() (*Identifier, error) {
	id, err := did.NewIdentifierWithMode("dlt4eu", "", did.ModeUUID)
	if err != nil {
		return nil, err
	}
	if err := id.AddNewKey("master", did.KeyTypeEd, did.EncodingBase64); err != nil {
		return nil, err
	}
	if err := id.AddVerificationMethod(id.GetReference("master"), did.AuthenticationVM); err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	return &Identifier{
		id:      id,
		created: now,
		updated: now,
	}, nil
}

// LoadIdentifier prepares an identifier instance from a provided DID document.
func LoadIdentifier(doc *did.Document) (*Identifier, error) {
	id, err := did.FromDocument(doc)
	if err != nil {
		return nil, errors.New("invalid DID document")
	}
	now := time.Now()
	created, err := time.Parse(time.RFC3339, doc.Created)
	if err != nil {
		created = now
	}
	return &Identifier{
		id:      id,
		created: created.Unix(),
		updated: now.Unix(),
	}, nil
}

// ID returns a textual representation of the DID instance.
func (el *Identifier) ID() string {
	return el.id.String()
}

// Created returns the date of the identifier's original creation.
func (el *Identifier) Created(format DateFormat) string {
	return formatDate(el.created, format)
}

// Updated returns the date of the identifier's last update operation.
func (el *Identifier) Updated(format DateFormat) string {
	return formatDate(el.updated, format)
}

// AuthenticationMethods enabled for the identifier.
func (el *Identifier) AuthenticationMethods() []string {
	return el.id.GetVerificationMethod(did.AuthenticationVM)
}

// Keys returns the registered keys on the identifier.
func (el *Identifier) Keys() []*PublicKey {
	list := make([]*PublicKey, len(el.id.Keys()))
	for i, k := range el.id.Keys() {
		list[i] = &PublicKey{
			ID:         k.ID,
			Kind:       k.Type.String(),
			Controller: k.Controller,
			Value:      k.ValueBase64,
		}
	}
	return list
}

// Document returns the DID document for the identifier instance. The returned document
// remove any private key material present, making the document safe to be published and
// shared.
func (el *Identifier) Document(mode DocumentMode) string {
	var ld []byte
	var err error
	switch mode {
	case DocumentModeNormalized:
		ld, err = el.id.Document(true).NormalizedLD()
	case DocumentModeExpanded:
		ld, err = el.id.Document(true).ExpandedLD()
	default:
		return ""
	}
	if err != nil {
		return ""
	}
	return string(ld)
}

// ProduceProof will generate a valid linked data proof for the provided data.
// https://w3c-dvcg.github.io/ld-proofs
func (el *Identifier) ProduceProof(req *ProofRequest) (*Proof, error) {
	key := el.id.Key("master")
	src, err := key.ProduceProof([]byte(req.Data), req.Purpose, req.Domain)
	if err != nil {
		return nil, err
	}
	p := &Proof{}
	p.withSource(src)
	p.Controller = el.ID()
	return p, nil
}
