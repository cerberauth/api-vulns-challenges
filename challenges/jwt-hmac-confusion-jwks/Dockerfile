FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-hmac-confusion-jwks/ ./challenges/jwt-hmac-confusion-jwks/

WORKDIR /app/challenges/jwt-hmac-confusion-jwks
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-hmac-confusion-jwks .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-hmac-confusion-jwks /usr/bin/jwt-hmac-confusion-jwks

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-hmac-confusion-jwks"]
CMD ["serve"]
