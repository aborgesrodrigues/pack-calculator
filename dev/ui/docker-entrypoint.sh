#!/bin/sh
set -e

# Railway: private DNS + api service name. Local Compose: Docker service "app".
if [ -n "${RAILWAY_ENVIRONMENT:-}" ] || [ -n "${RAILWAY_SERVICE_ID:-}" ]; then
  export API_UPSTREAM="${API_UPSTREAM:-api.railway.internal:8080}"
  export NGINX_RESOLVER="${NGINX_RESOLVER:-[fd12::10] ipv6=on valid=1s}"
else
  export API_UPSTREAM="${API_UPSTREAM:-app:8080}"
  export NGINX_RESOLVER="${NGINX_RESOLVER:-127.0.0.11 valid=10s}"
fi

envsubst '${API_UPSTREAM} ${NGINX_RESOLVER}' \
  < /etc/nginx/templates/default.conf.template \
  > /etc/nginx/conf.d/default.conf

nginx -t

exec nginx -g 'daemon off;'
