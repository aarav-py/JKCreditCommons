# Pilot Deployment Guide

## Prerequisites
- Kubernetes cluster (cloud-agnostic)
- PostgreSQL (managed or in-cluster)
- Object storage for encrypted credential vault
- KMS for envelope encryption
- TLS certificates for mTLS

## Services
- `issuer-service`
- `consent-service`
- `registry-service`
- `api-gateway`
- `audit-service`

## Deployment Steps
1. Provision PostgreSQL and initialize schemas.
2. Configure KMS and encryption keys.
3. Deploy services with environment variables:
   - `DATABASE_URL`
   - `KMS_KEY_ID`
   - `OAUTH_ISSUER_URL`
   - `MTLS_TRUST_STORE`
4. Configure API Gateway routes for OpenAPI endpoints.
5. Enable WORM storage for audit logs.
6. Apply network policies for service isolation.
7. Run smoke tests on discovery, consent, and credential retrieval.

## Observability
- Centralized logging (ELK/Opensearch)
- Metrics (Prometheus/Grafana)
- Audit log retention per policy
