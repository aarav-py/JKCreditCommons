package main

import (
	"log"
	"net/http"

	"jkcreditcommons/internal/api"
	"jkcreditcommons/internal/audit"
	"jkcreditcommons/internal/config"
	"jkcreditcommons/internal/crypto"
	"jkcreditcommons/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	seed := []store.BorrowerRecord{
		{Reference: "ref-001", District: "Shopian", Sector: "Handicrafts"},
		{Reference: "ref-002", District: "Anantnag", Sector: "Agriculture"},
	}
	memoryStore := store.NewMemoryStore(seed)

	service := &api.Service{
		Store:          memoryStore,
		Audit:          audit.New(cfg.AuditLogPath),
		Signer:         crypto.JWTSigner{Issuer: cfg.Issuer, PrivKey: cfg.Ed25519PrivateKey},
		Envelope:       crypto.NewEnvelope(cfg.CredentialKey),
		Issuer:         cfg.Issuer,
		ConsentTTL:     cfg.ConsentTTL,
		DiscoveryLimit: cfg.DiscoveryPageLimit,
	}

	mux := http.NewServeMux()
	mux.Handle("/discovery", api.RequireBearer(cfg.ExpectedBearer, http.HandlerFunc(service.Discovery)))
	mux.Handle("/consents", api.RequireBearer(cfg.ExpectedBearer, http.HandlerFunc(service.CreateConsent)))
	mux.Handle("/credentials/crc", api.RequireBearer(cfg.ExpectedBearer, http.HandlerFunc(service.GetCredential)))
	mux.Handle("/admin/issue/crc", api.RequireBearer(cfg.ExpectedBearer, http.HandlerFunc(service.IssueCredential)))

	log.Printf("JK Credit Commons API listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
