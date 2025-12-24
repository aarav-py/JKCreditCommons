package store

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type BorrowerRecord struct {
	Reference string
	District  string
	Sector    string
}

type Consent struct {
	ID          string
	BorrowerRef string
	LenderID    string
	Scope       []string
	Purpose     string
	AssistedBy  string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	Status      string
}

type CredentialRecord struct {
	BorrowerRef   string
	EncryptedVC   []byte
	Hash          string
	IssuedAt      time.Time
	ExpiresAt     time.Time
	Revoked       bool
	RegistryNotes map[string]string
}

type Store interface {
	ListBorrowers(district, sector string, limit int) ([]BorrowerRecord, error)
	SaveConsent(consent Consent) error
	GetConsent(id string) (Consent, error)
	SaveCredential(record CredentialRecord) error
	GetCredential(borrowerRef string) (CredentialRecord, error)
	RevokeCredential(borrowerRef string) error
}

type MemoryStore struct {
	mu          sync.RWMutex
	borrowers   []BorrowerRecord
	consents    map[string]Consent
	credentials map[string]CredentialRecord
}

func NewMemoryStore(seed []BorrowerRecord) *MemoryStore {
	return &MemoryStore{
		borrowers:   seed,
		consents:    make(map[string]Consent),
		credentials: make(map[string]CredentialRecord),
	}
}

func (m *MemoryStore) ListBorrowers(district, sector string, limit int) ([]BorrowerRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]BorrowerRecord, 0, limit)
	for _, borrower := range m.borrowers {
		if district != "" && borrower.District != district {
			continue
		}
		if sector != "" && borrower.Sector != sector {
			continue
		}
		results = append(results, borrower)
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func (m *MemoryStore) SaveConsent(consent Consent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.consents[consent.ID] = consent
	return nil
}

func (m *MemoryStore) GetConsent(id string) (Consent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	consent, ok := m.consents[id]
	if !ok {
		return Consent{}, ErrNotFound
	}
	return consent, nil
}

func (m *MemoryStore) SaveCredential(record CredentialRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.credentials[record.BorrowerRef] = record
	return nil
}

func (m *MemoryStore) GetCredential(borrowerRef string) (CredentialRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	credential, ok := m.credentials[borrowerRef]
	if !ok {
		return CredentialRecord{}, ErrNotFound
	}
	return credential, nil
}

func (m *MemoryStore) RevokeCredential(borrowerRef string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	credential, ok := m.credentials[borrowerRef]
	if !ok {
		return ErrNotFound
	}
	credential.Revoked = true
	m.credentials[borrowerRef] = credential
	return nil
}
