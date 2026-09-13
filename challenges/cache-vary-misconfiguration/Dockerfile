FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-vary-misconfiguration/ ./challenges/cache-vary-misconfiguration/

WORKDIR /app/challenges/cache-vary-misconfiguration
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-vary-misconfiguration .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-vary-misconfiguration /usr/bin/cache-vary-misconfiguration

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-vary-misconfiguration"]
CMD ["serve"]
