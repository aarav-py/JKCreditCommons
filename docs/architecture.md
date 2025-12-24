# JK Credit Commons Architecture

## Overview
JK Credit Commons is a state-owned, neutral DPI that **certifies facts**, not risk. It issues verifiable credentials (VCs) about MSME readiness, manages borrower consent, and exposes lender-agnostic APIs without acting as a credit marketplace or intermediary.

## Components
1. **Issuer Service**
   - Issues Credit Readiness Credentials (CRCs) as W3C JWT VCs (Ed25519).
   - Validates source data from Mission Yuva and community verification partners.
   - Supports revocation and expiry.

2. **Consent Service**
   - Captures borrower-controlled consent (self or assisted by CSC/SHG reps).
   - Stores consent artefacts separately from credentials.
   - Enforces scope, duration, and lender identity.

3. **Registry Service**
   - Stores only hashes + metadata for credentials (no plaintext VC).
   - Append-only audit log (WORM).

4. **Credential Vault**
   - Encrypted storage of VCs.
   - Envelope encryption with KMS-managed keys.

5. **API Gateway**
   - OAuth 2.0 + mTLS, rate limiting, and audit logging.
   - Lender APIs for discovery, consent initiation, and credential retrieval.

6. **Audit & Monitoring**
   - Immutable event log for all VC issuance, consent changes, and access.
   - Real-time monitoring with alerts for anomalous access.

## Data Flow (High-Level)
1. Enterprise data validated by issuing authority.
2. Issuer Service creates CRC and signs with Ed25519.
3. Credential hash registered in Registry Service.
4. Encrypted credential stored in Vault.
5. Lender requests consent; Consent Service records and issues consent artefact.
6. Upon consent, lender retrieves credential via API Gateway.

## Stateless & Scalable
All services are stateless. Any state is persisted to PostgreSQL and encrypted blob storage. Horizontal scaling supported behind Kubernetes Service/Ingress.

## Data Retention
- Consent artefacts: retained until expiry + audit retention window.
- Credential registry: append-only hashes for auditability.
- Credential vault: encrypted, with configurable retention.
