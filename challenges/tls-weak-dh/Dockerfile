FROM nginx:1.27-alpine

# hadolint ignore=DL3018
RUN apk add --no-cache openssl \
    && mkdir -p /etc/nginx/tls \
    && openssl req -x509 -newkey rsa:2048 -nodes -days 3650 -subj "/CN=localhost" \
        -addext "subjectAltName=DNS:localhost,IP:127.0.0.1" \
        -keyout /etc/nginx/tls/server.key -out /etc/nginx/tls/server.crt \
    && openssl dhparam -out /etc/nginx/tls/dhparam1024.pem 1024

COPY challenges/tls-weak-dh/conf/nginx.vulnerable.conf /etc/nginx/conf.d/nginx.vulnerable.conf
COPY challenges/tls-weak-dh/conf/nginx.fixed.conf /etc/nginx/conf.d/nginx.fixed.conf
COPY challenges/tls-weak-dh/docker-entrypoint.sh /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh

EXPOSE 8080

ENV VULNERABLE=true

ENTRYPOINT ["/docker-entrypoint.sh"]
