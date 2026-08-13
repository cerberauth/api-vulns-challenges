FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-apple-token-relay/ ./challenges/jwt-apple-token-relay/

WORKDIR /app/challenges/jwt-apple-token-relay
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-apple-token-relay .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-apple-token-relay/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-apple-token-relay /usr/bin/jwt-apple-token-relay

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-apple-token-relay"]
CMD ["serve"]
