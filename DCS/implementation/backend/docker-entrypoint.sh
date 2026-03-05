#!/bin/bash
set -e

# Generate frontend runtime configuration
CONFIG_FILE="/app/web/dist/config.js"
INDEX_FILE="/app/web/dist/index.html"
BASE_PATH="${DCS_UI_BASE_PATH:-/ui/}"
API_BASE_URL="${DCS_API_BASE_URL:-/}"

if [[ "${BASE_PATH}" != /* ]]; then
  BASE_PATH="/${BASE_PATH}"
fi
if [[ "${BASE_PATH}" != */ ]]; then
  BASE_PATH="${BASE_PATH}/"
fi

cat > "$CONFIG_FILE" << EOF
window.DCS_CONFIG = {
  API_BASE_URL: '${API_BASE_URL}',
}
EOF

if [ -f "$INDEX_FILE" ]; then
  sed -i "s|__DCS_UI_BASE_PATH__|${BASE_PATH}|g" "$INDEX_FILE"
fi

# If custom CA certificates exist, update the CA bundle
if [ -d "/usr/local/share/ca-certificates/custom" ] && [ "$(ls -A /usr/local/share/ca-certificates/custom 2>/dev/null)" ]; then
    if [ -f /usr/local/share/ca-certificates/custom/*.crt ]; then
        cat /etc/ssl/certs/ca-certificates.crt /usr/local/share/ca-certificates/custom/*.crt > /tmp/ca-bundle.crt
        export SSL_CERT_FILE=/tmp/ca-bundle.crt
    fi
fi

exec "$@"
