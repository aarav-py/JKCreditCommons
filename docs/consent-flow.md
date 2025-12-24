# Consent Flow

## Principles
- Borrower-controlled consent with explicit scope, duration, and lender identity.
- Consent artefacts are stored separately from credentials.
- Assisted consent supported via CSC/SHG representatives.

## Actors
- **Borrower (MSME)**
- **Assisted Representative** (CSC/SHG)
- **Lender**
- **Consent Service**

## Consent Artefact
Captured as a signed record:
- `consent_id`
- `borrower_reference`
- `lender_id`
- `scope` (credential types, fields)
- `purpose`
- `duration` (start/end)
- `assisted_by` (optional representative id)
- `issued_at`
- `status` (active/revoked/expired)

## Flow
1. **Discovery**: Lender queries anonymised discovery API.
2. **Consent Request**: Lender initiates consent request for a borrower reference.
3. **Borrower Approval**:
   - Self-service: borrower authenticates and approves.
   - Assisted: CSC/SHG rep captures borrower approval and identity proof.
4. **Consent Artefact Issued**: Consent Service stores artefact and returns consent token.
5. **Credential Access**: Lender presents consent token to retrieve credential.
6. **Revocation**: Borrower or regulator can revoke consent at any time.

## Consent Storage Separation
Consent artefacts are stored in a dedicated table and encrypted separately from credential storage, with a strict access policy.
