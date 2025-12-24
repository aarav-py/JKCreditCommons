# JK Credit Commons

JK Credit Commons is a state-owned, lender-agnostic Digital Public Infrastructure (DPI) that converts existing government and community data into consent-based **credit readiness credentials** for micro and nano enterprises in Jammu & Kashmir. It certifies observable facts such as enterprise existence, continuity of activity, and community validation. The platform does not compute credit scores, intermediate loans, or influence pricing, ensuring all credit decisions remain market-driven.

## Quickstart
```bash
go run ./cmd/server
```

Environment configuration:
- `ADDR`: bind address (default `:8080`)
- `VC_ISSUER`: issuer DID (default `did:jk:credit-commons`)
- `AUDIT_LOG_PATH`: audit log file (default `audit.log`)
- `CREDENTIAL_KEY`: base64 32-byte AES key for encrypted credential storage
- `ED25519_PRIVATE_KEY`: base64 Ed25519 private key for signing JWT VCs
- `EXPECTED_BEARER`: expected bearer token for OAuth placeholder validation

## Usage (Local Pilot Demo)
Issue a credential for a borrower reference:
```bash
curl -X POST "http://localhost:8080/admin/issue/crc?borrower_reference=ref-001" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "enterprise_exists": true,
    "months_active": 18,
    "validation_source": "Mission Yuva",
    "sector": "Handicrafts",
    "district": "Shopian"
  }'
```

Discover anonymised borrower references:
```bash
curl "http://localhost:8080/discovery?district=Shopian" \
  -H "Authorization: Bearer <token>" \
  -H "X-Lender-Id: lender-001"
```

Initiate consent (self or assisted):
```bash
curl -X POST "http://localhost:8080/consents" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -H "X-Lender-Id: lender-001" \
  -d '{
    "borrower_reference": "ref-001",
    "scope": ["crc"],
    "purpose": "loan_evaluation",
    "duration_days": 30,
    "assisted_by": "csc-operator-9"
  }'
```

Retrieve the credential (post-consent):
```bash
curl "http://localhost:8080/credentials/crc?consent_id=consent-<id>" \
  -H "Authorization: Bearer <token>" \
  -H "X-Lender-Id: lender-001"
```

## Contents
- `cmd/server`: reference API server (stateless, in-memory storage for pilot)
- `internal`: core services (consent, credential issuance, registry hashing)
- `openapi/openapi.yaml`: Lender-facing API specification
- `schemas/credit-readiness-credential.json`: W3C Verifiable Credential (JWT VC) schema
- `docs/architecture.md`: System architecture & components
- `docs/consent-flow.md`: Borrower consent flow & artefacts
- `docs/security-threat-model.md`: Security controls & threat model
- `docs/pilot-deployment-guide.md`: K8s-ready deployment steps
- `docs/lender-integration-walkthrough.md`: Example lender integration

## Non-goals
The platform does **not** compute credit scores, underwrite, price loans, disburse funds, rank lenders/borrowers, or operate as a marketplace.
