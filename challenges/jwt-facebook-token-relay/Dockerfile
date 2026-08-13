FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-facebook-token-relay/ ./challenges/jwt-facebook-token-relay/

WORKDIR /app/challenges/jwt-facebook-token-relay
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-facebook-token-relay .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-facebook-token-relay/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-facebook-token-relay /usr/bin/jwt-facebook-token-relay

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-facebook-token-relay"]
CMD ["serve"]
