#!/usr/bin/env bash
set -euo pipefail

# helper to deploy Keycloak and expose a shared issuer URL via Traefik.
# Defaults:
#   KEYCLOAK_HOST=keycloak.xfsc.local
#   KEYCLOAK_NAMESPACE=keycloak
#   TLS_CERT_FILE=./certs/dev.crt
#   TLS_KEY_FILE=./certs/dev.key

KEYCLOAK_HOST="${KEYCLOAK_HOST:-keycloak.xfsc.local}"
KEYCLOAK_NAMESPACE="${KEYCLOAK_NAMESPACE:-keycloak}"
TLS_CERT_FILE="${TLS_CERT_FILE:-./certs/dev.crt}"
TLS_KEY_FILE="${TLS_KEY_FILE:-./certs/dev.key}"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

log "ℹ️ Creating namespace: $KEYCLOAK_NAMESPACE"
kubectl create namespace "$KEYCLOAK_NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

log "ℹ️ Deploying Keycloak (dev) into namespace: $KEYCLOAK_NAMESPACE"
kubectl create -f https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/refs/heads/main/kubernetes/keycloak.yaml -n "$KEYCLOAK_NAMESPACE" 2>/dev/null || kubectl apply -f https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/refs/heads/main/kubernetes/keycloak.yaml -n "$KEYCLOAK_NAMESPACE"

log "ℹ️ Creating TLS secret for Keycloak"
if [[ ! -f "$TLS_CERT_FILE" ]] || [[ ! -f "$TLS_KEY_FILE" ]]; then
  log "❌ TLS cert or key file not found: $TLS_CERT_FILE, $TLS_KEY_FILE"
  exit 1
fi
kubectl create secret tls dev-wildcard-tls \
  --cert="$TLS_CERT_FILE" \
  --key="$TLS_KEY_FILE" \
  -n "$KEYCLOAK_NAMESPACE" \
  --dry-run=client -o yaml | kubectl apply -f -
log "✅ TLS secret created or updated"

log "ℹ️ Applying Traefik ingress for Keycloak host: $KEYCLOAK_HOST"
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: keycloak-ingress
  namespace: ${KEYCLOAK_NAMESPACE}
  annotations:
    kubernetes.io/ingress.allow-http: "true"
    traefik.ingress.kubernetes.io/router.entrypoints: websecure
spec:
  ingressClassName: traefik
  tls:
  - hosts:
    - ${KEYCLOAK_HOST}
    secretName: dev-wildcard-tls
  rules:
  - host: ${KEYCLOAK_HOST}
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: keycloak
            port:
              number: 8080
EOF

log "✅ Keycloak ingress applied with TLS"
log "ℹ️ Add host mapping on your machine:"
log "   <your-dev-ip> ${KEYCLOAK_HOST}"
log "ℹ️ Use issuer URL:"
log "   https://${KEYCLOAK_HOST}/realms/dcs"
