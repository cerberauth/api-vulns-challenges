FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-jwk-header-injection/ ./challenges/jwt-jwk-header-injection/

WORKDIR /app/challenges/jwt-jwk-header-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-jwk-header-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-jwk-header-injection/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-jwk-header-injection /usr/bin/jwt-jwk-header-injection

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-jwk-header-injection"]
CMD ["serve"]
