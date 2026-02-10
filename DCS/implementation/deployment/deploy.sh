#!/usr/bin/env bash
set -euo pipefail

#----------------------------------------
# Functions
#----------------------------------------
log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

usage() {
  echo "Usage: $0 <kubeconfig> <private_key_path> <crt_path> <domain> <path> <realm> <oidc_client_id>"
  exit 1
}

#----------------------------------------
# Input validation
#----------------------------------------
[ "$#" -ne 7 ] && usage
export KUBECONFIG="$1"
KEY_FILE="$2"
CRT_FILE="$3"
DOMAIN="$4"
URL_PATH="$5"
REALM="$6"
OIDC_CLIENT_ID="$7"

#----------------------------------------
# Image Registry (from environment)
#----------------------------------------
DOCKER_REGISTRY="${DOCKER_REGISTRY:-}"
DOCKER_REPO="${DOCKER_REPO:-}"

IMAGE_NAME="digital-contracting-service"
if [[ -n "$DOCKER_REGISTRY" && -n "$DOCKER_REPO" ]]; then
  IMAGE_NAME="$DOCKER_REGISTRY/$DOCKER_REPO/digital-contracting-service"
fi
log "ℹ️ Image name: $IMAGE_NAME"

#----------------------------------------
# OIDC Configuration
#----------------------------------------
# Use in-cluster Keycloak URL for the backend (if OIDC_ISSUER_URL not set)
# This allows the backend to reach Keycloak from inside the cluster
if [[ -z "${OIDC_ISSUER_URL:-}" ]]; then
  # Default to in-cluster service URL
  OIDC_ISSUER_URL="http://keycloak.default.svc.cluster.local:8080/auth/realms/${REALM}"
  log "ℹ️ Using in-cluster OIDC issuer URL (override with OIDC_ISSUER_URL env var)"
fi

log "ℹ️ OIDC Configuration:"
log "  - Issuer URL (for backend): $OIDC_ISSUER_URL"
log "  - Client ID: $OIDC_CLIENT_ID"
log "  - Keycloak API URL (for deploy script): $KEYCLOAK_URL"

# Check if kubeconfig file exists
if [[ ! -f "$KUBECONFIG" ]]; then
  log "❌ Kubeconfig file not found: $KUBECONFIG"
  exit 1
fi


#----------------------------------------
# Cleanup local helm artifacts
#----------------------------------------
if [ -f Chart.lock ]; then
  rm Chart.lock
  log "✅ Removed Chart.lock"
fi

if [ -d charts ]; then
  rm -rf charts
  log "✅ Removed charts/ directory"
fi

#----------------------------------------
# Check dependencies
#----------------------------------------
for cmd in kubectl helm jq curl sed trap; do
  if ! command -v "$cmd" &>/dev/null; then
    log "❌ '$cmd' is not installed. Please install it and retry."
    exit 1
  else
    log "✅ Found '$cmd'"
  fi
done

#----------------------------------------
# Verify ingress class traefik is installed
#----------------------------------------
log "ℹ Checking for ingressClass traefik..."
if ! kubectl get ingressclass traefik &>/dev/null; then
  log "❌ Ingress class traefik not found"
  exit 1
else
    log "✅ Ingress class traefik found"
fi

#----------------------------------------
# Continue with deployment
#----------------------------------------
log "ℹ️ Continuing with deployment..."

#----------------------------------------
# Generate and validate namespace from path
#----------------------------------------
NAMESPACE="digital-contracting-service-${URL_PATH}"
log "ℹ️ Using namespace: $NAMESPACE"

# Create namespace first
kubectl create namespace "$NAMESPACE" --kubeconfig "$KUBECONFIG" 2>/dev/null || true
log "✅ Namespace created or already exists"

#----------------------------------------
# Prepare temporary values file
#----------------------------------------
TMP_VALUES="$(mktemp -t values.XXXXXX.yaml)" || TMP_VALUES="/tmp/values-$$.yaml"
cp values.yaml "$TMP_VALUES"
log "ℹ️ Replacing placeholders in $TMP_VALUES"
sed -i \
  -e "s|\[domain-name\]|${DOMAIN}|g" \
  -e "s|\[path\]|${URL_PATH}|g" \
  -e "s|\[namespace\]|${NAMESPACE}|g" \
  -e "s|\[Admin_Username\]|${ADMIN_USER}|g" \
  -e "s|\[Admin_Password\]|${ADMIN_PASS}|g" \
  -e "s|\[oidc-issuer-url\]|${OIDC_ISSUER_URL}|g" \
  -e "s|\[oidc-client-id\]|${OIDC_CLIENT_ID}|g" \
  -e "s|\[registry\]|${IMAGE_NAME}|g" \
  "$TMP_VALUES"
log "✅ Placeholders replaced in $TMP_VALUES"

#----------------------------------------
# Helm dependency build & install
#----------------------------------------

log "ℹ️ Running: helm dependency build"
helm dependency build . --kubeconfig "$KUBECONFIG"

log "ℹ️ Installing digital-contracting-service via Helm"
helm install digital-contracting-service . \
  --namespace "$NAMESPACE" \
  --create-namespace \
  --kubeconfig "$KUBECONFIG" \
  -f "$TMP_VALUES"
log "✅ digital-contracting-service Helm release deployed"

#----------------------------------------
# Create TLS secret
#----------------------------------------
log "ℹ️ Creating TLS secret 'certificates'"
kubectl create secret tls certificates \
  --namespace "$NAMESPACE" \
  --key "$KEY_FILE" \
  --cert "$CRT_FILE" \
  --kubeconfig "$KUBECONFIG"
log "✅ TLS secret created"

#----------------------------------------
# Wait for Deployment to be ready
#----------------------------------------
log "ℹ️ Waiting for digital-contracting-service deployment to be ready (max 2m)..."
if ! kubectl rollout status deployment/digital-contracting-service \
     -n "$NAMESPACE" \
     --timeout=300s \
     --kubeconfig "$KUBECONFIG"; then
  log "❌ Timeout waiting for digital-contracting-service. Pod statuses:"
  kubectl get pods -n "$NAMESPACE" -o wide --kubeconfig "$KUBECONFIG"
  exit 1
fi
log "✅ digital-contracting-service deployment is ready"

#----------------------------------------
# Final output
#----------------------------------------
log "🎉 All operations completed successfully!"
echo
echo "🔹 DCS URL: https://${DOMAIN}/${URL_PATH}"
echo ""
log "ℹ️ Before accessing the service, ensure Keycloak is configured:"
log "   1. Create realm: ${REALM}"
log "   2. Create client: ${OIDC_CLIENT_ID}"
log "   3. Configure admin-cli client with realm-admin role"
log "   4. Create users and assign roles in Keycloak admin console"
log ""
log "ℹ️ See README.md for detailed Keycloak setup instructions"
