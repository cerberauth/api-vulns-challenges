FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-tls-misconfiguration/ ./challenges/proxy-tls-misconfiguration/

WORKDIR /app/challenges/proxy-tls-misconfiguration
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-tls-misconfiguration .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-tls-misconfiguration /usr/bin/proxy-tls-misconfiguration

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-tls-misconfiguration"]
CMD ["serve"]
