package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

type JWTSigner struct {
	Issuer  string
	PrivKey ed25519.PrivateKey
}

func (s JWTSigner) Sign(vcPayload map[string]interface{}) (string, error) {
	header := map[string]string{
		"alg": "EdDSA",
		"typ": "JWT",
	}

	payload := map[string]interface{}{
		"iss": s.Issuer,
		"vc":  vcPayload,
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal header: %w", err)
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := strings.Join([]string{encodedHeader, encodedPayload}, ".")

	signature := ed25519.Sign(s.PrivKey, []byte(unsigned))
	encodedSig := base64.RawURLEncoding.EncodeToString(signature)

	return strings.Join([]string{unsigned, encodedSig}, "."), nil
}
