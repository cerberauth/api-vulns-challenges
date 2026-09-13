#!/bin/sh
set -e

VULNERABLE="${VULNERABLE:-true}"

if [ "$VULNERABLE" = "true" ]; then
    cp /etc/nginx/conf.d/nginx.vulnerable.conf /etc/nginx/nginx.conf
else
    cp /etc/nginx/conf.d/nginx.fixed.conf /etc/nginx/nginx.conf
fi

exec nginx -g "daemon off;"
