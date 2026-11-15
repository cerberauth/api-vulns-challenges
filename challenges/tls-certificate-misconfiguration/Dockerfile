FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/tls-certificate-misconfiguration/ ./challenges/tls-certificate-misconfiguration/

WORKDIR /app/challenges/tls-certificate-misconfiguration
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /tls-certificate-misconfiguration .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /tls-certificate-misconfiguration /usr/bin/tls-certificate-misconfiguration

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["tls-certificate-misconfiguration"]
CMD ["serve"]
