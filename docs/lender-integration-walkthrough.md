# Lender Integration Walkthrough

## 1. Register Lender
- Provide lender certificate for mTLS.
- Obtain OAuth 2.0 client credentials.

## 2. Discover Borrowers
```bash
curl -X GET https://api.jkcreditcommons.gov/discovery?district=Shopian&sector=Handicrafts \
  --cert lender.crt --key lender.key \
  -H "Authorization: Bearer <access_token>"
```

## 3. Initiate Consent
```bash
curl -X POST https://api.jkcreditcommons.gov/consents \
  --cert lender.crt --key lender.key \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_reference": "ref-12345",
    "scope": ["crc"],
    "purpose": "loan_evaluation",
    "duration_days": 30
  }'
```

## 4. Retrieve Credential
```bash
curl -X GET https://api.jkcreditcommons.gov/credentials/crc?consent_id=consent-abc \
  --cert lender.crt --key lender.key \
  -H "Authorization: Bearer <access_token>"
```

## 5. Verify Credential
- Validate JWT signature using issuer public key.
- Check expiry and revocation status.
