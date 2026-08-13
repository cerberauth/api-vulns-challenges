FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-jwks-issuer-injection/ ./challenges/jwt-jwks-issuer-injection/

WORKDIR /app/challenges/jwt-jwks-issuer-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-jwks-issuer-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-jwks-issuer-injection/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-jwks-issuer-injection /usr/bin/jwt-jwks-issuer-injection

EXPOSE 8080
EXPOSE 8090

USER nonroot:nonroot

ENTRYPOINT ["jwt-jwks-issuer-injection"]
CMD ["serve"]
