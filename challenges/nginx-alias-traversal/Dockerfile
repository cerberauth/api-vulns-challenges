FROM nginx:1.27-alpine

COPY challenges/nginx-alias-traversal/conf/nginx.vulnerable.conf /etc/nginx/conf.d/nginx.vulnerable.conf
COPY challenges/nginx-alias-traversal/conf/nginx.fixed.conf /etc/nginx/conf.d/nginx.fixed.conf
COPY challenges/nginx-alias-traversal/html/ /usr/share/nginx/html/
COPY challenges/nginx-alias-traversal/docker-entrypoint.sh /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh

EXPOSE 8080

ENV VULNERABLE=true

ENTRYPOINT ["/docker-entrypoint.sh"]
