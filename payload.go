package copilot

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
)

const (
	SseEventNameConfirmation string = "copilot_confirmation"
	SseEventNameReferences   string = "copilot_references"
	SseEventNameErrors       string = "copilot_errors"
)

// PayloadVerifier is for verifying a payload using ECDSA to ensure it's from GitHub.
type PayloadVerifier interface {
	Verify(body []byte, sig string) (bool, error)
}

// NewVerifier returns a new Verifier.
func NewPayloadVerifier() (PayloadVerifier, error) {
	k, err := fetchPublicKey()
	if err != nil {
		return nil, err
	}
	return &verifier{k}, nil
}

// NewPayloadVerifierWithKey returns a new Verifier with the given public key.
func NewPayloadVerifierWithKey(pubKey string) (PayloadVerifier, error) {
	k, err := parsePubKey(pubKey)
	if err != nil {
		return nil, err
	}
	return &verifier{k}, nil
}

type verifier struct {
	pubKey *ecdsa.PublicKey
}

// Verify checks if the payload is valid.
func (a *verifier) Verify(data []byte, sig string) (bool, error) {
	// Parse the Signature
	parsedSig := asn1Signature{}
	asnSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false, err
	}
	rest, err := asn1.Unmarshal(asnSig, &parsedSig)
	if err != nil || len(rest) != 0 {
		return false, err
	}

	// Verify the SHA256 encoded payload against the signature with GitHub's Key
	digest := sha256.Sum256(data)
	keyOk := ecdsa.Verify(a.pubKey, digest[:], parsedSig.R, parsedSig.S)

	return keyOk, nil
}

func parsePubKey(pubKey string) (*ecdsa.PublicKey, error) {
	pubPemStr := strings.ReplaceAll(pubKey, "\\n", "\n")
	// Decode the Public Key
	block, _ := pem.Decode([]byte(pubPemStr))
	if block == nil {
		return nil, errors.New("error parsing PEM block with GitHub public key")
	}

	// Create our ECDSA Public Key
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Because of documentation, we know it's a *ecdsa.PublicKey
	ecdsaKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("GitHub key is not ECDSA")
	}

	return ecdsaKey, nil
}

// fetchPublicKey fetches the keys used to sign messages from copilot.
// Checking the signature with one of these keys verifies that the request
// came from GitHub and not elsewhere on the internet.
func fetchPublicKey() (*ecdsa.PublicKey, error) {
	resp, err := http.Get("https://api.github.com/meta/public_keys/copilot_api")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch public key: %s", resp.Status)
	}

	var respBody struct {
		PublicKeys []struct {
			Identifier string `json:"key_identifier"`
			Key        string `json:"key"`
			IsCurrent  bool   `json:"is_current"`
		} `json:"public_keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	var rawKey string
	for _, pk := range respBody.PublicKeys {
		if pk.IsCurrent {
			rawKey = pk.Key
			break
		}
	}
	if rawKey == "" {
		return nil, fmt.Errorf("could not find current public key")
	}

	return parsePubKey(rawKey)
}

// asn1Signature is a struct for ASN.1 serializing/parsing signatures.
type asn1Signature struct {
	R *big.Int
	S *big.Int
}
