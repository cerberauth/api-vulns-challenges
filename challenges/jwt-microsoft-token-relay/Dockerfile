FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-microsoft-token-relay/ ./challenges/jwt-microsoft-token-relay/

WORKDIR /app/challenges/jwt-microsoft-token-relay
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-microsoft-token-relay .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-microsoft-token-relay/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-microsoft-token-relay /usr/bin/jwt-microsoft-token-relay

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-microsoft-token-relay"]
CMD ["serve"]
