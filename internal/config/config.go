package config

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Addr               string
	Issuer             string
	AuditLogPath       string
	CredentialKey      []byte
	Ed25519PrivateKey  ed25519.PrivateKey
	Ed25519PublicKey   ed25519.PublicKey
	ConsentTTL         time.Duration
	OAuthAudience      string
	ExpectedBearer     string
	DiscoveryPageLimit int
}

func Load() (*Config, error) {
	cfg := &Config{
		Addr:               getEnv("ADDR", ":8080"),
		Issuer:             getEnv("VC_ISSUER", "did:jk:credit-commons"),
		AuditLogPath:       getEnv("AUDIT_LOG_PATH", "audit.log"),
		OAuthAudience:      getEnv("OAUTH_AUDIENCE", "jk-credit-commons"),
		ExpectedBearer:     getEnv("EXPECTED_BEARER", ""),
		DiscoveryPageLimit: getEnvInt("DISCOVERY_PAGE_LIMIT", 100),
	}

	cfg.ConsentTTL = getEnvDuration("CONSENT_TTL", 30*24*time.Hour)

	key, err := loadCredentialKey()
	if err != nil {
		return nil, err
	}
	cfg.CredentialKey = key

	priv, pub, err := loadEd25519Keys()
	if err != nil {
		return nil, err
	}
	cfg.Ed25519PrivateKey = priv
	cfg.Ed25519PublicKey = pub

	return cfg, nil
}

func loadCredentialKey() ([]byte, error) {
	keyB64 := os.Getenv("CREDENTIAL_KEY")
	if keyB64 == "" {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, fmt.Errorf("generate credential key: %w", err)
		}
		return key, nil
	}
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return nil, fmt.Errorf("decode credential key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("credential key must be 32 bytes")
	}
	return key, nil
}

func loadEd25519Keys() (ed25519.PrivateKey, ed25519.PublicKey, error) {
	privB64 := os.Getenv("ED25519_PRIVATE_KEY")
	if privB64 == "" {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("generate ed25519 key: %w", err)
		}
		return priv, pub, nil
	}
	privBytes, err := base64.StdEncoding.DecodeString(privB64)
	if err != nil {
		return nil, nil, fmt.Errorf("decode ed25519 private key: %w", err)
	}
	if len(privBytes) != ed25519.PrivateKeySize {
		return nil, nil, fmt.Errorf("ed25519 private key must be %d bytes", ed25519.PrivateKeySize)
	}
	priv := ed25519.PrivateKey(privBytes)
	pub := priv.Public().(ed25519.PublicKey)
	return priv, pub, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		var parsed int
		_, err := fmt.Sscanf(value, "%d", &parsed)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
