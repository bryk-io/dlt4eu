package model

import (
	"encoding/base64"
	"time"

	"go.bryk.io/x/ccg/did"
)

// Proof add authentication and integrity protection to linked data documents through
// the use of mathematical algorithms.
// https://w3c-ccg.github.io/ld-proofs/
type Proof struct {
	src *did.ProofLD

	// A URI that identifies the digital proof suite used.
	Kind string `json:"kind"`

	// A link to a machine-readable object, such as a DID Document, that contains
	// authorization relations that explicitly permit the use of certain verification
	// methods for specific purposes. For example, a controller object could contain
	// statements that restrict a public key to being used only for signing Verifiable
	// Credentials and no other kinds of documents.
	Controller string `json:"controller"`

	// A string value that specifies the operational domain of a digital proof.
	// This may be an Internet domain name like "example.com", a ad-hoc value such
	// as "corp-level3-access", or a very specific transaction value like "8zF6T$mqP".
	// A signer may include a domain in its digital proof to restrict its use to
	// particular target, identified by the specified domain.
	Domain string `json:"domain"`

	// A string value that is included in the digital proof and MUST only be used
	// once for a particular domain and window of time. This value is used to mitigate
	// replay attacks.
	Nonce string `json:"nonce"`

	// The specific intent for the proof, the reason why an entity created it.
	// Acts as a safeguard to prevent the proof from being misused for a purpose
	// other than the one it was intended for. For example, a proof can be used
	// for purposes of authentication, for asserting control of a Verifiable
	// Credential (assertionMethod), and several others.
	//
	// Common values include: authentication, assertionMethod, keyAgreement,
	// capabilityInvocation, capabilityDelegation.
	//
	// https://w3c-ccg.github.io/ld-proofs/#proof-purpose
	Purpose string `json:"purpose"`

	// A set of parameters required to independently verify the proof, such as
	// an identifier for a public/private key pair that would be used in the
	// proof.
	VerificationMethod string `json:"verificationMethod"`

	// A random or pseudo-random value used by some authentication protocols to
	// mitigate replay attacks. Optional.
	Challenge string `json:"challenge"`
}

// Created returns the proof creation's date.
func (el *Proof) Created(format DateFormat) string {
	t, err := time.Parse(time.RFC3339, el.src.Created)
	if err != nil {
		return ""
	}
	return formatDate(t.Unix(), format)
}

// Document produces an RDF dataset on the proof's JSON-LD document, the algorithm used
// is "URDNA2015" and the format "application/n-quads".
// https://json-ld.github.io/normalization/spec
func (el *Proof) Document() string {
	ld, err := el.src.NormalizedLD()
	if err != nil {
		return ""
	}
	return string(ld)
}

// Value generated for the proof, encoded en base64 as defined in RFC 4648.
func (el *Proof) Value() string {
	return base64.StdEncoding.EncodeToString(el.src.Value)
}

func (el *Proof) withSource(src *did.ProofLD) {
	el.Kind = src.Type
	el.Controller = src.Controller
	el.Domain = src.Domain
	el.Nonce = src.Nonce
	el.Purpose = src.Purpose
	el.VerificationMethod = src.VerificationMethod
	el.Challenge = src.Challenge
	el.src = src
}
