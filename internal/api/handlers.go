package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"jkcreditcommons/internal/audit"
	"jkcreditcommons/internal/crypto"
	"jkcreditcommons/internal/registry"
	"jkcreditcommons/internal/store"
	"jkcreditcommons/internal/vc"
)

type Service struct {
	Store          store.Store
	Audit          *audit.Logger
	Signer         crypto.JWTSigner
	Envelope       *crypto.Envelope
	Issuer         string
	ConsentTTL     time.Duration
	DiscoveryLimit int
}

type DiscoveryResponse struct {
	Results []store.BorrowerRecord `json:"results"`
}

type ConsentRequest struct {
	BorrowerReference string   `json:"borrower_reference"`
	Scope             []string `json:"scope"`
	Purpose           string   `json:"purpose"`
	DurationDays      int      `json:"duration_days"`
	AssistedBy        string   `json:"assisted_by,omitempty"`
}

type ConsentResponse struct {
	ConsentID string    `json:"consent_id"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CredentialResponse struct {
	JWTVc string `json:"jwt_vc"`
}

func (s *Service) Discovery(w http.ResponseWriter, r *http.Request) {
	district := r.URL.Query().Get("district")
	sector := r.URL.Query().Get("sector")
	results, err := s.Store.ListBorrowers(district, sector, s.DiscoveryLimit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "discovery_failed")
		return
	}
	_ = s.Audit.Append(audit.Event{
		Timestamp: time.Now().UTC(),
		Action:    "discovery",
		Actor:     actorFromRequest(r),
		Metadata:  map[string]string{"district": district, "sector": sector},
	})
	respondJSON(w, http.StatusOK, DiscoveryResponse{Results: results})
}

func (s *Service) CreateConsent(w http.ResponseWriter, r *http.Request) {
	var req ConsentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.BorrowerReference == "" || len(req.Scope) == 0 || req.Purpose == "" || req.DurationDays <= 0 {
		respondError(w, http.StatusBadRequest, "missing_fields")
		return
	}
	consentID := "consent-" + randomID()
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(time.Duration(req.DurationDays) * 24 * time.Hour)

	consent := store.Consent{
		ID:          consentID,
		BorrowerRef: req.BorrowerReference,
		LenderID:    actorFromRequest(r),
		Scope:       req.Scope,
		Purpose:     req.Purpose,
		AssistedBy:  req.AssistedBy,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
		Status:      "active",
	}
	if err := s.Store.SaveConsent(consent); err != nil {
		respondError(w, http.StatusInternalServerError, "consent_store_failed")
		return
	}
	_ = s.Audit.Append(audit.Event{
		Timestamp: time.Now().UTC(),
		Action:    "consent_created",
		Actor:     consent.LenderID,
		Metadata:  map[string]string{"consent_id": consentID},
	})
	respondJSON(w, http.StatusCreated, ConsentResponse{ConsentID: consentID, Status: consent.Status, ExpiresAt: expiresAt})
}

func (s *Service) GetCredential(w http.ResponseWriter, r *http.Request) {
	consentID := r.URL.Query().Get("consent_id")
	if consentID == "" {
		respondError(w, http.StatusBadRequest, "missing_consent_id")
		return
	}
	consent, err := s.Store.GetConsent(consentID)
	if err != nil {
		respondError(w, http.StatusNotFound, "consent_not_found")
		return
	}
	if consent.Status != "active" || time.Now().UTC().After(consent.ExpiresAt) {
		respondError(w, http.StatusForbidden, "consent_inactive")
		return
	}
	if actor := actorFromRequest(r); actor != "" && actor != consent.LenderID {
		respondError(w, http.StatusForbidden, "consent_lender_mismatch")
		return
	}

	credential, err := s.Store.GetCredential(consent.BorrowerRef)
	if err != nil {
		respondError(w, http.StatusNotFound, "credential_not_found")
		return
	}
	if credential.Revoked || time.Now().UTC().After(credential.ExpiresAt) {
		respondError(w, http.StatusGone, "credential_revoked_or_expired")
		return
	}
	decrypted, err := s.Envelope.Decrypt(credential.EncryptedVC)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "credential_decrypt_failed")
		return
	}
	var vcPayload map[string]interface{}
	if err := json.Unmarshal(decrypted, &vcPayload); err != nil {
		respondError(w, http.StatusInternalServerError, "credential_decode_failed")
		return
	}
	jwt, err := s.Signer.Sign(vcPayload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "credential_sign_failed")
		return
	}
	_ = s.Audit.Append(audit.Event{
		Timestamp: time.Now().UTC(),
		Action:    "credential_retrieved",
		Actor:     consent.LenderID,
		Metadata:  map[string]string{"consent_id": consentID, "borrower_reference": consent.BorrowerRef},
	})
	respondJSON(w, http.StatusOK, CredentialResponse{JWTVc: jwt})
}

func (s *Service) IssueCredential(w http.ResponseWriter, r *http.Request) {
	var subject vc.CreditReadinessSubject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	borrowerRef := r.URL.Query().Get("borrower_reference")
	if borrowerRef == "" {
		respondError(w, http.StatusBadRequest, "missing_borrower_reference")
		return
	}

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(365 * 24 * time.Hour)

	credential := vc.CreditReadinessCredential{
		Context:        []string{"https://www.w3.org/2018/credentials/v1"},
		Type:           []string{"VerifiableCredential", "CreditReadinessCredential"},
		Issuer:         s.Issuer,
		IssuanceDate:   issuedAt,
		ExpirationDate: expiresAt,
		CredentialSub:  subject,
	}

	payload, err := credential.Marshal()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "credential_marshal_failed")
		return
	}
	encrypted, err := s.Envelope.Encrypt(payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "credential_encrypt_failed")
		return
	}
	record := store.CredentialRecord{
		BorrowerRef: borrowerRef,
		EncryptedVC: encrypted,
		Hash:        registry.HashCredential(payload),
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
		RegistryNotes: map[string]string{
			"issuer": s.Issuer,
		},
	}
	if err := s.Store.SaveCredential(record); err != nil {
		respondError(w, http.StatusInternalServerError, "credential_store_failed")
		return
	}
	_ = s.Audit.Append(audit.Event{
		Timestamp: time.Now().UTC(),
		Action:    "credential_issued",
		Actor:     actorFromRequest(r),
		Metadata:  map[string]string{"borrower_reference": borrowerRef, "hash": record.Hash},
	})
	respondJSON(w, http.StatusCreated, map[string]string{"hash": record.Hash})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, code string) {
	respondJSON(w, status, map[string]string{"error": code})
}

func actorFromRequest(r *http.Request) string {
	return r.Header.Get("X-Lender-Id")
}

func randomID() string {
	return strings.ReplaceAll(time.Now().UTC().Format("20060102150405.000000"), ".", "")
}

var errUnauthorized = errors.New("unauthorized")

func RequireBearer(expected string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expected == "" {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondError(w, http.StatusUnauthorized, "missing_bearer")
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != expected {
			respondError(w, http.StatusUnauthorized, "invalid_bearer")
			return
		}
		next.ServeHTTP(w, r)
	})
}
