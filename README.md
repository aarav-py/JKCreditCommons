# JK Credit Commons

Pilot-grade Digital Public Infrastructure (DPI) backend blueprint for consented, verifiable MSME credentials in Jammu & Kashmir.

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
