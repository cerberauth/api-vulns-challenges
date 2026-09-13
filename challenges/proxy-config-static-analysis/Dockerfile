FROM busybox:1.36

WORKDIR /configs

COPY challenges/proxy-config-static-analysis/configs/ ./

EXPOSE 8080

# this challenge ships static config files for a config-level lint/static
# analysis subcommand rather than a running proxy; the container just serves
# the files over HTTP so they can be fetched for out-of-band analysis
ENTRYPOINT ["httpd", "-f", "-p", "8080", "-h", "/configs"]
