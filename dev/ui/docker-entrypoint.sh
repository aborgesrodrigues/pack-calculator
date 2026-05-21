#!/bin/sh
set -e

export API_UPSTREAM="${API_UPSTREAM:-app:8080}"
export NGINX_RESOLVER="${NGINX_RESOLVER:-127.0.0.11 valid=10s}"

envsubst '${API_UPSTREAM} ${NGINX_RESOLVER}' \
  < /etc/nginx/templates/default.conf.template \
  > /etc/nginx/conf.d/default.conf

exec nginx -g 'daemon off;'
