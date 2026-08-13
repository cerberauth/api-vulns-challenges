FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-google-token-relay/ ./challenges/jwt-google-token-relay/

WORKDIR /app/challenges/jwt-google-token-relay
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-google-token-relay .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-google-token-relay/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-google-token-relay /usr/bin/jwt-google-token-relay

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-google-token-relay"]
CMD ["serve"]
