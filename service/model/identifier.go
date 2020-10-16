package model

import (
	"time"

	"go.bryk.io/x/ccg/did"
)

type Identifier struct {
	id      *did.Identifier
	created int64
	updated int64
}

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

func (el *Identifier) ID() string {
	return el.id.String()
}

func (el *Identifier) Created(format DateFormat) string {
	return formatDate(el.created, format)
}

func (el *Identifier) Updated(format DateFormat) string {
	return formatDate(el.updated, format)
}

func (el *Identifier) AuthenticationMethods() []string {
	return el.id.GetVerificationMethod(did.AuthenticationVM)
}

func (el *Identifier) Keys() []*PublicKey {
	var list []*PublicKey
	for _, k := range el.id.Keys() {
		list = append(list, &PublicKey{
			ID:         k.ID,
			Kind:       k.Type.String(),
			Controller: k.Controller,
			Value:      k.ValueBase64,
		})
	}
	return list
}

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
