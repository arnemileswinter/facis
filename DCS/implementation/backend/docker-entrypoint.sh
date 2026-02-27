#!/bin/bash
set -e

# If custom CA certificates exist, update the CA bundle
if [ -d "/usr/local/share/ca-certificates/custom" ] && [ "$(ls -A /usr/local/share/ca-certificates/custom 2>/dev/null)" ]; then
    echo "Custom CA certificates found, updating CA bundle..."
    # This requires root, so we'll handle it differently
    # Instead, we'll concatenate custom certs to the system bundle
    if [ -f /usr/local/share/ca-certificates/custom/*.crt ]; then
        cat /etc/ssl/certs/ca-certificates.crt /usr/local/share/ca-certificates/custom/*.crt > /tmp/ca-bundle.crt
        export SSL_CERT_FILE=/tmp/ca-bundle.crt
        echo "Custom CA certificates loaded. SSL_CERT_FILE=$SSL_CERT_FILE"
    fi
fi

# Execute the main command
exec "$@"
