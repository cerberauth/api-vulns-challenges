FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/tls-known-cves/ ./challenges/tls-known-cves/

WORKDIR /app/challenges/tls-known-cves
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /tls-known-cves .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /tls-known-cves /usr/bin/tls-known-cves

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["tls-known-cves"]
CMD ["serve"]
