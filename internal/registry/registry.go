package registry

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashCredential(payload []byte) string {
	hash := sha256.Sum256(payload)
	return hex.EncodeToString(hash[:])
}
