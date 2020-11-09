// +build js,wasm

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	"go.bryk.io/x/ccg/did"
)

// Restore DID instance from its document
func loadDID(contents string) (*did.Identifier, error) {
	// Get DID from document
	doc := &did.Document{}
	if err := json.Unmarshal([]byte(contents), doc); err != nil {
		return nil, errors.New("invalid DID document")
	}
	id, err := did.FromDocument(doc)
	if err != nil {
		return nil, errors.New("invalid DID document")
	}
	return id, nil
}

// Return a properly encoded error message
func encodeError(err error) interface{} {
	msg := map[string]string{
		"error": err.Error(),
	}
	output, _ := json.Marshal(msg)
	return js.ValueOf(fmt.Sprintf("%s", output)).String()
}

// Create a new DID instance and return its complete
// JSON-encoded document.
func CreateDID(this js.Value, _ []js.Value) interface{} {
	// Generate DID of requested method, add an Ed25519 master
	// key authentication key and prepare a document proof
	id, err := did.NewIdentifierWithMode("dlt4eu", "", did.ModeUUID)
	if err != nil {
		return encodeError(err)
	}
	if err = id.AddNewKey("master", did.KeyTypeEd, did.EncodingBase64); err != nil {
		return encodeError(err)
	}
	if err = id.AddVerificationMethod(id.GetReference("master"), did.AuthenticationVM); err != nil {
		return encodeError(err)
	}

	// Return JSON-encoded document with and without private keys
	result := map[string]interface{}{}
	result["withPrivateKeys"] = id.Document(false)
	result["withoutPrivateKeys"] = id.Document(true)
	output, _ := json.MarshalIndent(result, "", "  ")
	return js.ValueOf(fmt.Sprintf("%s", output)).String()
}

// Generate a signature LD document.
// Parameters:
// - JSON-encoded did document (string)
// - contents to sign (string)
// - domain value (string)
func CreateProof(this js.Value, args []js.Value) interface{} {
	// Get parameters
	if len(args) != 3 {
		return encodeError(errors.New("missing required parameters"))
	}
	doc := args[0].String()
	data := args[1].String()
	domain := args[2].String()

	// Get DID from document
	id, err := loadDID(doc)
	if err != nil {
		return encodeError(err)
	}

	// Produce proof
	key := id.Key("master")
	proof, err := key.ProduceProof([]byte(data), "authentication", domain)
	if err != nil {
		return encodeError(err)
	}

	// Return JSON-encoded proof document
	output, _ := json.MarshalIndent(proof, "", "  ")
	return js.ValueOf(fmt.Sprintf("%s", output)).String()
}

func main() {
	// Register "exported" methods
	js.Global().Set("createDID", js.FuncOf(CreateDID))
	js.Global().Set("createProof", js.FuncOf(CreateProof))

	// Block and prevent program to exit
	select {}
}
