FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-issuer-not-validated/ ./challenges/jwt-issuer-not-validated/

WORKDIR /app/challenges/jwt-issuer-not-validated
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-issuer-not-validated .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-issuer-not-validated /usr/bin/jwt-issuer-not-validated

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-issuer-not-validated"]
CMD ["serve"]
