# Security & Threat Model

## Security Controls
- **OAuth 2.0 + mTLS** for lender authentication.
- **Ed25519 signing** for verifiable credentials.
- **Envelope encryption** for credential storage (KMS-managed keys).
- **WORM audit logs** for all issuance, consent, and access.
- **Rate limiting** and anomaly detection at API gateway.
- **RBAC** for internal operator access.

## Threat Model
| Threat | Mitigation |
| --- | --- |
| Credential tampering | Ed25519 signatures + hash registry | 
| Unauthorized lender access | OAuth 2.0 + mTLS + RBAC | 
| Consent spoofing | Borrower authentication + consent artefact signatures | 
| Data exfiltration | Encryption at rest + least privilege | 
| Replay attacks | Short-lived access tokens + nonce on consent access | 
| Insider misuse | WORM audit logs + privileged access monitoring | 
| Credential leakage | Field-level scope enforcement | 

## Privacy
- No credit scoring or underwriting data stored.
- Only hashes and metadata in registry.
- Consent stored separately from credentials.

## Compliance
- DPI-aligned neutrality: no marketplace or pricing behavior.
- Data minimization: only required credential fields stored and shared.
