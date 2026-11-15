FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-missing-expiry/ ./challenges/jwt-missing-expiry/

WORKDIR /app/challenges/jwt-missing-expiry
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-missing-expiry .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-missing-expiry /usr/bin/jwt-missing-expiry

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-missing-expiry"]
CMD ["serve"]
