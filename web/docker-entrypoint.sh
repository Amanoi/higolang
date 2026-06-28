#!/bin/sh
set -e

# Default domain
DOMAIN="${DOMAIN:-localhost}"

# Choose template based on SSL configuration
if [ -n "$SSL_CERT" ] && [ -f "$SSL_CERT" ]; then
    echo "SSL enabled for domain: $DOMAIN"
    envsubst '${DOMAIN} ${SSL_CERT} ${SSL_KEY}' < /etc/nginx/templates/nginx.conf.template > /etc/nginx/conf.d/default.conf
else
    echo "SSL not configured, running HTTP only for domain: $DOMAIN"
    envsubst '${DOMAIN}' < /etc/nginx/templates/nginx-http.conf.template > /etc/nginx/conf.d/default.conf
fi

# Start nginx
exec nginx -g "daemon off;"
